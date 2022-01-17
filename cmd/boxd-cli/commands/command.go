package commands

import (
	"regexp"
	"strings"

	"github.com/twystd/boxd/box"
)

type Command interface {
	Name() string
	Execute(box.Box) error
}

func clean(s string) string {
	return regexp.MustCompile(`[\s\t]+`).ReplaceAllString(strings.ToLower(s), "")
}
