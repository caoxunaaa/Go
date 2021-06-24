package main

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Services"
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

	Services.KafkaInit()
	err = Services.InitCron()
	if err != nil {
		logx.Error(err)
		return
	}
	defer Services.CloseCron()
	select {}
}
