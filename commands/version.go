package commands

import (
	"flag"
	"fmt"

	"github.com/twystd/unboxd/box"
)

type Version struct {
	Version string
}

func (cmd Version) Name() string {
	return "version"
}

func (cmd *Version) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd Version) Help() {
}

func (cmd Version) Execute(flagset *flag.FlagSet, b box.Box) error {
	fmt.Println()
	fmt.Printf("   unboxd %v\n", cmd.Version)
	fmt.Println()

	return nil
}
