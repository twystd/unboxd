package box

import (
	"github.com/twystd/unboxd/box/folders"
	"github.com/twystd/unboxd/box/lib"
)

func (b *Box) ListFolders(glob string) ([]folders.Folder, error) {
	list, err := listFolders("0", "", b.token.Token)
	if err != nil {
		return nil, err
	}

	g := lib.NewGlob(glob)
	matched := []folders.Folder{}
	for _, f := range list {
		if g.Match(f.Path) {
			matched = append(matched, f)
		}
	}

	return matched, nil
}

func listFolders(folderID string, prefix string, token string) ([]folders.Folder, error) {
	list := []folders.Folder{}

	l, err := folders.List(folderID, token)
	if err != nil {
		return nil, err
	}

	for k, v := range l {
		path := prefix + "/" + v
		list = append(list, folders.Folder{
			ID:   k,
			Name: v,
			Path: path,
		})
	}

	for _, f := range list {
		if l, err := listFolders(f.ID, f.Path, token); err != nil {
			return nil, err
		} else {
			list = append(list, l...)
		}
	}

	return list, nil
}
