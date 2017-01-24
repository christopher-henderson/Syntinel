package DAO

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syntinel_executor/utils"
)

type Docker struct {
	ID int
}

func (d *Docker) Save(content []byte) {
	path := d.Path()
	tmp := d.tmpPath()
	defer func() {
		if err := recover(); err != nil {
			// If there was an error, attempt to move the original (now called 'tmp')
			// back to where it was.
			utils.FileCopy(tmp, path)
			utils.FileRemove(tmp)
			panic(err)
		}
	}()
	// Copy the current Dockerfile to a temporary file.
	if err := utils.FileCopy(path, tmp); err != nil {
		log.Println(err)
	}
	// Copy the incoming Dockerfile to its final destination.
	if err := ioutil.WriteFile(path, content, 0770); err != nil {
		// This write was critical.
		log.Fatalln(err)
	}
	// Remove the temporary file.
	if err := utils.FileRemove(tmp); err != nil {
		log.Println(err)
	}
}

// Delete deletes the Dockerfile from the filesystem.
func (d *Docker) Delete() {
	if err := utils.FileRemove(d.Path()); err != nil {
		log.Println(err)
	}
}

// Path returns the absolute Dockerfile of the script on the filesystem.
func (d *Docker) Path() string {
	return fmt.Sprintf("%v%v%v", utils.DockerFileDirectory(), string(os.PathSeparator), strconv.Itoa(d.ID))
}

// tmpPath returns what Path does, but as a hidden file.
func (d *Docker) tmpPath() string {
	return fmt.Sprintf("%v%v.%v", utils.DockerFileDirectory(), string(os.PathSeparator), strconv.Itoa(d.ID))
}
