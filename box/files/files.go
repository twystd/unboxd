package files

type File struct {
	ID       string
	Filename string
	Tags     []string
}

const fetchSize = 128
