package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/credentials"
)

var UploadFileCmd = UploadFile{
	command: command{
		name:  "upload-file",
		delay: 500 * time.Millisecond,
	},
}

type UploadFile struct {
	command
}

func (cmd *UploadFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd UploadFile) Execute(flagset *flag.FlagSet, c credentials.ICredentials) error {
	credentials := c["box"].(box.Credentials)

	b := box.NewBox()
	if err := b.Authenticate(credentials); err != nil {
		return err
	}

	args := flagset.Args()
	if len(args) < 1 {
		return fmt.Errorf("missing file argument")
	}

	if len(args) < 2 {
		return fmt.Errorf("missing folder argument")
	}

	file := args[0]
	folder := args[1]

	fileID, err := cmd.exec(b, file, folder)
	if err != nil {
		return err
	}

	infof("upload-file", "%v  %v  uploaded", fileID, file)

	return nil
}

func (cmd UploadFile) exec(b box.Box, file string, folder string) (string, error) {
	if fileID, err := b.UploadFile(file, folder); err != nil {
		return "", err
	} else {
		return fileID, nil
	}
}
