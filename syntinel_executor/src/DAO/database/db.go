package database

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"syntinel_executor/DAO/database/entities"
	"syntinel_executor/utils"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile = "executor.sqlite3"
	driver = "sqlite3"
)

func InitDB() {
	WriteDDL(DDL)
}

func WriteDDL(ddl string) {
	db := getDB()
	_, err := db.Exec(ddl)
	if err != nil {
		panic(err)
	}
}

func GetDockerfile(id int) (*entities.DockerfileEntity, error) {
	dockerfile := &entities.DockerfileEntity{}
	err := ExecuteTransactionalSingleRowQuery(
		GetDockerfileStatement, []interface{}{id}, &dockerfile.ID,
		&dockerfile.Content,
		&dockerfile.Hash)
	return dockerfile, err
}

func InsertDockerfile(id int, content string) error {
	hash := sha256.Sum256([]byte(content))
	return ExecuteTransactionalDDL(InsertDockerfileStatement, id, content, hash[:])
}

func UpdateDockerfile(id int, content string) error {
	hash := sha256.Sum256([]byte(content))
	return ExecuteTransactionalDDL(UpdateDockerfileStatement, content, hash[:], id)
}

func DeleteDockerfile(id int) error {
	return ExecuteTransactionalDDL(DeleteDockerfileStatement, id)
}

func GetScript(id int) (*entities.ScriptEntity, error) {
	script := &entities.ScriptEntity{}
	err := ExecuteTransactionalSingleRowQuery(
		GetScriptStatement, []interface{}{id}, &script.ID,
		&script.Content,
		&script.Hash)
	return script, err
}

func InsertScript(id int, content string) error {
	hash := sha256.Sum256([]byte(content))
	return ExecuteTransactionalDDL(InsertScriptStatement, id, content, hash[:])
}

func UpdateScript(id int, content string) error {
	hash := sha256.Sum256([]byte(content))
	return ExecuteTransactionalDDL(UpdateScriptStatement, content, hash[:], id)
}

func DeleteScript(id int) error {
	return ExecuteTransactionalDDL(DeleteScriptStatement, id)
}

func GetTest(id int) (*entities.TestEntity, error) {
	test := &entities.TestEntity{}
	err := ExecuteTransactionalSingleRowQuery(
		GetTestStatement, []interface{}{id}, &test.ID,
		&test.Dockerfile,
		&test.Script)
	return test, err
}

func InsertTest(id, dockerfile, script int) error {
	return ExecuteTransactionalDDL(InsertTestStatement, id, dockerfile, script)
}

func UpdateTest(id int, dockerfile, script int) error {
	return ExecuteTransactionalDDL(UpdateTestStatement, dockerfile, script, id)
}

func DeleteTest(id int) error {
	return ExecuteTransactionalDDL(DeleteTestStatement, id)
}

func GetTestRun(id int) (*entities.TestRunEntity, error) {
	testRun := &entities.TestRunEntity{}
	err := ExecuteTransactionalSingleRowQuery(
		GetTestRunStatement, []interface{}{id}, &testRun.ID,
		&testRun.Test,
		&testRun.EnvironmentVariables,
		&testRun.Dockerfile,
		&testRun.Script)
	return testRun, err
}

func InsertTestRun(id int, test int, environmentVariables, dockerfile, script string) error {
	return ExecuteTransactionalDDL(InsertTestRunStatement, id, test, environmentVariables, dockerfile, script)
}

func DeleteTestRun(id int) error {
	return ExecuteTransactionalDDL(DeleteTestRunStatement, id)
}

func ExecuteTransactionalDDL(query string, args ...interface{}) error {
	transaction, err := getDB().Begin()
	defer transaction.Commit()
	if err != nil {
		return err
	}
	stmt, err := transaction.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(args...); err != nil {
		transaction.Rollback()
		return err
	}
	return nil
}

func ExecuteTransactionalSingleRowQuery(query string, selection []interface{}, targets ...interface{}) error {
	transaction, err := getDB().Begin()
	if err != nil {
		return err
	}
	defer transaction.Commit()
	statement, err := transaction.Prepare(query)
	if err != nil {
		return err
	}
	row := statement.QueryRow(selection...)
	if err := row.Scan(targets...); err != nil {
		transaction.Rollback()
		return err
	}
	return nil
}

var getDB = func() func() *sql.DB {
	db, err := sql.Open(driver, fmt.Sprintf("%v%v", utils.DatabaseDirectory(), dbFile))
	if err != nil {
		panic(err)
	}
	return func() *sql.DB {
		return db
	}
}()
