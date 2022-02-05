package lib

import (
	"strings"
)

type Glob interface {
	Match(string) bool
}

type glob struct {
	match func(string) bool
}

func NewGlob(g string) Glob {
	var match func(string) bool

	switch {
	case strings.HasSuffix(g, "/**"):
		l := len(g)
		match = func(p string) bool {
			return strings.HasPrefix(p, g[:l-2])
		}

	case strings.HasSuffix(g, "/*"):
		l := len(g)
		match = func(p string) bool {
			return strings.HasPrefix(p, g[:l-1]) && !strings.Contains(p[l:], "/")
		}

	case g == "/":
		l := len(g)
		match = func(p string) bool {
			return len(p) > 1 && strings.HasPrefix(p, "/") && !strings.Contains(p[l:], "/")
		}

	case strings.HasSuffix(g, "/"):
		l := len(g)
		match = func(p string) bool {
			return strings.HasPrefix(p, g) && !strings.Contains(p[l:], "/")
		}

	case g != "":
		match = func(p string) bool {
			return p == g
		}

	default:
		match = func(p string) bool {
			return true
		}
	}

	return glob{
		match: match,
	}
}

func (g glob) Match(s string) bool {
	return g.match(s)
}
