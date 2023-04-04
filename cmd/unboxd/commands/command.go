package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/twystd/unboxd/box"
	"github.com/twystd/unboxd/log"
)

type Command interface {
	Name() string
	Execute(box.Box) error
}

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

func warnf(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-20v %v", tag, format)

	log.Warnf(f, args...)
}
