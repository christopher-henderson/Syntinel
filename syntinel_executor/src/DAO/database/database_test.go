package database

import (
	"database/sql"
	"log"
	"os"
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

func TestInsertDockerfile(t *testing.T) {
	defer clearDB()
	if err := InsertDockerfile(1, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
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
	if _, err := GetDockerfile(1); err == nil {
		t.Errorf("Got a deleted Dockerfile.")
	}
}

func TestDeleteNonExistentDockerfile(t *testing.T) {
	if err := DeleteDockerfile(1); err != nil {
		t.Errorf("Error deleting Dockerfile when no such dockerfile exists: %v", err)
	}
}

func TestUpdateDockefile(t *testing.T) {
	defer clearDB()
	newContent := "this is new content"
	if err := InsertDockerfile(1, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
	}
	if err := UpdateDockerfile(1, newContent); err != nil {
		t.Errorf("Got an error updating a Dockerfile: %v", err)
	}
	if dockerfile, err := GetDockerfile(1); err != nil {
		t.Errorf("Got an error retrieving Dockerfile: %v", err)
	} else if dockerfile.Content != newContent {
		t.Errorf("Failed to update a Dockerfile. Expected it to change to %v, got %v", newContent, dockerfile.Content)
	}
}

func TestInsertScript(t *testing.T) {
	defer clearDB()
	if err := InsertScript(1, script); err != nil {
		t.Errorf("Got error inserting Script: %v", err)
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

func TestDeleteScript(t *testing.T) {
	defer clearDB()
	if err := InsertScript(1, script); err != nil {
		t.Errorf("Got error inserting Script: %v", err)
	}
	if err := DeleteScript(1); err != nil {
		t.Errorf("Got an error deleting a script: %v", err)
	}
	if s, err := GetScript(1); err == nil {
		t.Errorf("Failed to get an error when retrieving a deleted script. Got: %v", s)
	}
}

func TestDeleteNonExistentScript(t *testing.T) {
	clearDB()
	if err := DeleteScript(1000); err != nil {
		t.Errorf("Error deleting Script when no such script exists: %v", err)
	}
}

func TestUpdateScript(t *testing.T) {
	defer clearDB()
	newContent := "this is new content"
	if err := InsertScript(1, script); err != nil {
		t.Errorf("Got error inserting script %v", err)
	}
	if err := UpdateScript(1, newContent); err != nil {
		t.Errorf("Got an error updating a Script: %v", err)
	}
	if s, err := GetScript(1); err != nil {
		t.Errorf("Got an error retrieving script: %v", err)
	} else if s.Content != newContent {
		t.Errorf("Failed to update a script. Expected it to change to %v, got %v", newContent, s.Content)
	}
}

func TestInsertTest(t *testing.T) {
	defer clearDB()
	if err := InsertScript(1, script); err != nil {
		t.Errorf("Got error inserting script %v", err)
	}
	if err := InsertDockerfile(1, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
	}
	if err := InsertTest(1, 1, 1); err != nil {
		t.Errorf("Failed to insert test: %v", err)
	}
}

func TestInsertTestUnique(t *testing.T) {
	defer clearDB()
	if err := InsertScript(1, script); err != nil {
		t.Errorf("Got error inserting script %v", err)
	}
	if err := InsertDockerfile(1, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
	}
	if err := InsertTest(1, 1, 1); err != nil {
		t.Errorf("Failed to insert test: %v", err)
	}
	if err := InsertTest(1, 1, 1); err == nil {
		t.Errorf("Failed to enforce unique constraint on Test")
	}
}

func TestDeleteTest(t *testing.T) {
	defer clearDB()
	if err := InsertScript(1, script); err != nil {
		t.Errorf("Got error inserting script %v", err)
	}
	if err := InsertDockerfile(1, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
	}
	if err := InsertTest(1, 1, 1); err != nil {
		t.Errorf("Failed to insert test: %v", err)
	}
	if err := DeleteTest(1); err != nil {
		t.Errorf("Failed to delete Test: %v", err)
	}
	if _, err := GetTest(1); err == nil {
		t.Errorf("Failed to get an error when retrieving a deleted test. Got: %v", err)
	}
}

func TestDeleteNonExistentTest(t *testing.T) {
	if err := DeleteTest(1); err != nil {
		t.Errorf("Error deleting Test when no such test exists: %v", err)
	}
}

func TestUpdateTest(t *testing.T) {
	clearDB()
	defer clearDB()
	if err := InsertScript(1, script); err != nil {
		t.Errorf("Got error inserting script %v", err)
	}
	if err := InsertDockerfile(1, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
	}
	if err := InsertTest(1, 1, 1); err != nil {
		t.Errorf("Failed to insert test: %v", err)
	}
	if err := InsertScript(2, script); err != nil {
		t.Errorf("Got error inserting script %v", err)
	}
	if err := InsertDockerfile(2, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
	}
	if err := UpdateTest(1, 2, 2); err != nil {
		t.Errorf("Got an error updating a Test: %v", err)
	}
	if test, err := GetTest(1); err != nil {
		t.Errorf("Got an error retrieving a Test: %v", err)
	} else if test.Script != 2 || test.Dockerfile != 2 {
		t.Errorf("Failed to update Test. Expected Script to be %v, Dockerfile to be %v. Got %v and %v respectively.", 2, 2, test.Script, test.Dockerfile)
	}
}

// @TODO Constraints in general
// func TestInsertTestFKConstraint(t *testing.T) {
// 	defer clearDB()
// 	if err := InsertTest(1, 1, 1); err == nil {
// 		test, _ := GetTest(1)
// 		t.Errorf("Inserted a Test with a non-existent Script and Dockerfile. Test is %v", test)
// 	}
// }

func TestInsertTestRun(t *testing.T) {
	defer clearDB()
	if err := InsertTestRun(1, 1, "a=b", dockerfile, script); err != nil {
		t.Errorf("Failed to insert TestRun: %v", err)
	}
	if _, err := GetTestRun(1); err != nil {
		t.Errorf("Failed to retrieve TestRun: %v", err)
	}
}

func TestInsertTestRunUnique(t *testing.T) {
	defer clearDB()
	if err := InsertTestRun(1, 1, "a=b", dockerfile, script); err != nil {
		t.Errorf("Failed to insert TestRun: %v", err)
	}
	if err := InsertTestRun(1, 1, "a=b", dockerfile, script); err == nil {
		t.Errorf("Failed to enforce unique constraint on TestRun")
	}
}

func TestDeleteTestRun(t *testing.T) {
	defer clearDB()
	if err := InsertTestRun(1, 1, "a=b", dockerfile, script); err != nil {
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
	if err := InsertDockerfile(1, dockerfile); err == nil {
		t.Errorf("No error with deleted database.")
	}
	if _, err := GetDockerfile(1); err == nil {
		t.Errorf("No error with deleted database.")
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
