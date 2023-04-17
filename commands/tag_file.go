package commands

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/twystd/unboxd/box"
)

var TagFileCmd = TagFile{
	command: command{
		application: APP,
		name:        "tag-file",
		delay:       500 * time.Millisecond,
	},
}

type TagFile struct {
	command
}

func (cmd TagFile) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %v [--debug] --credentials <file> tag-file <file-id> <tag>\n", APP)
	fmt.Println()
	fmt.Println("  Adds a tag to a file stored in a Box folder.")
	fmt.Println()
	fmt.Println("    --credentials <file>  JSON file with Box credentials (required)")
	fmt.Println("      <file-id>           Box file ID")
	fmt.Println("      <tag>               Tag to add to file")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    --debug  Enable debugging information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Printf("    %v --debug --credentials .credentials tag-file 135789086421 hogwarts\n", APP)
	fmt.Println()
}

func (cmd *TagFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd TagFile) Execute(flagset *flag.FlagSet, b box.Box) error {
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

	log.Printf("%v  %v added tag %v\n", cmd.Name(), fileID, tag)

	return nil
}

func (cmd TagFile) exec(b box.Box, fileID uint64, tag string) error {
	return b.TagFile(fileID, tag)
}
