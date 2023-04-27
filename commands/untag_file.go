package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/credentials"
)

var UntagFileCmd = UntagFile{
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

func (cmd UntagFile) Execute(c any, flagset *flag.FlagSet) error {
	b := box.NewBox()
	if err := b.Authenticate(c.(credentials.Credentials)); err != nil {
		return err
	}

	var fileID uint64
	var tag string

	args := flagset.Args()

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

	infof("untag-file", "%v removed tag %v\n", fileID, tag)

	return nil
}

func (cmd UntagFile) exec(b box.Box, fileID uint64, tag string) error {
	return b.UntagFile(fileID, tag)
}
