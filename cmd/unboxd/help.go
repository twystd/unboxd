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
	fmt.Printf("  Usage: %v help [command]\n", APP)
	fmt.Println()
	fmt.Println("  Displays the command line and a list of the available commands. If a command is")
	fmt.Println("  specified, displays the help information for that command.")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Printf("    %v help\n", APP)
	fmt.Printf("    %v help list-folders\n", APP)
	fmt.Println()
}

func (h Help) Execute(flagset *flag.FlagSet, box box.Box) error {
	for _, c := range CLI {
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
	fmt.Printf("  Usage: %v <command> <options>\n", APP)
	fmt.Println()
	fmt.Println("  Supported commands:")

	for _, c := range CLI {
		fmt.Printf("    %v\n", c.Name())
	}

	fmt.Println()
	fmt.Printf("  Use '%v help <command>' for command specific information.\n", APP)
	fmt.Println()
}
