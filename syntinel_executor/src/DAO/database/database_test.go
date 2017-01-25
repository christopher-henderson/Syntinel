package database

import "testing"

const DDL = "CREATE TABLE IF NOT EXISTS `Test` (" +
	"	`id`	INTEGER NOT NULL UNIQUE," +
	"	`name`	TEXT NOT NULL," +
	"	`dockerfile`	INTEGER NOT NULL," +
	"	`script`	INTEGER NOT NULL," +
	"	PRIMARY KEY(`id`)," +
	"	FOREIGN KEY(`dockerfile`) REFERENCES `Dockerfile`," +
	"	FOREIGN KEY(`script`) REFERENCES `Script`" +
	");" +

	"CREATE TABLE IF NOT EXISTS `Dockerfile` (" +
	"	`ID`	INTEGER NOT NULL UNIQUE," +
	"	`Content`	TEXT NOT NULL," +
	"	`Hash`	TEXT NOT NULL," +
	"	PRIMARY KEY(`id`)" +
	");" +

	"CREATE TABLE IF NOT EXISTS `Script` (" +
	"	`id`	INTEGER NOT NULL UNIQUE," +
	"	`content`	TEXT NOT NULL," +
	"	`hash`	TEXT NOT NULL," +
	"	PRIMARY KEY(`id`)" +
	");" +

	"CREATE TABLE IF NOT EXISTS `TestRun` (" +
	"	`id`	INTEGER NOT NULL UNIQUE," +
	"	`test`	INTEGER NOT NULL," +
	"	`environemntVariables`	TEXT," +
	"	`dockerfile`	TEXT NOT NULL," +
	"	`script`	TEXT NOT NULL," +
	"	PRIMARY KEY(`id`)," +
	"	FOREIGN KEY(`test`) REFERENCES `Test`" +
	");"

const dockerfile = `
 FROM docker.io/centos

 MAINTAINER Christopher Henderson

 RUN yum install -y go git wget
 COPY script.sh $HOME/script.sh
 CMD chmod +x script.sh && ./script.sh
 `

const script = `#!/usr/bin/env bash
 git clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover`

func TestMain(m *testing.M) {
	WriteDDL(DDL)
	m.Run()
}

func TestInsertDockerfile(t *testing.T) {
	defer clearDB()
	if err := InsertDockerfile(1, dockerfile); err != nil {
		t.Errorf("Got error inserting Dockerfile: %v", err)
	}
}

func TestInsertDockerfileRace(t *testing.T) {
	defer clearDB()
	for i := 1; i < 5; i++ {
		go func() {
			if err := InsertDockerfile(i, dockerfile); err != nil {
				t.Errorf("Got error inserting Dockerfile: %v", err)
			}
		}()
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
}

func TestInsertScript(t *testing.T) {
	defer clearDB()
	if err := InsertScript(1, script); err != nil {
		t.Errorf("Got error inserting Script: %v", err)
	}
}

func TestInsertScriptRace(t *testing.T) {
	defer clearDB()
	for i := 1; i < 5; i++ {
		go func() {
			if err := InsertScript(i, script); err != nil {
				t.Errorf("Got error inserting Script: %v", err)
			}
		}()
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
