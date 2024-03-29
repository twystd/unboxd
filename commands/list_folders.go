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
	"github.com/twystd/unboxd/credentials"
)

var ListFoldersCmd = ListFolders{
	command: command{
		name:  "list-folders",
		delay: 500 * time.Millisecond,
	},

	file:       "",
	checkpoint: ".checkpoint",
	tags:       false,
	restart:    false,
	batch:      0,
}

type ListFolders struct {
	command
	file       string
	checkpoint string
	tags       bool
	restart    bool
	batch      uint
}

type folder struct {
	ID   uint64   `json:"ID"`
	Name string   `json:"name"`
	Path string   `json:"path"`
	Tags []string `json:"tags,omitempty"`
}

func (cmd *ListFolders) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.BoolVar(&cmd.tags, "tags", cmd.tags, "Include tags in folder information")
	flagset.StringVar(&cmd.file, "file", cmd.file, "TSV file to which to write folder information")
	flagset.StringVar(&cmd.checkpoint, "checkpoint", cmd.checkpoint, "Specifies the path for the checkpoint file")
	flagset.DurationVar(&cmd.delay, "delay", cmd.delay, "Delay between multiple requests to reduce traffic to Box API")
	flagset.BoolVar(&cmd.restart, "no-resume", cmd.restart, "Retrieves folder list from the beginning")
	flagset.UintVar(&cmd.batch, "batch-size", cmd.batch, "Number of calls to the Box API")

	return flagset
}

func (cmd ListFolders) Execute(flagset *flag.FlagSet, c credentials.ICredentials) error {
	credentials := c["box"].(box.Credentials)

	b := box.NewBox()
	if err := b.Authenticate(credentials); err != nil {
		return err
	}

	var base string

	args := flagset.Args()
	if len(args) > 0 {
		base = args[0]
	} else {
		base = ""
	}

	hash := cmd.hash("list-folders", b.Hash(), base)

	// .. get folder list
	list, err := cmd.exec(b, base, hash)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		warnf("list-folder", "no folders matching '%s", base)
	}

	// .. dedupe
	folders := []folder{}
	dedupe := map[uint64]bool{}
	for _, f := range list {
		if duplicate := dedupe[f.ID]; !duplicate {
			dedupe[f.ID] = true
			folders = append(folders, f)
		}
	}

	// .. save/print
	if cmd.file != "" {
		return cmd.save(folders)
	} else {
		return cmd.print(folders)
	}
}

func (cmd ListFolders) exec(b box.Box, glob string, hash string) ([]folder, error) {
	list := []folder{}

	folders, err := cmd.listFolders(b, 0, "", hash)
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
	sort.Slice(folders, func(i, j int) bool { return folders[i].Path < folders[j].Path })

	table := [][]string{
		[]string{"ID", "Path", "Tags"},
	}

	widths := []int{
		len(table[0][0]),
		len(table[0][1]),
		len(table[0][2]),
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
	infof("list-folders", "saving %v folders to TSV file %v\n", len(folders), cmd.file)

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
		w.Comma = '\t'
		w.WriteAll(records)

		return w.Error()
	}
}

func (cmd ListFolders) listFolders(b box.Box, folderID uint64, prefix string, hash string) ([]folder, error) {
	pipe, folders, _, err := resume(cmd.checkpoint, hash, cmd.restart)
	if err != nil {
		return nil, err
	}

	if len(pipe) > 0 {
		infof("list-folders", "Resuming last operation")
	} else {
		pipe = append(pipe, QueueItem{ID: folderID, Path: prefix})
	}

	count := uint(0)
	tail := 0
	for tail < len(pipe) {
		item := pipe[tail]
		if l, err := b.ListFolders(item.ID); err != nil {
			if errx := checkpoint(cmd.checkpoint, pipe[tail:], folders, []file{}, hash); errx != nil {
				warnf("list-folders", "%v", errx)
			}

			return folders, err
		} else {
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
		}

		count++
		if cmd.batch != 0 && count > cmd.batch {
			break
		}

		tail++
		if tail < len(pipe) {
			time.Sleep(cmd.delay)
		}
	}

	// ... incomplete?
	if len(pipe[tail:]) > 0 {
		if err := checkpoint(cmd.checkpoint, pipe[tail:], folders, []file{}, hash); err != nil {
			return folders, err
		} else {
			return folders, fmt.Errorf("interrupted")
		}
	}

	// ... complete!
	if err := checkpoint(cmd.checkpoint, []QueueItem{}, []folder{}, []file{}, ""); err != nil {
		return folders, err
	}

	return folders, nil
}
