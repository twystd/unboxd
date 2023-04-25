package commands

import (
	"crypto/sha256"
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
}

type command struct {
	name  string
	delay time.Duration
}

func (cmd command) Name() string {
	return cmd.name
}

func (cmd command) hash(command string, credentials string, root string) string {
	s := fmt.Sprintf("%v:%v:%v", command, credentials, root)
	hash := sha256.Sum256([]byte(s))

	return fmt.Sprintf("%x", hash)
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

func infof(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-20v %v", tag, format)

	log.Infof(f, args...)
}

func warnf(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-20v %v", tag, format)

	log.Warnf(f, args...)
}
