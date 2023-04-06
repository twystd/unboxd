package commands

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/twystd/unboxd/box"
)

type DeleteFile struct {
}

func (cmd DeleteFile) Name() string {
	return "delete-file"
}

func (cmd *DeleteFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd DeleteFile) Help() {
}

func (cmd DeleteFile) Execute(flagset *flag.FlagSet, b box.Box) error {
	args := flagset.Args()[1:]
	if len(args) == 0 {
		return fmt.Errorf("missing file ID argument")
	}

	for _, file := range args {
		var fileID uint64
		if !regexp.MustCompile("^[0-9]+$").MatchString(file) {
			return fmt.Errorf("invalid file ID")
		} else if v, err := strconv.ParseUint(file, 10, 64); err != nil {
			return fmt.Errorf("invalid file ID %v)", err)
		} else {
			fileID = uint64(v)
		}

		if err := cmd.exec(b, fileID); err != nil {
			return err
		}

		log.Printf("%v  %v deleted\n", cmd.Name(), file)
	}

	return nil
}

func (cmd DeleteFile) exec(b box.Box, fileID uint64) error {
	file := fmt.Sprintf("%v", fileID)

	return b.DeleteFile(file)
}
