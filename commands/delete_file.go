package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/twystd/unboxd/box"
)

var DeleteFileCmd = DeleteFile{
	command: command{
		name:  "delete-file",
		delay: 500 * time.Millisecond,
	},
}

type DeleteFile struct {
	command
}

func (cmd *DeleteFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd DeleteFile) Execute(flagset *flag.FlagSet, b box.Box) error {
	args := flagset.Args()
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

		infof("delete-file", "%v deleted\n", file)
	}

	return nil
}

func (cmd DeleteFile) exec(b box.Box, fileID uint64) error {
	file := fmt.Sprintf("%v", fileID)

	return b.DeleteFile(file)
}
