package Databases

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
)

var SuperxonProductionLineProductStatisticDevice *sqlx.DB
var SuperxonProductionLineOracleRelationDb *sqlx.DB
var SuperxonHumanResourcesDb *sqlx.DB
var SuperxonUserDb *sqlx.DB

func InitMysql() error {
	var err error
	SuperxonProductionLineProductStatisticDevice, err = sqlx.Open("mysql", "superxon:superxon@(172.20.3.12:3306)/production-line-product-statistic?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return err
	}
	SuperxonProductionLineProductStatisticDevice.SetMaxIdleConns(10)
	SuperxonProductionLineProductStatisticDevice.SetMaxOpenConns(50)

	SuperxonProductionLineOracleRelationDb, err = sqlx.Open("mysql", "superxon:superxon@(172.20.3.12:3306)/production-line-oracle-relation?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		fmt.Println("open production-line-oracle-relation failed,", err)
		return err
	}
	SuperxonProductionLineOracleRelationDb.SetMaxIdleConns(10)
	SuperxonProductionLineOracleRelationDb.SetMaxOpenConns(50)

	SuperxonHumanResourcesDb, err = sqlx.Open("mysql", "superxon:superxon@(172.20.3.12:3306)/superxon-human-resources?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		fmt.Println("open superxon-human-resources failed,", err)
		return err
	}
	SuperxonHumanResourcesDb.SetMaxIdleConns(10)
	SuperxonHumanResourcesDb.SetMaxOpenConns(50)

	SuperxonUserDb, err = sqlx.Open("mysql", "superxon:superxon@(172.20.3.12:3306)/superxon-user?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		fmt.Println("open superxon-user failed,", err)
		return err
	}
	SuperxonUserDb.SetMaxIdleConns(10)
	SuperxonUserDb.SetMaxOpenConns(50)
	return nil
}

func CloseMysql() {
	_ = SuperxonProductionLineProductStatisticDevice.Close()
	_ = SuperxonProductionLineOracleRelationDb.Close()
	_ = SuperxonHumanResourcesDb.Close()
	_ = SuperxonUserDb.Close()
}
