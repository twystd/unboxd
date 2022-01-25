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

type CreateTemplate struct {
}

func (cmd CreateTemplate) Name() string {
	return "create-template"
}

func (cmd CreateTemplate) Execute(b box.Box) error {
	args := flag.Args()[1:]
	if len(args) < 1 {
		return fmt.Errorf("missing template JSON file argument")
	}

	file := args[0]
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	schema := templates.Schema{}
	if err := json.Unmarshal(bytes, &schema); err != nil {
		return err
	}

	if err := cmd.exec(b, schema); err != nil {
		return err
	}

	log.Printf("%v  %v created\n", cmd.Name(), schema.Name)

	return nil
}

func (cmd CreateTemplate) exec(b box.Box, schema templates.Schema) error {
	if _, err := b.CreateTemplate(schema); err != nil {
		return err
	}

	return nil
}
