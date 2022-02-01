package box

import (
	"strings"
)

type Glob struct {
	Match func(string) bool
}

func NewGlob(glob string) Glob {
	var match func(string) bool

	switch {
	case strings.HasSuffix(glob, "/**"):
		l := len(glob)
		match = func(p string) bool {
			return strings.HasPrefix(p, glob[:l-2])
		}

	case strings.HasSuffix(glob, "/*"):
		l := len(glob)
		match = func(p string) bool {
			return strings.HasPrefix(p, glob[:l-1]) && !strings.Contains(p[l:], "/")
		}

	case glob == "/":
		l := len(glob)
		match = func(p string) bool {
			return len(p) > 1 && strings.HasPrefix(p, "/") && !strings.Contains(p[l:], "/")
		}

	case strings.HasSuffix(glob, "/"):
		l := len(glob)
		match = func(p string) bool {
			return strings.HasPrefix(p, glob) && !strings.Contains(p[l:], "/")
		}

	case glob != "":
		match = func(p string) bool {
			return p == glob
		}

	default:
		match = func(p string) bool {
			return true
		}
	}

	return Glob{
		Match: match,
	}
}
