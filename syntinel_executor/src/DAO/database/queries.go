package database

const (
	GetScriptStatement    = "SELECT ID, Content, Hash FROM Script WHERE ID=?"
	InsertScriptStatement = "INSERT INTO Script(id, content, hash) VALUES (? ,?, ?)"
	UpdateScriptStatement = "UPDATE Script SET Content=?, Hash=? WHERE ID=?"
	DeleteScriptStatement = "DELETE FROM Script WHERE ID=?"

	GetDockerfileStatement    = "SELECT ID, Content, Hash FROM Dockerfile WHERE ID=?"
	InsertDockerfileStatement = "INSERT INTO Dockerfile(id, content, hash) VALUES (? ,?, ?)"
	UpdateDockerfileStatement = "UPDATE Dockerfile SET Content=?, Hash=? WHERE ID=?"
	DeleteDockerfileStatement = "DELETE FROM Dockerfile WHERE ID=?"

	GetTestStatement    = "SELECT ID, dockerfile, script FROM Test WHERE ID=?"
	InsertTestStatement = "INSERT INTO Test(id, dockerfile, script) VALUES (? ,?, ?)"
	UpdateTestStatement = "UPDATE Test SET dockerfile=?, script=? WHERE ID=?"
	DeleteTestStatement = "DELETE FROM Test WHERE ID=?"

	GetTestRunStatement    = "SELECT ID, test, environmentVariables, dockerfile, script FROM TestRun WHERE ID=?"
	InsertTestRunStatement = "INSERT INTO TestRun(id, test, environmentVariables, dockerfile, script) VALUES (? ,?, ?, ?, ?)"
	DeleteTestRunStatement = "DELETE FROM TestRun WHERE ID=?"
)
