package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/commands"
)

var VERSION = "v0.0.x"

var options = struct {
	credentials string
	debug       bool
}{
	credentials: ".credentials.json",
	debug:       false,
}

var cli = []commands.Command{
	&commands.ListFoldersCmd,

	&commands.ListFilesCmd,
	&commands.UploadFile{},
	&commands.DeleteFile{},
	&commands.TagFile{},
	&commands.UntagFile{},
	&commands.RetagFile{},

	&commands.ListTemplates{},
	&commands.GetTemplate{},
	&commands.CreateTemplate{},
	&commands.DeleteTemplate{},
}

func main() {
	// ... parse command line
	cmd, flagset, err := parse()
	if err != nil {
		fmt.Printf("ERROR: unable to parse command line (%v)\n", err)
		return
	}

	// if cmd == "help" {
	// 	usage()
	// 	os.Exit(0)
	// }
	//
	// if cmd == "version" {
	// 	version()
	// 	os.Exit(0)
	// }

	if cmd == nil {
		usage()
		os.Exit(1)
	}

	credentials, err := NewCredentials(options.credentials)
	if err != nil {
		log.Fatalf("Error reading credentials from %s (%v)", options.credentials, err)
	}

	box := box.NewBox()
	if err := box.Authenticate(credentials); err != nil {
		log.Fatalf("%v", err)
	} else if err := cmd.Execute(flagset, box); err != nil {
		log.Fatalf("%v  %v", cmd.Name(), err)
	}
}

func usage() {
	fmt.Println()
	fmt.Println("   Usage: unboxd [--debug] --credentials <file> <command>")
	fmt.Println()
	fmt.Println("   Commands:")
	fmt.Println()

	for _, c := range cli {
		fmt.Printf("     %v\n", c.Name())
	}

	fmt.Println()
}

func version() {
	fmt.Println()
	fmt.Printf("   boxd-cli %v\n", VERSION)
	fmt.Println()
}

func parse() (commands.Command, *flag.FlagSet, error) {
	flagset := flag.NewFlagSet("unboxd", flag.ExitOnError)

	flagset.StringVar(&options.credentials, "credentials", options.credentials, "(required) JSON file with Box credentials")
	flagset.BoolVar(&options.debug, "debug", options.debug, "(optional) enable debugging information")
	flagset.Parse(os.Args[1:])

	args := flagset.Args()
	if len(args) > 1 {
		for _, c := range cli {
			if c.Name() == args[0] {
				cmd := c
				flagset = cmd.Flagset(flagset)
				if err := flagset.Parse(args[1:]); err != nil {
					return cmd, flagset, err
				} else {
					return cmd, flagset, nil
				}
			}
		}
	}

	return nil, flagset, nil
}
