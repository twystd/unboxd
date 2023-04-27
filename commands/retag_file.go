package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/credentials"
)

var RetagFileCmd = RetagFile{
	command: command{
		name:  "retag-file",
		delay: 500 * time.Millisecond,
	},
}

type RetagFile struct {
	command
}

func (cmd *RetagFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd RetagFile) Execute(c any, flagset *flag.FlagSet) error {
	b := box.NewBox()
	if err := b.Authenticate(c.(credentials.Credentials)); err != nil {
		return err
	}

	var fileID uint64
	var oldTag string
	var newTag string

	args := flagset.Args()

	if len(args) < 1 {
		return fmt.Errorf("missing file ID")
	} else if v, err := getFileID(args[0]); err != nil {
		return err
	} else {
		fileID = v
	}

	if len(args) < 2 {
		return fmt.Errorf("missing 'old' tag")
	} else {
		oldTag = args[1]
	}

	if len(args) < 3 {
		return fmt.Errorf("missing 'new' tag")
	} else {
		newTag = args[2]
	}

	if err := cmd.exec(b, fileID, oldTag, newTag); err != nil {
		return err
	}

	infof("retag-file", "%v replaced tag %v with %v\n", fileID, oldTag, newTag)

	return nil
}

func (cmd RetagFile) exec(b box.Box, fileID uint64, oldTag, newTag string) error {
	return b.RetagFile(fileID, oldTag, newTag)
}
