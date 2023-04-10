package commands

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/lib"
)

var ListFoldersCmd = ListFolders{
	command: command{
		name:  "list-folders",
		delay: 500 * time.Millisecond,
	},

	file:       "",
	checkpoint: ".checkpoint",
	tags:       false,
}

type ListFolders struct {
	command
	file       string
	checkpoint string
	tags       bool
}

type folder struct {
	ID   uint64   `json:"ID"`
	Name string   `json:"name"`
	Path string   `json:"path"`
	Tags []string `json:"tags,omitempty"`
}

func (cmd *ListFolders) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.BoolVar(&cmd.tags, "tags", cmd.tags, "(optional) include tags in folder information")
	flagset.StringVar(&cmd.file, "file", cmd.file, "(optional) TSV file to which to write folder information")
	flagset.StringVar(&cmd.checkpoint, "checkpoint", cmd.checkpoint, "(optional) specifies the path for the checkpoint file")
	flag.DurationVar(&cmd.delay, "delay", cmd.delay, "(optional) delay between multiple requests to reduce traffic to Box API")

	return flagset
}

func (cmd ListFolders) Help() {
}

func (cmd ListFolders) Execute(flagset *flag.FlagSet, b box.Box) error {
	folder := ""

	args := flagset.Args()
	if len(args) > 0 {
		folder = args[0]
	}

	folders, err := cmd.exec(b, folder)
	if err != nil {
		return err
	}

	if len(folders) == 0 {
		return fmt.Errorf("no folders matching path '%s", folder)
	}

	sort.Slice(folders, func(i, j int) bool { return folders[i].Path < folders[j].Path })

	if cmd.file != "" {
		return cmd.save(folders)
	} else {
		return cmd.print(folders)
	}
}

func (cmd ListFolders) exec(b box.Box, glob string) ([]folder, error) {
	list := []folder{}

	folders, err := listFolders(b, 0, "", cmd.checkpoint, cmd.delay)
	if err != nil {
		return nil, err
	}

	g := lib.NewGlob(glob)
	for _, f := range folders {
		if g.Match(f.Path) {
			list = append(list, f)
		}
	}

	return list, nil
}

func (cmd ListFolders) print(folders []folder) error {
	infof("list-folders", "saving %v folders to TSV file %v\n", len(folders), cmd.file)

	sort.Slice(folders, func(i, j int) bool { return folders[i].Path < folders[j].Path })

	widths := []int{0, 0, 0}
	table := [][]string{
		[]string{"ID", "Path", "Tags"},
	}

	for _, f := range folders {
		id := fmt.Sprintf("%v", f.ID)
		path := fmt.Sprintf("%v", f.Path)
		tags := fmt.Sprintf("%v", strings.Join(f.Tags, ";"))

		if N := len(id); N > widths[0] {
			widths[0] = N
		}

		if N := len(path); N > widths[1] {
			widths[1] = N
		}

		if cmd.tags {
			if N := len(tags); N > widths[2] {
				widths[2] = N
			}
		}

		table = append(table, []string{id, path, tags})
	}

	var format string
	if cmd.tags {
		format = fmt.Sprintf("%%-%vv  %%-%vv  %%-%vv\n", widths[0], widths[1], widths[2])
	} else {
		format = fmt.Sprintf("%%-%vv  %%-%vv\n", widths[0], widths[1])
	}

	if cmd.tags {
		for _, row := range table {
			fmt.Printf(format, row[0], row[1], row[2])
		}
	} else {
		for _, row := range table {
			fmt.Printf(format, row[0], row[1])
		}
	}

	return nil
}

func (cmd ListFolders) save(folders []folder) error {
	sort.Slice(folders, func(i, j int) bool { return folders[i].Path < folders[j].Path })

	records := [][]string{
		[]string{"ID", "Path"},
	}

	if cmd.tags {
		records = [][]string{
			[]string{"ID", "Path", "Tags"},
		}
	}

	for _, f := range folders {
		id := fmt.Sprintf("%v", f.ID)
		path := fmt.Sprintf("%v", f.Path)
		tags := fmt.Sprintf("%v", strings.Join(f.Tags, ";"))

		if cmd.tags {
			records = append(records, []string{id, path, tags})
		} else {
			records = append(records, []string{id, path})
		}
	}

	if err := os.MkdirAll(filepath.Dir(cmd.file), 0750); err != nil {
		return err
	} else if f, err := os.Create(cmd.file); err != nil {
		return err
	} else {
		w := csv.NewWriter(f)
		w.WriteAll(records)

		return w.Error()
	}
}

func listFolders(b box.Box, folderID uint64, prefix string, chkpt string, delay time.Duration) ([]folder, error) {
	tail := 0

	if pipe, folders, err := resume(chkpt); err != nil {
		return nil, err
	} else {
		if len(pipe) > 0 {
			infof("list-folders", "Resuming last operation")
		}

		if len(pipe) == 0 {
			pipe = append(pipe, QueueItem{ID: folderID, Path: prefix})
		}

		for tail < len(pipe) {
			time.Sleep(delay)

			item := pipe[tail]
			l, err := b.ListFolders(item.ID)

			if err != nil {
				if errx := checkpoint(chkpt, pipe[tail:], folders); errx != nil {
					warnf("list-folders", "%v", errx)
				}

				return folders, err
			}

			for _, f := range l {
				path := item.Path + "/" + f.Name
				folders = append(folders, folder{
					ID:   f.ID,
					Name: f.Name,
					Tags: f.Tags,
					Path: path,
				})

				pipe = append(pipe, QueueItem{ID: f.ID, Path: path})
			}

			tail++
		}

		// ... incomplete?
		if len(pipe[tail:]) > 0 {
			if err := checkpoint(chkpt, pipe[tail:], folders); err != nil {
				return folders, err
			} else {
				return folders, fmt.Errorf("interrupted")
			}
		}

		// ... complete!
		if err := checkpoint(chkpt, []QueueItem{}, []folder{}); err != nil {
			return folders, err
		}

		return folders, nil
	}
}