package main

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Router"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
)

func main() {
	var err error

	err = Databases.InitOracle()
	if err != nil {
		logx.Error(err)
		return
	}
	defer Databases.CloseOracle()

	err = Databases.InitMysql()
	if err != nil {
		logx.Error(err)
		return
	}
	defer Databases.CloseMysql()

	Databases.RedisInit()
	defer Databases.RedisClose()

	r := Router.Init()
	if err = r.Run("0.0.0.0:8002"); err != nil {
		fmt.Printf("startup service failed, err:%v\n\n", err)
	}
}
