package main

import (
	"flag"
	"fmt"
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
}

func (h Help) Execute(flagset *flag.FlagSet) error {
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
