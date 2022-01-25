package commands

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/files"
)

type ListFiles struct {
}

func (cmd ListFiles) Name() string {
	return "list-files"
}

func (cmd ListFiles) Execute(b box.Box) error {
	args := flag.Args()[1:]
	if len(args) < 1 {
		return fmt.Errorf("missing folder argument")
	}

	folder := args[0]

	files, err := cmd.exec(b, folder)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no files")
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Filename < files[j].Filename })

	widths := []int{0, 0, 0}
	table := [][3]string{}

	for _, f := range files {
		id := fmt.Sprintf("%v", f.ID)
		filename := fmt.Sprintf("%v", f.Filename)
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

func (cmd ListFiles) exec(b box.Box, folder string) ([]files.File, error) {
	return b.ListFiles(folder)
}
