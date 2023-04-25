package main

import (
	"github.com/twystd/unboxd/commands"
)

var APP = "unboxd"
var VERSION = "v0.0.x"

var version = commands.Version{
	APP:     APP,
	Version: VERSION,
}

var help = commands.Help{
	APP: APP,
}

var cli = []commands.Command{
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

	&version,
	&help,
}

func init() {
	help.CLI = cli
}

func main() {
	exec(cli)
}
