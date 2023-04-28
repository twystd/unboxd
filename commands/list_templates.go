package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/templates"
	"github.com/twystd/unboxd/credentials"
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

func (cmd *ListTemplates) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd ListTemplates) Execute(flagset *flag.FlagSet, c credentials.ICredentials) error {
	credentials := c["box"].(box.Credentials)

	b := box.NewBox()
	if err := b.Authenticate(credentials); err != nil {
		return err
	}

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
