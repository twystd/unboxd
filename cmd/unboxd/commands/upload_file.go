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

func (cmd UploadFile) Execute(b box.Box) error {
	args := flag.Args()[1:]
	if len(args) < 1 {
		return fmt.Errorf("missing file argument")
	}

	if len(args) < 2 {
		return fmt.Errorf("missing folder argument")
	}

	file := args[0]
	folder := args[1]

	if err := cmd.exec(b, file, folder); err != nil {
		return err
	}

	log.Printf("%v  %v uploaded to %v\n", cmd.Name(), file, folder)

	return nil
}

func (cmd UploadFile) exec(b box.Box, file string, folder string) error {
	if _, err := b.UploadFile(file, folder); err != nil {
		return err
	}

	return nil
}
