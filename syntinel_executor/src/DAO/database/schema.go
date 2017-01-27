package database

const Schema = `CREATE TABLE IF NOT EXISTS 'TestRun' (
	'id'	INTEGER NOT NULL UNIQUE,
	'testID'	INTEGER NOT NULL,
	'dockerfile'	TEXT NOT NULL,
	'script'	TEXT NOT NULL,
	'environmentVariables'	TEXT NOT NULL DEFAULT "",
	PRIMARY KEY('id')
);`
