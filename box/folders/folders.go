package folders

import (
	"fmt"

	"github.com/twystd/unboxd/log"
)

const fetchSize = 500

type Folder struct {
	ID   uint64
	Name string
	Tags []string
}

func debugf(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-20v %v", tag, format)

	log.Debugf(f, args...)
}
