package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/commands"
	"github.com/twystd/unboxd/log"
)

var options = struct {
	credentials string
	debug       bool
}{
	credentials: ".credentials.json",
	debug:       false,
}

func exec(cli []commands.Command) {
	// ... parse command line
	cmd, flagset, err := parse(cli)
	if err != nil {
		fmt.Printf("ERROR: unable to parse command line (%v)\n", err)
		return
	}

	if cmd == nil {
		usage(cli)
		os.Exit(1)
	}

	if cmd.Name() == "help" {
		cmd.Execute(flagset, box.Box{})
		os.Exit(0)
	}

	if cmd.Name() == "version" {
		cmd.Execute(flagset, box.Box{})
		os.Exit(0)
	}

	if options.debug {
		log.SetLevel("debug")
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

func usage(cli []commands.Command) {
	fmt.Println()
	fmt.Printf("   Usage: %v [--debug] --credentials <file> <command>\n", APP)
	fmt.Println()
	fmt.Println("   Commands:")
	fmt.Println()

	for _, c := range cli {
		fmt.Printf("     %v\n", c.Name())
	}

	fmt.Println()
}

func parse(cli []commands.Command) (commands.Command, *flag.FlagSet, error) {
	flagset := flag.NewFlagSet(APP, flag.ExitOnError)

	flagset.StringVar(&options.credentials, "credentials", options.credentials, "(required) JSON file with Box credentials")
	flagset.BoolVar(&options.debug, "debug", options.debug, "(optional) Enable debugging information")
	flagset.Parse(os.Args[1:])

	args := flagset.Args()

	if len(args) > 0 {
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
