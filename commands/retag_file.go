package commands

import (
	"flag"
	"fmt"
	"log"

	"github.com/twystd/unboxd/box"
)

type RetagFile struct {
}

func (cmd RetagFile) Name() string {
	return "retag-file"
}

func (cmd *RetagFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd RetagFile) Help() {
}

func (cmd RetagFile) Execute(flagset *flag.FlagSet, b box.Box) error {
	var fileID uint64
	var oldTag string
	var newTag string

	args := flagset.Args()[1:]

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

	log.Printf("%v  %v replaced tag %v with %v\n", cmd.Name(), fileID, oldTag, newTag)

	return nil
}

func (cmd RetagFile) exec(b box.Box, fileID uint64, oldTag, newTag string) error {
	return b.RetagFile(fileID, oldTag, newTag)
}
