package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"strconv"
	"time"
)

func main() {

	CreateItems()
}

type TestItem struct {
	Id       string
	Interval string
	NextRun  string
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func Random(min, max int) int {
	rand.Seed(time.Now().Unix() * 3)
	return rand.Intn(max-min) + min
}

func StoreItem(db *sql.DB, items []TestItem) {
	sql_additem := `
	INSERT OR REPLACE INTO items(
		Id,
		Interval,
		NextRun
	) values(?, ?, ?)
	`

	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.Id, item.Interval, item.NextRun)
		if err2 != nil {
			panic(err2)
		}
	}
}

func CreateItems() {
	const dbpath = "foo.db"
	db := InitDB(dbpath)
	CreateTable(db)

	for i := 0; i < 10; i++ {
		db := InitDB(dbpath)
		id := "fakeTest" + strconv.Itoa(i)
		interval := Random(2, 4)
		temp := strconv.Itoa(interval)
		item := []TestItem{{id, temp, "2017-01-14T00:00:42Z"}}
		StoreItem(db, item)
		time.Sleep(1 * time.Second)
	}
	//db := InitDB(dbpath)
	readItems := ReadItem(db)
	for i := 0; i < len(readItems); i++ {
		fmt.Println(readItems[i])
	}
}

func ReadItem(db *sql.DB) []TestItem {
	sql_readall := `
	SELECT Id, Interval, NextRun FROM items
	ORDER BY datetime(NextRun) DESC
	`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []TestItem
	for rows.Next() {
		item := TestItem{}
		err2 := rows.Scan(&item.Id, &item.Interval, &item.NextRun)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}

func CreateTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS items(
		Id TEXT NOT NULL PRIMARY KEY,
		Interval TEXT,
		NextRun DATETIME
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}
