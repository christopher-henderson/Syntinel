package database

import (
	"crypto/sha256"
	"database/sql"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile = "executor.sqlite3"
	driver = "sqlite3"
)

type DockerEntity struct {
	ID      int
	Content string
	Hash    string
}

type ScriptEntity struct {
	ID      int
	Content string
	Hash    string
}

func InitDB(file string) {
	db := getDB()
	defer db.Close()
	ddl, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	ddlContent, err := ioutil.ReadAll(ddl)
	if err != nil {
		panic(err)
	}
	WriteDDL(string(ddlContent))
}

func WriteDDL(ddl string) {
	db := getDB()
	_, err := db.Exec(ddl)
	if err != nil {
		panic(err)
	}
}

func GetDockerfile(id int) (*DockerEntity, error) {
	dockerfile := &DockerEntity{}
	err := getByID("SELECT ID, Content, Hash FROM Dockerfile WHERE ID=?", id, &dockerfile.ID,
		&dockerfile.Content,
		&dockerfile.Hash)
	return dockerfile, err
}

func InsertDockerfile(id int, content string) error {
	hash := sha256.Sum256([]byte(content))
	return transactionalQuery("INSERT INTO Dockerfile(id, content, hash) VALUES (? ,?, ?)", id, content, hash[:])
}

func UpdateDockerfile(id int, content string) error {
	hash := sha256.Sum256([]byte(content))
	return transactionalQuery("UPDATE Dockerfile SET Content=?, Hash=? WHERE ID=?", content, hash[:], id)
}

func DeleteDockerfile(id int) error {
	return transactionalQuery("DELETE FROM Dockerfile WHERE ID=?", id)
}

func GetScript(id int) (*ScriptEntity, error) {
	script := &ScriptEntity{}
	err := getByID("SELECT ID, Content, Hash FROM Dockerfile WHERE ID=?", id, &script.ID,
		&script.Content,
		&script.Hash)
	return script, err
}

func InsertScript(id int, content string) error {
	hash := sha256.Sum256([]byte(content))
	return transactionalQuery("INSERT INTO Script(id, content, hash) VALUES (? ,?, ?)", id, content, hash[:])
}

func UpdateScript(id int, content string) error {
	hash := sha256.Sum256([]byte(content))
	return transactionalQuery("UPDATE Script SET Content=?, Hash=? WHERE ID=?", content, hash[:], id)
}

func getByID(query string, id int, args ...interface{}) error {
	db := getDB()
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()
	if err := statement.QueryRow(id).Scan(args...); err != nil {
		return err
	}
	return nil
}

func transactionalQuery(query string, args ...interface{}) error {
	db := getDB()
	transaction, err := db.Begin()
	defer transaction.Commit()
	if err != nil {
		return err
	}
	statement, err := transaction.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()
	if _, err := statement.Exec(args...); err != nil {
		return err
	}
	if err := transaction.Commit(); err != nil {
		transaction.Rollback()
		return err
	}
	return nil
}

var getDB = func() func() *sql.DB {
	db, err := sql.Open(driver, dbFile)
	if err != nil {
		panic(err)
	}
	return func() *sql.DB {
		return db
	}
}()

func clearDB() {
	db := getDB()
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		panic(err)
	}
	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			panic(err)
		}
		tables = append(tables, table)
	}
	for _, table := range tables {
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			panic(err)
		}
	}
}
