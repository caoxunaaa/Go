package main

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Models/DeviceManage"
	"SuperxonWebSite/Models/FileManage"
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Models/User"
	"SuperxonWebSite/Models/WaringDisplay"
	"SuperxonWebSite/Router"
	"SuperxonWebSite/Services"
	"fmt"
)

func main() {
	Databases.InitOracle()
	defer Databases.CloseOracle()

	Databases.InitMysql()
	defer Databases.CloseMysql()

	Databases.SuperxonDbDeviceOrm.AutoMigrate(
		&DeviceManage.DeviceBaseInfo{},
		&DeviceManage.DeviceRepairInfo{},
		&DeviceManage.DeviceMaintenanceItem{},
		&DeviceManage.DeviceMaintenanceCurrentInfo{},
		&DeviceManage.DeviceMaintenanceRecord{},
		&DeviceManage.DeviceTransmitInfo{},
		&DeviceManage.DeviceCategory{},
		&DeviceManage.SelfTest{},
		&ModuleRunDisplay.UndoneProjectPlanInfo{},
		&WaringDisplay.PnPassRateChartData{},
		&WaringDisplay.WarningCountChartData{},
		&User.Profile{},
		&FileManage.VideoInfo{},
	)
	_ = Databases.SuperxonDbDeviceOrm.Close()

	Databases.RedisInit()
	defer Databases.RedisClose()

	//Databases.InitMongoDb()
	//defer Databases.CloseMongoDb()

	Services.InitCron()
	defer Services.CloseCron()

	r := Router.Init()
	if err := r.Run("0.0.0.0:8002"); err != nil {
		fmt.Printf("startup service failed, err:%v\n\n", err)
	}
}
