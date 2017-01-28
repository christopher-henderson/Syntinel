package statistics

import "bufio"

type Top interface {
	Command() string
	Args() []string
	Parse(*bufio.Scanner)
}
