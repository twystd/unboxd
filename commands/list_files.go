package commands

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/box/lib"
)

var ListFilesCmd = ListFiles{
	command: command{
		name:  "list-files",
		delay: 500 * time.Millisecond,
	},

	file:       "",
	checkpoint: ".checkpoint",
	tags:       false,
	restart:    false,
}

type ListFiles struct {
	command
	file       string
	checkpoint string
	tags       bool
	restart    bool
}

type file struct {
	ID       uint64
	FileName string
	FilePath string
	Tags     []string
}

func (cmd ListFiles) Help() {
	fmt.Println()
	fmt.Println("  Usage: unboxd [--debug] --credentials <file> list-files [--tags] [--file <file>] [--checkpoint <file>] [--delay <duration>] [--no-resume] <filespec>")
	fmt.Println()
	fmt.Println("  Retrieves a list of files that match the file spec")
	fmt.Println()
	fmt.Println("  A filespec is a glob expression against which to match file paths e.g.:")
	fmt.Println("    /*         matches files in the top level folder")
	fmt.Println("    /**        matches all files recursively")
	fmt.Println("    /photos/*  matches all files in the /photos folder")
	fmt.Println()
	fmt.Println("  The default filespec is /** i.e. list all files recursively")
	fmt.Println()
	fmt.Println("    --credentials <file>  JSON file with Box credentials (required)")
	fmt.Println("    --tags                Include tags in file information")
	fmt.Println("    --file                TSV file to which to write file information")
	fmt.Println("    --no-resume           Retrieves file list from the beginning (default is to continue from last checkpoint")
	fmt.Println("    --checkpoint          Specifies the path for the checkpoint file (default is .checkpoint)")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    --delay  Delay between multiple requests to reduce traffic to Box API")
	fmt.Println("    --debug  Enable debugging information")
	fmt.Println()
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    unboxd --debug --credentials .credentials list-files --tags --file folders.tsv /**"`)
	fmt.Println()
}

func (cmd *ListFiles) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.BoolVar(&cmd.tags, "tags", cmd.tags, "Include tags in folder information")
	flagset.StringVar(&cmd.file, "file", cmd.file, "TSV file to which to write folder information")
	flagset.StringVar(&cmd.checkpoint, "checkpoint", cmd.checkpoint, "Specifies the path for the checkpoint file")
	flag.DurationVar(&cmd.delay, "delay", cmd.delay, "Delay between multiple requests to reduce traffic to Box API")
	flagset.BoolVar(&cmd.restart, "no-resume", cmd.restart, "Retrieves folder list from the beginning")

	return flagset
}

func (cmd ListFiles) Execute(flagset *flag.FlagSet, b box.Box) error {
	args := flagset.Args()[1:]
	if len(args) < 1 {
		return fmt.Errorf("missing folder argument")
	}

	folder := args[0]

	files, err := cmd.exec(b, folder)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no files")
	}

	sort.Slice(files, func(i, j int) bool { return files[i].FileName < files[j].FileName })

	widths := []int{0, 0, 0}
	table := [][3]string{}

	for _, f := range files {
		id := fmt.Sprintf("%v", f.ID)
		filename := fmt.Sprintf("%v", f.FileName)
		tags := strings.Join(f.Tags, ",")

		if N := len(id); N > widths[0] {
			widths[0] = N
		}

		if N := len(filename); N > widths[1] {
			widths[1] = N
		}

		if N := len(tags); N > widths[2] {
			widths[2] = N
		}

		table = append(table, [3]string{id, filename, tags})
	}

	format := fmt.Sprintf("%%-%vv  %%-%vv  %%-%vv\n", widths[0], widths[1], widths[2])
	for _, row := range table {
		fmt.Printf(format, row[0], row[1], row[2])
	}

	return nil
}

func (cmd ListFiles) exec(b box.Box, glob string) ([]file, error) {
	folders, err := listFolders(b, 0, "", cmd.checkpoint, cmd.delay, cmd.restart)
	if err != nil {
		return nil, err
	}

	files := []file{}

	for _, f := range folders {
		l, err := listFiles(b, f.ID, f.Path)
		if err != nil {
			return nil, err
		}

		files = append(files, l...)
	}

	list := []file{}

	g := lib.NewGlob(glob + "/")
	for _, f := range files {
		if g.Match(f.FilePath) {
			list = append(list, f)
		}
	}

	return list, nil
}

func listFiles(b box.Box, folderID uint64, prefix string) ([]file, error) {
	files := []file{}

	l, err := b.ListFiles(folderID)
	if err != nil {
		return nil, err
	}

	for _, f := range l {
		path := prefix + "/" + f.Name
		files = append(files, file{
			ID:       f.ID,
			FileName: f.Name,
			FilePath: path,
			Tags:     f.Tags,
		})
	}

	return files, nil
}
