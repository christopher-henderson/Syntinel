package database

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"syntinel_executor/utils"

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

type TestEntity struct {
	ID         int
	Dockerfile int
	Script     int
}

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

func GetDockerfile(id int) (*DockerEntity, error) {
	dockerfile := &DockerEntity{}
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

func GetScript(id int) (*ScriptEntity, error) {
	script := &ScriptEntity{}
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

func GetTest(id int) (*TestEntity, error) {
	test := &TestEntity{}
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
	return ExecuteTransactionalDDL(UpdateTestStatement, id, dockerfile, script)
}

func DeleteTest(id int) error {
	return ExecuteTransactionalDDL(DeleteTestStatement, id)
}

func GetTestRun(id int) (*TestRunEntity, error) {
	testRun := &TestRunEntity{}
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
	defer transaction.Commit()
	if err != nil {
		return err
	}
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

const GetScriptStatement = "SELECT ID, Content, Hash FROM Script WHERE ID=?"
const InsertScriptStatement = "INSERT INTO Script(id, content, hash) VALUES (? ,?, ?)"
const UpdateScriptStatement = "UPDATE Script SET Content=?, Hash=? WHERE ID=?"
const DeleteScriptStatement = "DELETE FROM Script WHERE ID=?"

const GetDockerfileStatement = "SELECT ID, Content, Hash FROM Dockerfile WHERE ID=?"
const InsertDockerfileStatement = "INSERT INTO Dockerfile(id, content, hash) VALUES (? ,?, ?)"
const UpdateDockerfileStatement = "UPDATE Dockerfile SET Content=?, Hash=? WHERE ID=?"
const DeleteDockerfileStatement = "DELETE FROM Dockerfile WHERE ID=?"

const GetTestStatement = "SELECT ID, dockerfile, script FROM Test WHERE ID=?"
const InsertTestStatement = "INSERT INTO Test(id, dockerfile, script) VALUES (? ,?, ?)"
const UpdateTestStatement = "UPDATE Test SET dockerfile=?, script=? WHERE ID=?"
const DeleteTestStatement = "DELETE FROM Test WHERE ID=?"

const GetTestRunStatement = "SELECT ID, test, environmentVariables, dockerfile, script FROM TestRun WHERE ID=?"
const InsertTestRunStatement = "INSERT INTO TestRun(id, test, environmentVariables, dockerfile, script) VALUES (? ,?, ?, ?, ?)"
const DeleteTestRunStatement = "DELETE FROM TestRun WHERE ID=?"
