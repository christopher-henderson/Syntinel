package database

import (
	"log"
	"testing"
)

const dockerfile = `
 FROM docker.io/centos

 MAINTAINER Christopher Henderson

 RUN yum install -y go git wget
 COPY script.sh $HOME/script.sh
 CMD chmod +x script.sh && ./script.sh
 `

const script = `#!/usr/bin/env bash
 git clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover`

func TestMain(m *testing.M) {
	WriteDDL(DDL)
	m.Run()
}

func TestInsertDockerfile(t *testing.T) {
	defer clearDB()
	if err := InsertDockerfile(1, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
	}
}

func TestInsertDockerfileRace(t *testing.T) {
	defer clearDB()
	for i := 1; i < 5; i++ {
		go func() {
			if err := InsertDockerfile(i, dockerfile); err != nil {
				t.Errorf("Got error inserting Dockerfile: %v", err)
			}
		}()
	}
}

func TestInsertDockerfileUnique(t *testing.T) {
	defer clearDB()
	if err := InsertDockerfile(1, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
	}
	if err := InsertDockerfile(1, dockerfile); err == nil {
		t.Errorf("Failed to enforce UNIQUE constraing on Dockerfile.ID")
	}
}

func TestDeleteDockerfile(t *testing.T) {
	defer clearDB()
	if err := InsertDockerfile(1, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
	}
	if err := DeleteDockerfile(1); err != nil {
		t.Errorf("Got error deleting dockerfile: %v", err)
	}
}

func TestInsertScript(t *testing.T) {
	defer clearDB()
	if err := InsertScript(1, script); err != nil {
		t.Errorf("Got error inserting Script: %v", err)
	}
}

func TestInsertScriptRace(t *testing.T) {
	defer clearDB()
	for i := 1; i < 5; i++ {
		go func() {
			if err := InsertScript(i, script); err != nil {
				t.Errorf("Got error inserting Script: %v", err)
			}
		}()
	}
}

func TestInsertScriptUnique(t *testing.T) {
	defer clearDB()
	if err := InsertScript(1, script); err != nil {
		t.Errorf("Got error inserting Script: %v", err)
	}
	if err := InsertScript(1, script); err == nil {
		t.Errorf("Failed to enforce UNIQUE constraing on Script.ID")
	}
}

func clearDB() {
	db := getDB()
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			log.Fatalln(err)
		}
		tables = append(tables, table)
	}
	rows.Close()
	for _, table := range tables {
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			log.Fatalln(err)
		}
	}
}
