package commands

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/twystd/unboxd/box"
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

func (cmd RetagFile) Help() {
	fmt.Println()
	fmt.Println("  Usage: unboxd [--debug] --credentials <file> retag-file <file-id> <old-tag> <new-tag>")
	fmt.Println()
	fmt.Println("  Replaces a tag on a file stored in a Box folder. The tag is only replaced if it exists.")
	fmt.Println()
	fmt.Println("    --credentials <file>  JSON file with Box credentials (required)")
	fmt.Println("      <file-id>           Box file ID")
	fmt.Println("      <old-tag>           Tag to be replaced")
	fmt.Println("      <new-tag>           Replacement tag")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    --debug  Enable debugging information")
	fmt.Println()
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    unboxd --debug --credentials .credentials retag-file 135789086421 hogwarts hogsmeade"`)
	fmt.Println()
}

func (cmd *RetagFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
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
