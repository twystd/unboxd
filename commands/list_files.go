package commands

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/lib"
	"github.com/twystd/unboxd/credentials"
)

var ListFilesCmd = ListFiles{
	command: command{
		name:  "list-files",
		delay: 500 * time.Millisecond,
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

var header = struct {
	normal   []string
	withTags []string
}{
	normal:   []string{"ID", "Folder", "Filename"},
	withTags: []string{"ID", "Folder", "Filename", "Tags"},
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

func (cmd ListFiles) Execute(flagset *flag.FlagSet, c credentials.ICredentials) error {
	credentials := c["box"].(box.Credentials)

	b := box.NewBox()
	if err := b.Authenticate(credentials); err != nil {
		return err
	}

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

	recalc := func(widths []int, record []string) []int {
		for i, field := range record {
			if N := len(field); N > widths[i] {
				widths[i] = N
			}
		}

		return widths
	}

	var hdr []string
	if cmd.tags {
		hdr = header.withTags
	} else {
		hdr = header.normal
	}

	widths := recalc(make([]int, len(hdr)), hdr)
	table := [][]string{hdr}
	for _, file := range files {
		record := cmd.toRecord(file)
		widths = recalc(widths, record)
		table = append(table, record)
	}

	columns := []string{}
	for _, w := range widths {
		columns = append(columns, fmt.Sprintf("%%-%vv", w))
	}

	format := fmt.Sprintf("%v\n", strings.Join(columns, "  "))
	for _, row := range table {
		args := []any{}
		for _, v := range row {
			args = append(args, v)
		}

		fmt.Printf(format, args...)
	}

	return nil
}

func (cmd ListFiles) save(files []file) error {
	infof("list-files", "saving %v files to TSV file %v\n", len(files), cmd.file)

	sort.Slice(files, func(i, j int) bool { return files[i].FilePath < files[j].FilePath })

	var hdr []string
	if cmd.tags {
		hdr = header.normal
	} else {
		hdr = header.withTags
	}

	table := [][]string{hdr}
	for _, file := range files {
		record := cmd.toRecord(file)
		table = append(table, record)
	}

	if err := os.MkdirAll(filepath.Dir(cmd.file), 0750); err != nil {
		return err
	} else if f, err := os.Create(cmd.file); err != nil {
		return err
	} else {
		w := csv.NewWriter(f)
		w.Comma = '\t'
		w.WriteAll(table)

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
				path := item.Path + "/" + f.Name
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

func (cmd ListFiles) toRecord(f file) []string {
	id := fmt.Sprintf("%v", f.ID)
	folder := path.Dir(f.FilePath)
	filename := f.FileName
	tags := strings.Join(f.Tags, "; ")

	var record []string
	if cmd.tags {
		record = []string{
			id,
			folder,
			filename,
			tags,
		}
	} else {
		record = []string{
			id,
			folder,
			filename,
		}
	}

	return record
}
