package box

import (
	"fmt"
	"strings"
)

type Box struct {
	token *AccessToken
}

type folders struct {
	photos  string
	pending string
}

type BoxFile struct {
	ID   FileID
	Name string
}

func NewBox() Box {
	return Box{}
}

const fetchSize = 128

func (b *Box) Authenticate(clientID, secret, user, userID string) error {
	if b.token != nil && b.token.IsValid() {
		return nil
	}

	token, err := authenticate(clientID, secret, user, userID)
	if err != nil {
		return err
	}

	b.token = token

	return nil
}

func (b *Box) ListFiles(folder string) ([]File, error) {
	prefix := ""
	folderID := FolderID("0")

loop:
	for {
		folders, err := listFolders(folderID, b.token.token)
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

	files, err := listFiles(folderID, b.token.token)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (b *Box) DeleteFile(fileID FileID) error {
	return deleteFile(fileID, b.token.token)
}

func (b *Box) ListTemplates() (map[string]TemplateKey, error) {
	return listTemplates(b.token.token)
}

func (b *Box) GetTemplate(key TemplateKey) (*Schema, error) {
	return getTemplate(key, b.token.token)
}

func (b *Box) CreateTemplate(schema Schema) (interface{}, error) {
	return createTemplate(schema.Name, schema.Fields, b.token.token)
}

func (b *Box) DeleteTemplate(key TemplateKey) error {
	return deleteTemplate(key, b.token.token)
}
