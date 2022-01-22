package commands

import (
	"flag"
	"fmt"
	"sort"

	"github.com/twystd/unboxd/box"
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

	widths := []int{0, 0}
	table := [][2]string{}

	for _, f := range files {
		id := fmt.Sprintf("%v", f.ID)
		filename := fmt.Sprintf("%v", f.Filename)

		if N := len(id); N > widths[0] {
			widths[0] = N
		}

		if N := len(filename); N > widths[1] {
			widths[1] = N
		}

		table = append(table, [2]string{id, filename})
	}

	format := fmt.Sprintf("%%-%vv  %%-%vv\n", widths[0], widths[1])
	for _, row := range table {
		fmt.Printf(format, row[0], row[1])
	}

	return nil
}

func (cmd ListFiles) exec(b box.Box, folder string) ([]box.File, error) {
	return b.ListFiles(folder)
}
