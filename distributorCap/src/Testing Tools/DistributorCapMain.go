package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
)

var runInterval = 1

func main() {
	for {
		fmt.Println("Running cycle")
		TestAll()
		time.Sleep(10 * time.Second)

	}

}

type TestItem struct {
	Id       string
	Interval string
	NextRun  string
}

func (f *TestItem) SetNextRun(nextRun string) {
	f.NextRun = nextRun
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

func TestAll() {
	const dbpath = "foo.db"
	db := InitDB(dbpath)
	defer db.Close()
	readItems := ReadItem(db)
	Distribute(readItems)
}

func Distribute(tests []TestItem) {
	layout := "2006-01-02T15:04:05Z07:00"
	for i := 0; i < len(tests); i++ {
		t, _ := time.Parse(layout, tests[i].NextRun)
		temp := FindDiffrenceInMin(t, time.Now())
		if temp < runInterval {
			setTime := FindDiffrenceInSec(time.Now(), t)
			fmt.Print("This one should run ")
			fmt.Print(tests[i].Id)
			fmt.Print(" in ")
			fmt.Print(setTime)
			fmt.Println(" seconds")
			fmt.Print("Next Execution time ")
			fmt.Println()
			temp1 := time.Now().Format(time.RFC3339)
			temp2, err := time.Parse(layout, temp1)
			if err != nil {
				//do something
			}
			ScheduleAndRun(tests[i], temp2, setTime)
		}
	}
}

func ScheduleAndRun(test TestItem, t time.Time, setTime int) {
	const dbpath = "foo.db"
	nextRun := AddTime(t, test.Interval)
	fmt.Println(nextRun)
	test.SetNextRun(ConvertTimeToSqlFormat(nextRun))
	item := []TestItem{
		test,
	}
	db := InitDB(dbpath)
	defer db.Close()
	StoreItem(db, item)
	go PassToExe(test.Id, setTime)
}

func PassToExe(id string, seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
	fmt.Print("Fired off ")
	fmt.Println(id)
}

func FindDiffrenceInSec(now time.Time, exeTime time.Time) int {
	i := exeTime.Sub(now)
	return int(i.Seconds())
}

func FindDiffrenceInMin(nextRun time.Time, inital time.Time) int {
	i := nextRun.Sub(inital)
	return int(i.Minutes())
}

func AddTime(t time.Time, n string) time.Time {
	q, _ := strconv.ParseInt(n, 10, 8)
	j := t.Add(time.Minute * time.Duration(q))
	return j
}
func ConvertTimeToSqlFormat(t time.Time) string {
	formatedTime := t.Format(time.RFC3339)
	return formatedTime
}
