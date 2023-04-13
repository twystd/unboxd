package commands

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/twystd/unboxd/box"
)

var UploadFileCmd = UploadFile{
	command: command{
		name:  "upload-file",
		delay: 500 * time.Millisecond,
	},
}

type UploadFile struct {
	command
}

func (cmd *UploadFile) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd UploadFile) Help() {
	fmt.Println()
	fmt.Println("  Usage: unboxd [--debug] --credentials <file> upload-file <file> <folder>")
	fmt.Println()
	fmt.Println("  Uploads a file to a Box folder.")
	fmt.Println()
	fmt.Println("    --credentials <file>  JSON file with Box credentials (required)")
	fmt.Println("      <file>              File to upload")
	fmt.Println("      <folder>            Destination folder")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    --debug  Enable debugging information")
	fmt.Println()
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    unboxd --debug --credentials .credentials upload-file photo.jpg /photos"`)
	fmt.Println()
}

func (cmd UploadFile) Execute(flagset *flag.FlagSet, b box.Box) error {
	args := flagset.Args()[1:]
	if len(args) < 1 {
		return fmt.Errorf("missing file argument")
	}

	if len(args) < 2 {
		return fmt.Errorf("missing folder argument")
	}

	file := args[0]
	folder := args[1]

	fileID, err := cmd.exec(b, file, folder)
	if err != nil {
		return err
	}

	log.Printf("%[1]v  %[2]v  %[3]v  uploaded", cmd.Name(), fileID, file)

	return nil
}

func (cmd UploadFile) exec(b box.Box, file string, folder string) (string, error) {
	if fileID, err := b.UploadFile(file, folder); err != nil {
		return "", err
	} else {
		return fileID, nil
	}
}
