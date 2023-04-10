package commands

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/twystd/unboxd/box"
)

var UntagFileCmd = TagFile{
	command: command{
		name:  "untag-file",
		delay: 500 * time.Millisecond,
	},
}

type UntagFile struct {
	command
}

func (cmd *UntagFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd UntagFile) Help() {
}

func (cmd UntagFile) Execute(flagset *flag.FlagSet, b box.Box) error {
	var fileID uint64
	var tag string

	args := flagset.Args()[1:]

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

	log.Printf("%v  %v removed tag %v\n", cmd.Name(), fileID, tag)

	return nil
}

func (cmd UntagFile) exec(b box.Box, fileID uint64, tag string) error {
	return b.UntagFile(fileID, tag)
}
