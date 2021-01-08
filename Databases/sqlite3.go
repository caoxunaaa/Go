package Databases

import (
	"database/sql"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

var SqliteDbEntry *sql.DB

func InitSqlite3() {
	var err error
	SqliteDbEntry, err = sql.Open("sqlite3", "\\\\172.16.0.52\\superxon\\IT\\1.产线软件\\data\\Test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	SqliteDbEntry.SetMaxIdleConns(10)
	SqliteDbEntry.SetMaxOpenConns(50)
}

func CloseSqlite3() {
	_ = SqliteDbEntry.Close()
}
