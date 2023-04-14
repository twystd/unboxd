package commands

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/templates"
)

var DeleteTemplateCmd = DeleteTemplate{
	command: command{
		name:  "delete-template",
		delay: 500 * time.Millisecond,
	},
}

type DeleteTemplate struct {
	command
}

func (cmd DeleteTemplate) Help() {
	fmt.Println()
	fmt.Println("  Usage: unboxd [--debug] --credentials <file> delete-template [--exact] [--key] [--out <file>] <template-id>")
	fmt.Println()
	fmt.Println("  Deletes a Box metadata template associated with the account.")
	fmt.Println()
	fmt.Println("    <template-id>  Metadata template name or Box ID")
	fmt.Println()
	fmt.Println("    --credentials <file>  JSON file with Box credentials (required)")
	fmt.Println("    --exact               Requires that a template name match the template ID exactly (defaults to 'approximately')")
	fmt.Println("    --key                 Requires that the template Box ID match the template ID")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    --debug  Enable debugging information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    unboxd --debug --credentials .credentials delete-template --exact HOGWARTS")
	fmt.Println()
}

func (cmd *DeleteTemplate) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd DeleteTemplate) Execute(flagset *flag.FlagSet, b box.Box) error {
	template := ""
	exactMatch := false
	byKey := false

	args := flagset.Args()[1:]
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

	list, err := b.ListTemplates()
	if err != nil {
		return err
	} else if list == nil {
		return fmt.Errorf("invalid template list")
	}

	keys := []templates.TemplateKey{}
	for k, v := range list {
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

func (cmd DeleteTemplate) exec(b box.Box, t templates.TemplateKey) error {
	return b.DeleteTemplate(t)
}
