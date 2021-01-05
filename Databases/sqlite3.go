package Databases

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var SqliteDbEntry *sql.DB
var SqliteDbDevice *sql.DB

func InitSqlite3() {
	var err error
	SqliteDbEntry, err = sql.Open("sqlite3", "\\\\172.16.0.52\\superxon\\IT\\1.产线软件\\data\\Test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	SqliteDbEntry.SetMaxIdleConns(10)
	SqliteDbEntry.SetMaxOpenConns(20)

	SqliteDbDevice, err = sql.Open("sqlite3", "..\\sql\\db.sqlite3")
	if err != nil {
		fmt.Println(err)
		return
	}
	SqliteDbDevice.SetMaxIdleConns(10)
	SqliteDbDevice.SetMaxOpenConns(20)
}

func CloseSqlite3() {
	_ = SqliteDbEntry.Close()
	_ = SqliteDbDevice.Close()
}
