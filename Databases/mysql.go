package Databases

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
)

var SuperxonDbDevice *sqlx.DB
var SuperxonProductionLineOracleRelationDb *sqlx.DB
var SuperxonDbDeviceOrm *gorm.DB

func InitMysql() {
	var err error
	SuperxonDbDeviceOrm, err = gorm.Open("mysql", "superxon:superxon@(172.20.3.12:3306)/superxon?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		fmt.Println(err)
		return
	}

	SuperxonDbDevice, err = sqlx.Open("mysql", "superxon:superxon@(172.20.3.12:3306)/superxon?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	SuperxonDbDevice.SetMaxIdleConns(10)
	SuperxonDbDevice.SetMaxOpenConns(50)

	SuperxonProductionLineOracleRelationDb, err = sqlx.Open("mysql", "superxon:superxon@(172.20.3.12:3306)/production-line-oracle-relation?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		fmt.Println("open production-line-oracle-relation failed,", err)
		return
	}
	SuperxonProductionLineOracleRelationDb.SetMaxIdleConns(10)
	SuperxonProductionLineOracleRelationDb.SetMaxOpenConns(50)
}

func CloseMysql() {
	_ = SuperxonDbDevice.Close()
	_ = SuperxonProductionLineOracleRelationDb.Close()
}
