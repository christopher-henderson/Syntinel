package DAO

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// The Script type keeps the primary key of the script on the remote master
// server and uses that as the identifier on the filesystem.
type Script struct {
	ID int
}

func (s *Script) Save(content []byte) {
	id := strconv.Itoa(s.ID)
	absolutePath := abspath()
	path := absolutePath + id
	tmpPath := absolutePath + "." + id
	defer cleanup(path, tmpPath)
	if err := copy(path, tmpPath); err != nil {
		log.Fatalln(err)
	}
	if err := ioutil.WriteFile(path, content, 0770); err != nil {
		log.Fatalln(err)
	}
	if err := remove(tmpPath); err != nil {
		log.Fatalln(err)
	}
}

func (s *Script) Delete() {
	path := abspath() + strconv.Itoa(s.ID)
	if err := remove(path); err != nil {
		log.Fatalln(err)
	}
}

func (s *Script) Path() string {
	return abspath() + strconv.Itoa(s.ID)
}

func abspath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}
	return path + string(os.PathSeparator) + "assets" +
		string(os.PathSeparator) + "scripts" +
		string(os.PathSeparator)
}

func copy(source string, destination string) error {
	if _, err := os.Stat(source); err != nil {
		return nil
	}
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

func remove(path string) error {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

func cleanup(path string, tmp string) {
	if err := recover(); err != nil {
		copy(tmp, path)
		remove(tmp)
		panic(err)
	}
}
