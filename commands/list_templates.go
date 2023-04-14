package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/templates"
)

var ListTemplatesCmd = ListTemplates{
	command: command{
		name:  "list-templates",
		delay: 500 * time.Millisecond,
	},
}

type ListTemplates struct {
	command
}

func (cmd ListTemplates) Help() {
	fmt.Println()
	fmt.Println("  Usage: unboxd [--debug] --credentials <file> list-templates")
	fmt.Println()
	fmt.Println("  Retrieves the full list of metadata templates associated with the account.")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    --debug  Enable debugging information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    unboxd --debug --credentials .credentials list-templates")
	fmt.Println()
}

func (cmd *ListTemplates) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd ListTemplates) Execute(flagset *flag.FlagSet, b box.Box) error {
	if templates, err := cmd.exec(b); err != nil {
		return err
	} else if len(templates) == 0 {
		return fmt.Errorf("no templates defined")
	} else {
		for k, v := range templates {
			fmt.Printf("%-16s  %s\n", k, v)
		}
	}

	return nil
}

func (cmd ListTemplates) exec(b box.Box) (map[string]templates.TemplateKey, error) {
	return b.ListTemplates()
}
