package database

const (
	GetTestRunStatement    = "SELECT id, testID, dockerfile, script, environmentVariables FROM TestRun WHERE ID=?"
	InsertTestRunStatement = "INSERT INTO TestRun(id, testID, dockerfile, script, environmentVariables) VALUES (? ,?, ?, ?, ?)"
	DeleteTestRunStatement = "DELETE FROM TestRun WHERE id=?"
)
