package Databases

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func InitSqlite() (err error){
	DB, err = gorm.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		fmt.Println("db.sqlite3数据库连接失败")
		return
	}
	return DB.DB().Ping()
}

func CloseSqlite() {
	DB.Close()
}