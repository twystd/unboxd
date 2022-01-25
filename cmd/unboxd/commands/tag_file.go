package commands

import (
	"flag"
	"fmt"
	"log"

	"github.com/twystd/unboxd/box"
)

type TagFile struct {
}

func (cmd TagFile) Name() string {
	return "tag-file"
}

func (cmd TagFile) Execute(b box.Box) error {
	var fileID uint64
	var tag string

	args := flag.Args()[1:]

	if len(args) < 1 {
		return fmt.Errorf("missing file ID")
	} else if v, err := getFileID(args[0]); err != nil {
		return err
	} else {
		fileID = v
	}

	if len(args) < 2 {
		return fmt.Errorf("missing tag")
	} else {
		tag = args[1]
	}

	if err := cmd.exec(b, fileID, tag); err != nil {
		return err
	}

	log.Printf("%v  %v tagged\n", cmd.Name(), fileID)

	return nil
}

func (cmd TagFile) exec(b box.Box, fileID uint64, tag string) error {
	file := box.FileID(fmt.Sprintf("%v", fileID))

	return b.TagFile(file, tag)
}
