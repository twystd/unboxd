package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/log"
)

type Command interface {
	Name() string
	Flagset(*flag.FlagSet) *flag.FlagSet
	Execute(*flag.FlagSet, box.Box) error
	Help()
}

type command struct {
	name  string
	delay time.Duration
}

func (cmd command) Name() string {
	return cmd.name
}

// func helpOptions(flagset *flag.FlagSet) {
// 	count := 0
// 	flag.VisitAll(func(f *flag.Flag) {
// 		count++
// 	})
//
// 	flagset.VisitAll(func(f *flag.Flag) {
// 		fmt.Printf("    --%-19s %s\n", f.Name, f.Usage)
// 	})
//
// 	fmt.Println()
// 	fmt.Println("  Options:")
// 	flag.VisitAll(func(f *flag.Flag) {
// 		fmt.Printf("    --%-6s %s\n", f.Name, f.Usage)
// 	})
//
// 	fmt.Printf("    --%-6s %s\n", "debug", "Enable debugging information")
// }

func clean(s string) string {
	return regexp.MustCompile(`[\s\t]+`).ReplaceAllString(strings.ToLower(s), "")
}

func getFileID(arg string) (uint64, error) {
	if !regexp.MustCompile("^[0-9]+$").MatchString(arg) {
		return 0, fmt.Errorf("invalid file ID")
	} else if v, err := strconv.ParseUint(arg, 10, 64); err != nil {
		return 0, fmt.Errorf("invalid file ID %v)", err)
	} else {
		return uint64(v), nil
	}
}

func infof(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-20v %v", tag, format)

	log.Infof(f, args...)
}

func warnf(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-20v %v", tag, format)

	log.Warnf(f, args...)
}
