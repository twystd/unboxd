package commands

import (
	"flag"
	"fmt"

	"github.com/twystd/unboxd/box"
)

var VersionCmd = Version{
	command: command{
		application: APP,
		name:        "version",
	},
	Version: VERSION,
}

type Version struct {
	command
	Version string
}

func (cmd *Version) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd Version) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %v version\n", APP)
	fmt.Println()
	fmt.Println("  Displays the current version information.")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Printf("    %v version\n", APP)
	fmt.Println()
}

func (cmd Version) Execute(flagset *flag.FlagSet, b box.Box) error {
	fmt.Println()
	fmt.Printf("   %v %v\n", APP, cmd.Version)
	fmt.Println()

	return nil
}
