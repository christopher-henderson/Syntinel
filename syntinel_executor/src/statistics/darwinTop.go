package statistics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type DarwinTop struct {
	//13
	command string
	args    []string
}

func NewDarwinTop() DarwinTop {
	return DarwinTop{"top", []string{"-pid", strconv.Itoa(os.Getppid())}}
}

func (dt DarwinTop) Command() string {
	return dt.command
}

func (dt DarwinTop) Args() []string {
	return dt.args
}

func (dt DarwinTop) Parse(out *bufio.Scanner) {
	for {
		fmt.Println("New batch.")
		for i := 0; i < 10; i++ {
			for out.Scan() {
				fmt.Println(string(out.Bytes()))
			}
		}
	}
}
