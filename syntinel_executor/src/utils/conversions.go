package utils

import (
	"log"
	"strconv"
)

func AtoI(str string) int {
	integer, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalln(err)
	}
	return integer
}
