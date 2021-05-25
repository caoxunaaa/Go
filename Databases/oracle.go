package Databases

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
)

var OracleDB *sql.DB

func InitOracle() {
	var err error
	//OracleDB, err = sql.Open("oci8", "PYD/123456@172.16.0.25:1521/Product_Test_DB")
	OracleDB, err = sql.Open("oci8", "CD8292/superxon@172.20.1.25:1521/Product_Test_DB")
	if err != nil {
		fmt.Println(err)
		return
	}
	OracleDB.SetMaxIdleConns(15)
	OracleDB.SetMaxOpenConns(200)
}

func CloseOracle() {
	_ = OracleDB.Close()
}
