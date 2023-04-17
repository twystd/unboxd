package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/templates"
)

var CreateTemplateCmd = CreateTemplate{
	command: command{
		application: APP,
		name:        "create-template",
		delay:       500 * time.Millisecond,
	},
}

type CreateTemplate struct {
	command
}

func (cmd CreateTemplate) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %v [--debug] --credentials <file> create-template  <template-file>\n", APP)
	fmt.Println()
	fmt.Println("  Creates a new Box metadata template from the definition in the template file.")
	fmt.Println()
	fmt.Println("    <template-file>  JSON metadata template definition")
	fmt.Println()
	fmt.Println("    --credentials <file>  JSON file with Box credentials (required)")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    --debug  Enable debugging information")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Printf("    %v --debug --credentials .credentials create-template hogwarts.json\n", APP)
	fmt.Println()
}

func (cmd *CreateTemplate) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (cmd CreateTemplate) Execute(flagset *flag.FlagSet, b box.Box) error {
	args := flagset.Args()[1:]
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
