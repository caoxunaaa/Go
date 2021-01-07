package Databases

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var SqliteDbEntry *sql.DB
var SqliteDbDevice *sqlx.DB
var SqliteDbDeviceOrm *gorm.DB

func InitSqlite3() {
	var err error
	SqliteDbEntry, err = sql.Open("sqlite3", "\\\\172.16.0.52\\superxon\\IT\\1.产线软件\\data\\Test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	SqliteDbEntry.SetMaxIdleConns(10)
	SqliteDbEntry.SetMaxOpenConns(20)

	SqliteDbDeviceOrm, err = gorm.Open("sqlite3", "D:\\CaoXun\\WorkProject\\Project\\Myself\\SuperxonWebSite\\sql\\db.sqlite3")
	if err != nil {
		fmt.Println(err)
		return
	}

	SqliteDbDevice, err = sqlx.Open("sqlite3", "D:\\CaoXun\\WorkProject\\Project\\Myself\\SuperxonWebSite\\sql\\db.sqlite3")
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
