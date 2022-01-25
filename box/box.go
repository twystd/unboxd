package box

import (
	"fmt"
	"strings"

	"github.com/twystd/unboxd/box/credentials"
	"github.com/twystd/unboxd/box/files"
	"github.com/twystd/unboxd/box/folders"
)

type Box struct {
	token *credentials.AccessToken
}

func NewBox() Box {
	return Box{}
}

const fetchSize = 128

func (b *Box) Authenticate(credentials credentials.Credentials) error {
	if b.token != nil && b.token.IsValid() {
		return nil
	}

	token, err := credentials.Authenticate()
	if err != nil {
		return err
	}

	b.token = token

	return nil
}

func (b *Box) ListFiles(folder string) ([]files.File, error) {
	prefix := ""
	folderID := "0"

loop:
	for {
		folders, err := folders.ListFolders(folderID, b.token.Token)
		if err != nil {
			return nil, err
		} else if len(folders) == 0 {
			return nil, fmt.Errorf("no folder found matching '%v'", folder)
		}

		for k, v := range folders {
			path := prefix + "/" + v
			switch {
			case path == folder:
				folderID = k
				break loop

			case strings.HasPrefix(folder, path):
				folderID = k
				prefix = path
				continue loop
			}
		}

		return nil, fmt.Errorf("no folder found matching '%v'", folder)
	}

	files, err := files.ListFiles(folderID, b.token.Token)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (b *Box) DeleteFile(fileID string) error {
	return files.DeleteFile(fileID, b.token.Token)
}

func (b *Box) TagFile(fileID string, tag string) error {
	return files.TagFile(fileID, tag, b.token.Token)
}

func (b *Box) UntagFile(fileID string, tag string) error {
	return files.UntagFile(fileID, tag, b.token.Token)
}

func (b *Box) ListTemplates() (map[string]TemplateKey, error) {
	return listTemplates(b.token.Token)
}

func (b *Box) GetTemplate(key TemplateKey) (*Schema, error) {
	return getTemplate(key, b.token.Token)
}

func (b *Box) CreateTemplate(schema Schema) (interface{}, error) {
	return createTemplate(schema.Name, schema.Fields, b.token.Token)
}

func (b *Box) DeleteTemplate(key TemplateKey) error {
	return deleteTemplate(key, b.token.Token)
}
