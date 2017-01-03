package DAO

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

type Docker struct {
	ID int
}

func (d *Docker) Save(file multipart.File) {
	id := strconv.Itoa(d.ID)
	absolutePath := absDockerPath()
	path := absolutePath + id
	tmpPath := absolutePath + "." + id
	defer cleanup(path, tmpPath)
	if err := copy(path, tmpPath); err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	if err := ioutil.WriteFile(path, data, 0770); err != nil {
		log.Fatalln(err)
	}
	if err := remove(tmpPath); err != nil {
		log.Fatalln(err)
	}
}

func (d *Docker) Delete() {
	path := absDockerPath() + strconv.Itoa(d.ID)
	if err := remove(path); err != nil {
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
