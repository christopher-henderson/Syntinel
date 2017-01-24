package ResultServer

import (
	"bufio"
	"fmt"
	"log"
)

const maxBackoff = 20

func Stream(ID int, stdout *bufio.Scanner) {
	go func() {
		for stdout.Scan() {
			log.Println(fmt.Sprintf("STDOUT ID %v: %v", ID, string(stdout.Bytes())))
		}
	}()
}
