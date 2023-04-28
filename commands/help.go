package commands

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/twystd/unboxd/credentials"
)

//go:embed help.txt
var HelpText string

type Help struct {
	APP      string
	CLI      []Command
	HelpText string
}

func (cmd Help) Name() string {
	return "help"
}

func (h *Help) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (h Help) Execute(flagset *flag.FlagSet, c credentials.ICredentials) error {
	command := flagset.Arg(0)
	info := map[string]string{
		"APP": h.APP,
	}

	helptext := HelpText
	if h.HelpText != "" {
		helptext = h.HelpText
	}

	templates := template.Must(template.New("help").Parse(helptext))
	if t := templates.Lookup(command); t == nil {
		h.help()
	} else if err := t.Execute(os.Stdout, info); err != nil {
		h.help()
	}

	return nil
}

func (h Help) help() {
	fmt.Println()
	fmt.Printf("  Usage: %v <command> <options>\n", h.APP)
	fmt.Println()
	fmt.Println("  Supported commands:")

	for _, c := range h.CLI {
		fmt.Printf("    %v\n", c.Name())
	}

	fmt.Println()
	fmt.Printf("  Use '%v help <command>' for command specific information.\n", h.APP)
	fmt.Println()
}
