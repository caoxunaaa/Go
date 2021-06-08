package main

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Models/FileManage"
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Models/WaringDisplay"
	"SuperxonWebSite/Router"
	"SuperxonWebSite/Services"
	"fmt"
)

func main() {
	var err error

	Databases.InitOracle()
	defer Databases.CloseOracle()

	Databases.InitMysql()
	defer Databases.CloseMysql()

	Databases.SuperxonDbDeviceOrm.AutoMigrate(
		&ModuleRunDisplay.UndoneProjectPlanInfo{},
		&WaringDisplay.PnPassRateChartData{},
		&WaringDisplay.WarningCountChartData{},
		&FileManage.VideoInfo{},
	)
	err = Databases.SuperxonDbDeviceOrm.Close()
	if err != nil {
		fmt.Println(err)
	}

	Databases.RedisInit()
	defer Databases.RedisClose()

	Services.KafkaInit()
	err = Services.InitCron()
	if err != nil {
		fmt.Println(err)
	}
	defer Services.CloseCron()

	r := Router.Init()
	if err := r.Run("0.0.0.0:8002"); err != nil {
		fmt.Printf("startup service failed, err:%v\n\n", err)
	}

}
