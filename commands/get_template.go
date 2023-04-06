package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/templates"
)

type GetTemplate struct {
}

func (cmd GetTemplate) Name() string {
	return "get-template"
}

func (cmd *GetTemplate) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd GetTemplate) Help() {
}

func (cmd GetTemplate) Execute(flagset *flag.FlagSet, b box.Box) error {
	template := ""
	file := ""
	exactMatch := false
	byKey := false

	args := flagset.Args()[1:]
	if len(args) == 1 {
		template = args[0]
		args = args[1:]
	} else if len(args) > 1 {
		arg := args[0]
		args = args[1:]

		switch arg {
		case "--exact":
			exactMatch = true
			arg := args[0]
			args = args[1:]
			template = arg

		case "--key":
			byKey = true
			arg := args[0]
			args = args[1:]
			template = arg

		default:
			template = arg
		}
	}

	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		switch arg {
		case "--out":
			if len(args) > 0 {
				arg := args[0]
				args = args[1:]
				file = arg
			}
		}
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
		if schema, err := cmd.exec(b, keys[0]); err != nil {
			return err
		} else if file != "" {
			if err := cmd.save(schema, file); err != nil {
				return err
			}

			log.Printf("saved template %v to file %v\n", template, file)
			return nil

		} else {
			return cmd.print(schema)
		}

	default:
		return fmt.Errorf("more than one template found matching %v (specify --exact for exact match, or --key to get by template key)", template)
	}
}

func (cmd GetTemplate) exec(b box.Box, t templates.TemplateKey) (*templates.Schema, error) {
	if schema, err := b.GetTemplate(t); err != nil {
		return nil, err
	} else if schema == nil {
		return nil, fmt.Errorf("invalid schema")
	} else {
		return schema, nil
	}
}

func (cmd GetTemplate) save(schema *templates.Schema, file string) error {
	if bytes, err := json.MarshalIndent(schema, "  ", "  "); err != nil {
		return err
	} else {
		return os.WriteFile(file, bytes, 0666)
	}
}

func (cmd GetTemplate) print(schema *templates.Schema) error {
	if bytes, err := json.MarshalIndent(schema, "  ", "  "); err != nil {
		return err
	} else {
		fmt.Printf("%v\n", string(bytes))
	}

	return nil
}
