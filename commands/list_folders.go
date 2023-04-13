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

func (cmd ListFolders) Help() {
	fmt.Println()
	fmt.Println("  Usage: unboxd [--debug] --credentials <file> list-folders [--tags] [--file <file>] [--checkpoint <file>] [--delay <duration>] [--no-resume] <folderspec>")
	fmt.Println()
	fmt.Println("  Retrieves a list of folders that match the folder spec.")
	fmt.Println()
	fmt.Println("  A folderspec is a glob expression against which to match folder paths e.g.:")
	fmt.Println("    /          matches top level folders")
	fmt.Println("    /**        matches all folders recursively")
	fmt.Println("    /photos/*  matches all folders in the /photos folder")
	fmt.Println()
	fmt.Println("  The default folderspec is /** i.e. list all folders recursively")
	fmt.Println()
	fmt.Println("    --credentials <file>  JSON file with Box credentials (required)")
	fmt.Println("    --tags                Include tags in folder information")
	fmt.Println("    --file                TSV file to which to write folder information")
	fmt.Println("    --no-resume           Retrieves folder list from the beginning (default is to continue from the last checkpoint)")
	fmt.Println("    --checkpoint          Specifies the path for the checkpoint file (default is .checkpoint)")
	fmt.Println("    --batch               Maximum number of calls to the Box API (defaults to no limit)")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    --delay  Delay between multiple requests to reduce traffic to Box API")
	fmt.Println("    --debug  Enable debugging information")
	fmt.Println()
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    unboxd --debug --credentials .credentials list-folders --tags --file folders.tsv /**"`)
	fmt.Println()
}

func (cmd *ListFolders) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.BoolVar(&cmd.tags, "tags", cmd.tags, "Include tags in folder information")
	flagset.StringVar(&cmd.file, "file", cmd.file, "TSV file to which to write folder information")
	flagset.StringVar(&cmd.checkpoint, "checkpoint", cmd.checkpoint, "Specifies the path for the checkpoint file")
	flag.DurationVar(&cmd.delay, "delay", cmd.delay, "Delay between multiple requests to reduce traffic to Box API")
	flagset.BoolVar(&cmd.restart, "no-resume", cmd.restart, "Retrieves folder list from the beginning")
	flagset.UintVar(&cmd.batch, "batch-size", cmd.batch, "Number of calls to the Box API")

	return flagset
}

func (cmd ListFolders) Execute(flagset *flag.FlagSet, b box.Box) error {
	var base string

	args := flagset.Args()
	if len(args) > 0 {
		base = args[0]
	} else {
		base = ""
	}

	// .. get folder list
	list, err := cmd.exec(b, base)
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

func (cmd ListFolders) exec(b box.Box, glob string) ([]folder, error) {
	list := []folder{}

	folders, err := listFolders(b, 0, "", cmd.checkpoint, cmd.delay, cmd.restart, cmd.batch)
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
		w.WriteAll(records)

		return w.Error()
	}
}

func listFolders(b box.Box, folderID uint64, prefix string, chkpt string, delay time.Duration, restart bool, batch uint) ([]folder, error) {
	pipe, folders, _, err := resume(chkpt, restart)
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
		time.Sleep(delay)

		item := pipe[tail]
		if l, err := b.ListFolders(item.ID); err != nil {
			if errx := checkpoint(chkpt, pipe[tail:], folders, []file{}); errx != nil {
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
		if batch != 0 && count > batch {
			break
		}

		tail++
	}

	// ... incomplete?
	if len(pipe[tail:]) > 0 {
		if err := checkpoint(chkpt, pipe[tail:], folders, []file{}); err != nil {
			return folders, err
		} else {
			return folders, fmt.Errorf("interrupted")
		}
	}

	// ... complete!
	if err := checkpoint(chkpt, []QueueItem{}, []folder{}, []file{}); err != nil {
		return folders, err
	}

	return folders, nil
}
