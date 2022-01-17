package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/twystd/boxd/box"
	"github.com/twystd/boxd/cmd/boxd-cli/commands"
)

var options = struct {
	credentials string
	debug       bool
}{
	credentials: ".credentials.json",
	debug:       false,
}

var cli = []commands.Command{
	commands.ListFiles{},
	commands.DeleteFile{},

	commands.ListTemplates{},
	commands.GetTemplate{},
	commands.CreateTemplate{},
	commands.DeleteTemplate{},
}

func main() {
	flag.StringVar(&options.credentials, "credentials", options.credentials, "(required) JSON file with Box credentials")
	flag.BoolVar(&options.debug, "debug", options.debug, "(optional) enables debug mode")
	flag.Parse()

	cmd := "help"
	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	} else {
		cmd = args[0]
	}

	if cmd == "help" {
		usage()
		os.Exit(0)
	}

	credentials := Credentials{}
	if err := credentials.Load(options.credentials); err != nil {
		log.Fatalf("Error reading credentials from %s (%v)", options.credentials, err)
	}

	box := box.NewBox()

	var f commands.Command
	for _, c := range cli {
		if cmd == c.Name() {
			f = c
			break
		}
	}

	if f == nil {
		usage()
		os.Exit(1)
	} else if err := box.Authenticate(credentials.ClientID, credentials.Secret, credentials.User, credentials.UserID); err != nil {
		log.Fatalf("%v", err)
	} else if err := f.Execute(box); err != nil {
		log.Fatalf("%v  %v", f.Name(), err)
	}
}

func usage() {
	fmt.Println()
	fmt.Println("   Usage: boxd-cli [--debug] --credentials <file> <command>")
	fmt.Println()
	fmt.Println("   Commands:")
	fmt.Println()

	for _, c := range cli {
		fmt.Printf("     %v\n", c.Name())
	}
}
