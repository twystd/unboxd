package commands

import (
	"flag"
	"fmt"
	"log"

	"github.com/twystd/unboxd/box"
)

type UploadFile struct {
}

func (cmd UploadFile) Name() string {
	return "upload-file"
}

func (cmd *UploadFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd UploadFile) Help() {
}

func (cmd UploadFile) Execute(flagset *flag.FlagSet, b box.Box) error {
	args := flagset.Args()[1:]
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

	log.Printf("%[1]v  %[2]v  %[3]v  uploaded", cmd.Name(), fileID, file)

	return nil
}

func (cmd UploadFile) exec(b box.Box, file string, folder string) (string, error) {
	if fileID, err := b.UploadFile(file, folder); err != nil {
		return "", err
	} else {
		return fileID, nil
	}
}
