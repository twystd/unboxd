package commands

import (
	"flag"
	"fmt"

	"github.com/twystd/unboxd/credentials"
)

type Version struct {
	APP     string
	Version string
}

func (cmd Version) Name() string {
	return "version"
}

func (cmd *Version) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd Version) Execute(flagset *flag.FlagSet, c credentials.ICredentials) error {
	fmt.Println()
	fmt.Printf("   %v %v\n", cmd.APP, cmd.Version)
	fmt.Println()

	return nil
}
