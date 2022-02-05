package commands

import (
	"flag"
	"fmt"
	"sort"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/lib"
)

type ListFolders struct {
}

type folder struct {
	ID   uint64
	Name string
	Path string
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

func (cmd ListFolders) exec(b box.Box, glob string) ([]folder, error) {
	list := []folder{}

	folders, err := listFolders(b, 0, "")
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

func listFolders(b box.Box, folderID uint64, prefix string) ([]folder, error) {
	folders := []folder{}

	l, err := b.ListFolders(folderID)
	if err != nil {
		return nil, err
	}

	for _, f := range l {
		path := prefix + "/" + f.Name
		folders = append(folders, folder{
			ID:   f.ID,
			Name: f.Name,
			Path: path,
		})
	}

	for _, f := range folders {
		if l, err := listFolders(b, f.ID, f.Path); err != nil {
			return nil, err
		} else {
			folders = append(folders, l...)
		}
	}

	return folders, nil
}
