package commands

import (
	"fmt"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/templates"
)

type ListTemplates struct {
}

func (cmd ListTemplates) Name() string {
	return "list-templates"
}

func (cmd ListTemplates) Execute(b box.Box) error {
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
