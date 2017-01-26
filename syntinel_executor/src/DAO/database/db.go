package database

import (
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
	ExecuteTransactionalDDL(Schema)
}

func GetTestRun(id int) (*entities.TestRunEntity, error) {
	testRun := &entities.TestRunEntity{}
	// "SELECT id, testID, dockerfile, script, environmentVariables FROM TestRun WHERE ID=?"
	err := ExecuteTransactionalSingleRowQuery(
		GetTestRunStatement,
		[]interface{}{id},
		&testRun.ID,
		&testRun.TestID,
		&testRun.Dockerfile,
		&testRun.Script,
		&testRun.EnvironmentVariables)
	return testRun, err
}

func InsertTestRun(testRun *entities.TestRunEntity) error {
	return ExecuteTransactionalDDL(
		InsertTestRunStatement,
		testRun.ID,
		testRun.TestID,
		testRun.Dockerfile,
		testRun.Script,
		testRun.EnvironmentVariables)
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
