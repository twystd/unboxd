package commands

import (
	"flag"
	"fmt"
	"log"
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

func (cmd DeleteFile) Help() {
	fmt.Println()
	fmt.Println("  Usage: unboxd [--debug] --credentials <file> delete-file <file-id>")
	fmt.Println()
	fmt.Println("  Deletes a file stored in a Box folder.")
	fmt.Println()
	fmt.Println("    --credentials <file>  JSON file with Box credentials (required)")
	fmt.Println("      <file-id>           Box file ID")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    --debug  Enable debugging information")
	fmt.Println()
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    unboxd --debug --credentials .credentials delete-file 135789086421"`)
	fmt.Println()
}

func (cmd *DeleteFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
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
