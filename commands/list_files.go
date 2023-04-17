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

var ListFilesCmd = ListFiles{
	command: command{
		application: APP,
		name:        "list-files",
		delay:       500 * time.Millisecond,
	},

	file:       "",
	checkpoint: ".checkpoint",
	tags:       false,
	restart:    false,
	batch:      0,
}

type ListFiles struct {
	command
	file       string
	checkpoint string
	tags       bool
	restart    bool
	batch      uint
}

type file struct {
	ID       uint64
	FileName string
	FilePath string
	Tags     []string
}

func (cmd ListFiles) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %v [--debug] --credentials <file> list-files [--tags] [--file <file>] [--checkpoint <file>] [--delay <duration>] [--no-resume] <filespec>\n", APP)
	fmt.Println()
	fmt.Println("  Retrieves a list of files that match the file spec.")
	fmt.Println()
	fmt.Println("  A filespec is a glob expression against which to match file paths e.g.:")
	fmt.Println("    /*         matches files in the top level folder")
	fmt.Println("    /**        matches all files recursively")
	fmt.Println("    /photos/*  matches all files in the /photos folder")
	fmt.Println()
	fmt.Println("  The default filespec is /** i.e. list all files recursively")
	fmt.Println()
	fmt.Println("    --credentials <file>  JSON file with Box credentials (required)")
	fmt.Println("    --tags                Include tags in file information")
	fmt.Println("    --file                TSV file to which to write file information")
	fmt.Println("    --no-resume           Retrieves file list from the beginning (default is to continue from last checkpoint")
	fmt.Println("    --checkpoint          Specifies the path for the checkpoint file (default is .checkpoint)")
	fmt.Println("    --batch               Maximum number of calls to the Box API (defaults to no limit)")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    --delay  Delay between multiple requests to reduce traffic to Box API")
	fmt.Println("    --debug  Enable debugging information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Printf("    %v --debug --credentials .credentials list-files --tags --file folders.tsv /**\n", APP)
	fmt.Println()
}

func (cmd *ListFiles) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.BoolVar(&cmd.tags, "tags", cmd.tags, "Include tags in folder information")
	flagset.StringVar(&cmd.file, "file", cmd.file, "TSV file to which to write folder information")
	flagset.StringVar(&cmd.checkpoint, "checkpoint", cmd.checkpoint, "Specifies the path for the checkpoint file")
	flagset.DurationVar(&cmd.delay, "delay", cmd.delay, "Delay between multiple requests to reduce traffic to Box API")
	flagset.BoolVar(&cmd.restart, "no-resume", cmd.restart, "Retrieves folder list from the beginning")
	flagset.UintVar(&cmd.batch, "batch-size", cmd.batch, "Number of calls to the Box API")

	return flagset
}

func (cmd ListFiles) Execute(flagset *flag.FlagSet, b box.Box) error {
	glob := ""

	args := flagset.Args()
	if len(args) > 0 {
		glob = args[0]
	}

	hash := cmd.hash("list-files", b.Hash(), glob)

	// .. get files
	list, err := cmd.exec(b, glob, hash)
	if err != nil {
		return err
	}

	// .. dedupe and sort
	files := []file{}
	dedupe := map[uint64]bool{}
	for _, f := range list {
		if duplicate := dedupe[f.ID]; !duplicate {
			files = append(files, f)
			dedupe[f.ID] = true
		}
	}

	// .. save/print
	if cmd.file != "" {
		return cmd.save(files)
	} else {
		return cmd.print(files)
	}
}

func (cmd ListFiles) exec(b box.Box, glob string, hash string) ([]file, error) {
	list := []file{}

	folders, err := cmd.listFiles(b, 0, "", hash)
	if err != nil {
		return nil, err
	}

	g := lib.NewGlob(glob)
	for _, f := range folders {
		if g.Match(f.FilePath) {
			list = append(list, f)
		}
	}

	return list, nil
}

func (cmd ListFiles) print(files []file) error {
	sort.Slice(files, func(i, j int) bool { return files[i].FilePath < files[j].FilePath })

	widths := []int{0, 0, 0}
	table := [][3]string{}

	for _, f := range files {
		id := fmt.Sprintf("%v", f.ID)
		filename := fmt.Sprintf("%v", f.FileName)
		tags := strings.Join(f.Tags, ",")

		if N := len(id); N > widths[0] {
			widths[0] = N
		}

		if N := len(filename); N > widths[1] {
			widths[1] = N
		}

		if N := len(tags); N > widths[2] {
			widths[2] = N
		}

		table = append(table, [3]string{id, filename, tags})
	}

	format := fmt.Sprintf("%%-%vv  %%-%vv  %%-%vv\n", widths[0], widths[1], widths[2])
	for _, row := range table {
		fmt.Printf(format, row[0], row[1], row[2])
	}

	return nil
}

func (cmd ListFiles) save(files []file) error {
	infof("list-files", "saving %v files to TSV file %v\n", len(files), cmd.file)

	sort.Slice(files, func(i, j int) bool { return files[i].FilePath < files[j].FilePath })

	records := [][]string{
		[]string{"ID", "Path"},
	}

	if cmd.tags {
		records = [][]string{
			[]string{"ID", "Path", "Tags"},
		}
	}

	for _, f := range files {
		id := fmt.Sprintf("%v", f.ID)
		path := fmt.Sprintf("%v", f.FilePath)
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

func (cmd ListFiles) listFiles(b box.Box, folderID uint64, prefix string, hash string) ([]file, error) {
	pipe, folders, files, err := resume(cmd.checkpoint, hash, cmd.restart)
	if err != nil {
		return nil, err
	}

	if len(pipe) > 0 {
		infof("list-files", "Resuming last operation")
	} else {
		pipe = append(pipe, QueueItem{ID: folderID, Path: prefix})
	}

	count := uint(0)
	tail := 0
	for tail < len(pipe) {
		item := pipe[tail]

		// get files for current folder
		if l, err := b.ListFiles(item.ID); err != nil {
			if errx := checkpoint(cmd.checkpoint, pipe[tail:], folders, files, hash); errx != nil {
				warnf("list-files", "%v", errx)
			}

			return files, err
		} else {
			for _, f := range l {
				path := prefix + "/" + f.Name
				files = append(files, file{
					ID:       f.ID,
					FileName: f.Name,
					FilePath: path,
					Tags:     f.Tags,
				})
			}
		}

		// get subfolders for current folder
		if l, err := b.ListFolders(item.ID); err != nil {
			if errx := checkpoint(cmd.checkpoint, pipe[tail:], folders, files, hash); errx != nil {
				warnf("list-files", "%v", errx)
			}

			return files, err
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
		if err := checkpoint(cmd.checkpoint, pipe[tail:], folders, files, hash); err != nil {
			return files, err
		} else {
			return files, fmt.Errorf("interrupted")
		}
	}

	// ... complete!
	if err := checkpoint(cmd.checkpoint, []QueueItem{}, []folder{}, []file{}, ""); err != nil {
		return files, err
	}

	return files, nil
}
