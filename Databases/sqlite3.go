package Databases

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var SqliteDB *sql.DB

func InitSqlite3() {
	var err error
	SqliteDB, err = sql.Open("sqlite3", "\\\\172.16.0.52\\superxon\\IT\\1.产线软件\\data\\Test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CloseSqlite3() {
	_ = SqliteDB.Close()
}
