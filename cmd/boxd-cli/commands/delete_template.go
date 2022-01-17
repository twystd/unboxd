package commands

import (
	"flag"
	"fmt"
	"log"

	"github.com/twystd/boxd/box"
)

type DeleteTemplate struct {
}

func (cmd DeleteTemplate) Name() string {
	return "delete-template"
}

func (cmd DeleteTemplate) Execute(b box.Box) error {
	template := ""
	exactMatch := false
	byKey := false

	args := flag.Args()[1:]
	if len(args) == 1 {
		template = args[0]
	} else if len(args) > 1 {
		switch args[0] {
		case "--exact":
			exactMatch = true

		case "--key":
			byKey = true
		}

		template = args[1]
	}

	if template == "" {
		return fmt.Errorf("missing template name argument")
	}

	templates, err := b.ListTemplates()
	if err != nil {
		return err
	} else if templates == nil {
		return fmt.Errorf("invalid template list")
	}

	keys := []box.TemplateKey{}
	for k, v := range templates {
		switch {
		case exactMatch:
			if template == k {
				keys = append(keys, v)
			}

		case byKey:
			if fmt.Sprintf("%v", v) == template {
				keys = append(keys, v)
			}

		default:
			if clean(template) == clean(k) {
				keys = append(keys, v)
			}
		}
	}

	switch len(keys) {
	case 0:
		if byKey {
			return fmt.Errorf("no template found for key %v", template)
		} else {
			return fmt.Errorf("no template found for name %v", template)
		}

	case 1:
		if err := cmd.exec(b, keys[0]); err != nil {
			return err
		} else {
			log.Printf("%v  %v deleted\n", cmd.Name(), template)
			return nil
		}

	default:
		return fmt.Errorf("more than one template found matching %v (specify --exact for exact match, or --key to delete by template key)", template)
	}
}

func (cmd DeleteTemplate) exec(b box.Box, t box.TemplateKey) error {
	return b.DeleteTemplate(t)
}
