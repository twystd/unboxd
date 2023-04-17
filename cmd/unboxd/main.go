package main

import (
	"github.com/twystd/unboxd/commands"
)

var APP = commands.APP
var VERSION = commands.VERSION

var CLI = []commands.Command{
	&commands.ListFoldersCmd,

	&commands.ListFilesCmd,
	&commands.UploadFileCmd,
	&commands.DeleteFileCmd,
	&commands.TagFileCmd,
	&commands.UntagFileCmd,
	&commands.RetagFileCmd,

	&commands.ListTemplatesCmd,
	&commands.GetTemplateCmd,
	&commands.CreateTemplateCmd,
	&commands.DeleteTemplateCmd,

	&commands.VersionCmd,
	&Help{},
}

func main() {
	exec(CLI)
	// // ... parse command line
	// cmd, flagset, err := parse()
	// if err != nil {
	// 	fmt.Printf("ERROR: unable to parse command line (%v)\n", err)
	// 	return
	// }

	// if cmd == nil {
	// 	usage()
	// 	os.Exit(1)
	// }

	// if cmd.Name() == "help" {
	// 	cmd.Execute(flagset, box.Box{})
	// 	os.Exit(0)
	// }

	// if cmd.Name() == "version" {
	// 	cmd.Execute(flagset, box.Box{})
	// 	os.Exit(0)
	// }

	// if options.debug {
	// 	log.SetLevel("debug")
	// }

	// credentials, err := NewCredentials(options.credentials)
	// if err != nil {
	// 	log.Fatalf("Error reading credentials from %s (%v)", options.credentials, err)
	// }

	// box := box.NewBox()
	// if err := box.Authenticate(credentials); err != nil {
	// 	log.Fatalf("%v", err)
	// } else if err := cmd.Execute(flagset, box); err != nil {
	// 	log.Fatalf("%v  %v", cmd.Name(), err)
	// }
}
