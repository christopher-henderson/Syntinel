package DAO

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Script struct {
	Id int
}

func abspath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return path + string(os.PathSeparator) + "assets" +
		string(os.PathSeparator) + "scripts" +
		string(os.PathSeparator)
}

func copy(a string, b string) error {
	if _, err := os.Stat(a); err != nil {
		return nil
	}
	src, err := os.Open(a)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(a)
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

func (s *Script) Save(content []byte) {
	id := strconv.Itoa(s.Id)
	path := abspath() + id
	tmpPath := abspath() + "." + id
	log.Println(path)
	log.Println(tmpPath)
	defer cleanup(path, tmpPath)
	if err := copy(path, tmpPath); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(path, content, 0770); err != nil {
		log.Fatal(err)
	}
	// if err := remove(tmpPath); err != nil {
	// 	log.Fatal(err)
	// }
}
