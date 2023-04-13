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
	batch:      0,
}

type ListFiles struct {
	command
	file       string
	checkpoint string
	tags       bool
	restart    bool
	batch      uint
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
	fmt.Println("  Retrieves a list of files that match the file spec.")
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
	fmt.Println("    --batch               Maximum number of calls to the Box API (defaults to no limit)")
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
	flagset.UintVar(&cmd.batch, "batch-size", cmd.batch, "Number of calls to the Box API")

	return flagset
}

func (cmd ListFiles) Execute(flagset *flag.FlagSet, b box.Box) error {
	glob := ""

	args := flagset.Args()
	if len(args) > 0 {
		glob = args[0]
	}

	files, err := cmd.exec(b, glob)
	if err != nil {
		return err
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
	list := []file{}

	folders, err := listFiles(b, 0, "", cmd.checkpoint, cmd.delay, cmd.restart, cmd.batch)
	if err != nil {
		return nil, err
	}

	g := lib.NewGlob(glob)
	for _, f := range folders {
		if g.Match(f.FilePath) {
			list = append(list, f)
		}
	}

	return list, nil
}

func listFiles(b box.Box, folderID uint64, prefix string, chkpt string, delay time.Duration, restart bool, batch uint) ([]file, error) {
	pipe, folders, files, err := resume(chkpt, restart)
	if err != nil {
		return nil, err
	}

	if len(pipe) > 0 {
		infof("list-files", "Resuming last operation")
	} else {
		pipe = append(pipe, QueueItem{ID: folderID, Path: prefix})
	}

	count := uint(0)
	tail := 0
	for tail < len(pipe) {
		time.Sleep(delay)

		item := pipe[tail]

		// get files for current folder
		if l, err := b.ListFiles(item.ID); err != nil {
			if errx := checkpoint(chkpt, pipe[tail:], folders, files); errx != nil {
				warnf("list-files", "%v", errx)
			}

			return files, err
		} else {
			for _, f := range l {
				path := prefix + "/" + f.Name
				files = append(files, file{
					ID:       f.ID,
					FileName: f.Name,
					FilePath: path,
					Tags:     f.Tags,
				})
			}
		}

		// get subfolders for current folder
		if l, err := b.ListFolders(item.ID); err != nil {
			if errx := checkpoint(chkpt, pipe[tail:], folders, files); errx != nil {
				warnf("list-files", "%v", errx)
			}

			return files, err
		} else {
			for _, f := range l {
				path := item.Path + "/" + f.Name
				folders = append(folders, folder{
					ID:   f.ID,
					Name: f.Name,
					Tags: f.Tags,
					Path: path,
				})

				pipe = append(pipe, QueueItem{ID: f.ID, Path: path})
			}
		}

		count++
		if batch != 0 && count > batch {
			break
		}

		tail++
	}

	// ... incomplete?
	if len(pipe[tail:]) > 0 {
		if err := checkpoint(chkpt, pipe[tail:], folders, files); err != nil {
			return files, err
		} else {
			return files, fmt.Errorf("interrupted")
		}
	}

	// ... complete!
	if err := checkpoint(chkpt, []QueueItem{}, []folder{}, []file{}); err != nil {
		return files, err
	}

	return files, nil
}
