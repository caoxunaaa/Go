package main

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Models/DeviceManage"
	"SuperxonWebSite/Models/FileManage"
	"SuperxonWebSite/Models/User"
	"SuperxonWebSite/Router"
	"SuperxonWebSite/Services"
	"fmt"
)

func main() {
	Databases.InitOracle()
	defer Databases.CloseOracle()

	Databases.InitSqlite3()
	defer Databases.CloseSqlite3()

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
		&User.Profile{},
		&FileManage.VideoInfo{},
	)
	_ = Databases.SuperxonDbDeviceOrm.Close()

	Databases.RedisInit()
	defer Databases.RedisClose()

	Databases.InitMongoDb()
	defer Databases.CloseMongoDb()

	Services.InitCron()
	defer Services.CloseCron()

	r := Router.Init()
	if err := r.Run("0.0.0.0:8002"); err != nil {
		fmt.Printf("startup service failed, err:%v\n\n", err)
	}
}
