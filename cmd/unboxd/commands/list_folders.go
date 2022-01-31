package commands

import (
	"flag"
	"fmt"
	"sort"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/folders"
)

type ListFolders struct {
}

func (cmd ListFolders) Name() string {
	return "list-folders"
}

func (cmd ListFolders) Execute(b box.Box) error {
	folder := ""

	args := flag.Args()[1:]
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

	widths := []int{0, 0}
	table := [][2]string{}

	for _, f := range folders {
		id := fmt.Sprintf("%v", f.ID)
		path := fmt.Sprintf("%v", f.Path)

		if N := len(id); N > widths[0] {
			widths[0] = N
		}

		if N := len(path); N > widths[1] {
			widths[1] = N
		}

		table = append(table, [2]string{id, path})
	}

	format := fmt.Sprintf("%%-%vv  %%-%vv\n", widths[0], widths[1])
	for _, row := range table {
		fmt.Printf(format, row[0], row[1])
	}

	return nil
}

func (cmd ListFolders) exec(b box.Box, folder string) ([]folders.Folder, error) {
	return b.ListFolders(folder)
}
