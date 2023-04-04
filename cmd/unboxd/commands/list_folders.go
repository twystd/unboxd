package commands

import (
	"flag"
	"fmt"
	"sort"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/lib"
)

var ListFoldersCmd = ListFolders{
	delay: 500 * time.Millisecond,
}

type ListFolders struct {
	delay time.Duration
}

type folder struct {
	ID   uint64 `json:"ID"`
	Name string `json:"name"`
	Path string `json:"path"`
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

	folders, err := listFolders(b, 0, "", cmd.delay)
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

func listFolders(b box.Box, folderID uint64, prefix string, delay time.Duration) ([]folder, error) {
	folders := []folder{}
	pipe := []uint64{}
	tail := 0

	pipe = append(pipe, folderID)

	for tail < len(pipe) {
		time.Sleep(delay)

		l, err := b.ListFolders(pipe[tail])

		if err != nil {
			if errx := checkpoint("./runtime/.checkpoint", pipe[tail:], folders); errx != nil {
				warnf("list-folders", "%v", errx)
			}

			return folders, err
		}

		for _, f := range l {
			path := prefix + "/" + f.Name
			folders = append(folders, folder{
				ID:   f.ID,
				Name: f.Name,
				Path: path,
			})

			pipe = append(pipe, f.ID)
		}

		tail++
	}

	if len(pipe[tail:]) > 0 {
		if err := checkpoint("./runtime/.checkpoint", pipe[tail:], folders); err != nil {
			return folders, err
		}
	} else {
		if err := checkpoint(".checkpoint", pipe[tail:], []folder{}); err != nil {
			return folders, err
		}
	}

	return folders, nil
}
