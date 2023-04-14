package main

import (
	"flag"
	"fmt"

	"github.com/twystd/unboxd/box"
)

type Help struct {
}

func (cmd Help) Name() string {
	return "help"
}

func (h *Help) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (h Help) Help() {
	fmt.Println()
	fmt.Println("  Usage: unboxd help [command]")
	fmt.Println()
	fmt.Println("  Displays the command line and a list of the available commands. If a command is")
	fmt.Println("  specified, displays the help information for that command.")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    unboxd help")
	fmt.Println("    unboxd help list-folders")
	fmt.Println()
}

func (h Help) Execute(flagset *flag.FlagSet, box box.Box) error {
	for _, c := range cli {
		if c.Name() == flagset.Arg(0) {
			c.Help()
			return nil
		}
	}

	h.help()

	return nil
}

func (h Help) help() {
	fmt.Println()
	fmt.Println("  Usage: unboxd <command> <options>")
	fmt.Println()
	fmt.Println("  Supported commands:")

	for _, c := range cli {
		fmt.Printf("    %v\n", c.Name())
	}

	fmt.Println()
	fmt.Println("  Use 'unboxd help <command>' for command specific information.")
	fmt.Println()
}
