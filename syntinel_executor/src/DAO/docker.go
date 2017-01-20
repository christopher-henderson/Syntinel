package DAO

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syntinel_executor/utils"
)

type Docker struct {
	ID int
}

func (d *Docker) Save(data io.Reader) {
	id := strconv.Itoa(d.ID)
	absolutePath := absDockerPath()
	path := absolutePath + id
	tmpPath := absolutePath + "." + id
	defer func() {
		if err := recover(); err != nil {
			// If there was an error, attempt to move the original (not 'tmp')
			// back to where it was.
			utils.FileCopy(tmpPath, path)
			utils.FileRemove(tmpPath)
			panic(err)
		}
	}()
	dst, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer dst.Close()
	if _, err := io.Copy(dst, data); err != nil {
		log.Fatalln(err)
	}
	if err := utils.FileRemove(tmpPath); err != nil {
		log.Fatalln(err)
	}
}

func (d *Docker) Delete() {
	path := absDockerPath() + strconv.Itoa(d.ID)
	if err := utils.FileRemove(path); err != nil {
		log.Fatalln(err)
	}
}

func (d *Docker) Path() string {
	return absDockerPath() + strconv.Itoa(d.ID)
}

func absDockerPath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}
	return path + string(os.PathSeparator) + "assets" +
		string(os.PathSeparator) + "dockers" +
		string(os.PathSeparator)
}
