package DAO

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"syntinel_executor/utils"
)

// The Script type keeps the primary key of the script on the remote master
// server and uses that as the identifier on the filesystem.
type Script struct {
	ID int
}

func (s *Script) Save(content []byte) {
	path := s.Path()
	tmp := s.tmpPath()
	defer func() {
		if err := recover(); err != nil {
			// If there was an error, attempt to move the original (now called 'tmp')
			// back to where it was.
			utils.FileCopy(tmp, path)
			utils.FileRemove(tmp)
			panic(err)
		}
	}()
	// Copy the current script to a temporary file.
	if err := utils.FileCopy(path, tmp); err != nil {
		log.Println(err)
	}
	// Copy the incoming script to its final destination.
	if err := ioutil.WriteFile(path, content, 0770); err != nil {
		// This write was critical.
		log.Fatalln(err)
	}
	// Remove the temporary file.
	if err := utils.FileRemove(tmp); err != nil {
		log.Println(err)
	}
}

// Delete deletes the script from the filesystem.
func (s *Script) Delete() {
	if err := utils.FileRemove(s.Path()); err != nil {
		log.Println(err)
	}
}

// Path returns the absolute path of the script on the filesystem.
func (s *Script) Path() string {
	return fmt.Sprintf("%v%v", utils.ScriptDirectory(), strconv.Itoa(s.ID))
}

// tmpPath returns what Path does, but as a hidden file.
func (s *Script) tmpPath() string {
	return fmt.Sprintf("%v.%v", utils.ScriptDirectory(), strconv.Itoa(s.ID))
}
