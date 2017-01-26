package database

const DDL = "BEGIN TRANSACTION;" +
	"CREATE TABLE IF NOT EXISTS \"TestRun\" (" +
	"`id`	INTEGER NOT NULL UNIQUE," +
	"`test`	INTEGER NOT NULL," +
	"`environmentVariables`	TEXT," +
	"`dockerfile`	TEXT NOT NULL," +
	"`script`	TEXT NOT NULL," +
	"PRIMARY KEY(`id`)," +
	"FOREIGN KEY(`test`) REFERENCES `Test`" +
	");" +
	"CREATE TABLE IF NOT EXISTS \"Test\" (" +
	"`id`	INTEGER NOT NULL UNIQUE," +
	"`dockerfile`	INTEGER NOT NULL," +
	"`script`	INTEGER NOT NULL," +
	"PRIMARY KEY(`id`)," +
	"FOREIGN KEY(`dockerfile`) REFERENCES `Dockerfile`," +
	"FOREIGN KEY(`script`) REFERENCES `Script`" +
	");" +
	"CREATE TABLE IF NOT EXISTS \"Script\" (" +
	"`id`	INTEGER NOT NULL UNIQUE," +
	"`content`	TEXT NOT NULL," +
	"`hash`	TEXT NOT NULL," +
	"PRIMARY KEY(`id`)" +
	");" +
	"CREATE TABLE IF NOT EXISTS \"Dockerfile\" (" +
	"`id`	INTEGER NOT NULL UNIQUE," +
	"`content`	TEXT NOT NULL," +
	"`hash`	TEXT NOT NULL," +
	"PRIMARY KEY(`id`)" +
	");" +
	"COMMIT;"
