package database

import (
	"database/sql"
	"log"
	"os"
	"syntinel_executor/DAO/database/entities"
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

var getDBTest = func() func() *sql.DB {
	db, err := sql.Open(driver, dbFile)
	if err != nil {
		panic(err)
	}
	return func() *sql.DB {
		return db
	}
}()

func TestMain(m *testing.M) {
	getDB = getDBTest
	InitDB()
	m.Run()
	os.Remove(dbFile)
}

func TestInsertTestRun(t *testing.T) {
	defer clearDB()
	tr := &entities.TestRunEntity{}
	tr.ID = 1
	tr.TestID = 1
	tr.Dockerfile = dockerfile
	tr.Script = script
	tr.EnvironmentVariables = ""
	if err := InsertTestRun(tr); err != nil {
		t.Errorf("Failed to insert TestRun: %v", err)
	}
	if tr, err := GetTestRun(1); err != nil {
		t.Errorf("Failed to retrieve TestRun: %v", err)
	} else if tr.Dockerfile != dockerfile || tr.Script != script {
		log.Println(tr.Dockerfile)
		t.Errorf("Got back a test run with bad Dockerfile and Script:\n%v\n%v", tr.Dockerfile, tr.Script)
	}
}

func TestInsertTestRunUnique(t *testing.T) {
	defer clearDB()
	tr := &entities.TestRunEntity{}
	tr.ID = 1
	tr.TestID = 1
	tr.Dockerfile = dockerfile
	tr.Script = script
	tr.EnvironmentVariables = ""
	if err := InsertTestRun(tr); err != nil {
		t.Errorf("Failed to insert TestRun: %v", err)
	}
	if err := InsertTestRun(tr); err == nil {
		t.Errorf("Failed to enforce unique constraint on TestRun")
	}
}

func TestDeleteTestRun(t *testing.T) {
	defer clearDB()
	tr := &entities.TestRunEntity{}
	tr.ID = 1
	tr.TestID = 1
	tr.Dockerfile = dockerfile
	tr.Script = script
	tr.EnvironmentVariables = ""
	if err := InsertTestRun(tr); err != nil {
		t.Errorf("Failed to insert TestRun: %v", err)
	}
	if err := DeleteTestRun(1); err != nil {
		t.Errorf("Faield to delete TestRun: %v", err)
	}
}

func TestBadDDLQuery(t *testing.T) {
	if err := ExecuteTransactionalDDL("nope"); err == nil {
		t.Errorf("Failed to throw error on bad query")
	}
}

func TestBadQuery(t *testing.T) {
	if err := ExecuteTransactionalSingleRowQuery("nope", []interface{}{"things"}); err == nil {
		t.Errorf("Failed to throw error on bad query.")
	}
}

func TestBadDB(t *testing.T) {
	os.RemoveAll(dbFile)
	defer func() {
		getDB().Close()
		getDB = func() func() *sql.DB {
			db, err := sql.Open(driver, dbFile)
			if err != nil {
				panic(err)
			}
			return func() *sql.DB {
				return db
			}
		}()
		InitDB()
	}()
	tr := &entities.TestRunEntity{}
	tr.ID = 1
	tr.TestID = 1
	tr.Dockerfile = dockerfile
	tr.Script = script
	tr.EnvironmentVariables = ""
	if err := InsertTestRun(tr); err == nil {
		t.Errorf("Did not recieve error on deleted database.")
	}
	if _, err := GetTestRun(1); err == nil {
		t.Errorf("Did not recieve error on deleted database.")
	}
}

// @TODO figure these two out. Read http://go-database-sql.org/index.html
//
// func TestBadTransaction(t *testing.T) {
// 	defer clearDB()
// 	tx, _ := getDB().Begin()
// 	defer tx.Commit()
// 	if err := InsertDockerfile(1, dockerfile); err == nil {
// 		t.Errorf("No error with deleted database.")
// 	}
// }
//
// func TestRollback(t *testing.T) {
// 	if err := InsertScript(1, script); err != nil {
// 		t.Errorf("Got error inserting script %v", err)
// 	}
// 	query := "UPDATE Script SET Content='derp' WHERE ID=1;UPDATE ASd SET Content='asd' where ID=1;"
// 	if err := ExecuteTransactionalDDL(query); err == nil {
// 		s, _ := GetScript(1)
// 		t.Errorf("No error on bad query, got script %v", s.Content)
// 	}
// }

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
