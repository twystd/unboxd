package box

import (
	"github.com/twystd/unboxd/box/credentials"
	"github.com/twystd/unboxd/box/files"
	"github.com/twystd/unboxd/box/folders"
	"github.com/twystd/unboxd/box/templates"
)

type Box struct {
	token *credentials.AccessToken
	hash  string
}

func NewBox() Box {
	return Box{}
}

func (b *Box) Authenticate(credentials credentials.Credentials) error {
	if b.token != nil && b.token.IsValid() {
		return nil
	}

	token, err := credentials.Authenticate()
	if err != nil {
		return err
	}

	b.token = token
	b.hash = credentials.Hash()

	return nil
}

func (b Box) Hash() string {
	return b.hash
}

func (b *Box) ListFolders(folderID uint64) ([]folders.Folder, error) {
	return folders.List(folderID, b.token.Token)
}

func (b *Box) ListFiles(folderID uint64) ([]files.File, error) {
	return files.List(folderID, b.token.Token)
}

func (b *Box) UploadFile(file string, folder string) (string, error) {
	return files.Upload(file, folder, b.token.Token)
}

func (b *Box) DeleteFile(fileID string) error {
	return files.Delete(fileID, b.token.Token)
}

func (b *Box) TagFile(fileID uint64, tag string) error {
	return files.Tag(fileID, tag, b.token.Token)
}

func (b *Box) UntagFile(fileID uint64, tag string) error {
	return files.Untag(fileID, tag, b.token.Token)
}

func (b *Box) RetagFile(fileID uint64, oldTag, newTag string) error {
	return files.Retag(fileID, oldTag, newTag, b.token.Token)
}

func (b *Box) ListTemplates() (map[string]templates.TemplateKey, error) {
	return templates.List(b.token.Token)
}

func (b *Box) GetTemplate(key templates.TemplateKey) (*templates.Schema, error) {
	return templates.Get(key, b.token.Token)
}

func (b *Box) CreateTemplate(schema templates.Schema) (interface{}, error) {
	return templates.Create(schema.Name, schema.Fields, b.token.Token)
}

func (b *Box) DeleteTemplate(key templates.TemplateKey) error {
	return templates.Delete(key, b.token.Token)
}
