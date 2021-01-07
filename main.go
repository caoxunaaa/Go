package main

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Models/DeviceManage"
	"SuperxonWebSite/Router"
	"fmt"
)

func main() {
	Databases.InitOracle()
	defer Databases.CloseOracle()

	Databases.InitSqlite3()
	Databases.SqliteDbDeviceOrm.AutoMigrate(
		&DeviceManage.DeviceBaseInfo{},
		&DeviceManage.DeviceRepairInfo{},
		&DeviceManage.DeviceMaintenanceItem{},
		&DeviceManage.DeviceMaintenanceCurrentInfo{},
		&DeviceManage.DeviceMaintenanceRecord{},
		&DeviceManage.DeviceTransmitInfo{},
		&DeviceManage.DeviceCategory{},
		&DeviceManage.SelfTest{},
	)
	_ = Databases.SqliteDbDeviceOrm.Close()
	defer Databases.CloseSqlite3()

	Databases.RedisInit()
	defer Databases.RedisClose()

	r := Router.Init()
	if err := r.Run("0.0.0.0:8002"); err != nil {
		fmt.Printf("startup service failed, err:%v\n\n", err)
	}
}
