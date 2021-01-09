package DeviceManage

import (
	"SuperxonWebSite/Databases"
	"database/sql"
	"fmt"
	"time"
)

type DeviceBaseInfo struct {
	ID                  uint           `gorm:"primary_key" db:"id"`
	StorageTime         time.Time      `db:"storage_time"`
	Name                string         `db:"name"`
	Sort                sql.NullString `db:"sort"`
	Sn                  string         `gorm:"unique;not null" db:"sn"`
	Assets              string         `gorm:"unique;not null" db:"assets"`
	CategoryRoot        string         `db:"category_root"`
	CategoryChild       string         `db:"category_child"`
	Owner               string         `db:"owner"`
	InternalCoding      sql.NullString `db:"internal_coding"`
	CalibrationType     string         `db:"calibration_type"`
	Supplier            sql.NullString `db:"supplier"`
	StatusOfRepair      string         `gorm:"default:'正常'" db:"status_of_repair"`       //正常，维修中，报废
	StatusOfMaintenance string         `gorm:"default:'未绑定'" db:"status_of_maintenance"` //未绑定，正常，待保养，保养超时, 维修或报废
}

type SelfTest struct {
	ID   uint   `gorm:"primary_key" db:"id"`
	Name string `db:"name"`
}

func GetAllDeviceBaseInfoList() (deviceBaseInfoList []*DeviceBaseInfo, err error) {
	sqlStr := "SELECT * FROM device_base_infos ORDER BY (CASE status_of_repair WHEN '维修中' THEN 1 WHEN '正常' THEN 2 ELSE 3 END), (CASE status_of_maintenance WHEN '保养超时' THEN 1 WHEN '待保养' THEN 2 WHEN '未绑定' THEN 3 ELSE 4 END) ASC"
	err = Databases.SuperxonDbDevice.Select(&deviceBaseInfoList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func GetDeviceBaseInfo(snAssetsIc string) (deviceBaseInfo *DeviceBaseInfo, err error) {
	deviceBaseInfo = new(DeviceBaseInfo)
	sqlStr := "select * from device_base_infos where sn = ? or assets = ? or internal_coding = ?"
	err = Databases.SuperxonDbDevice.Get(deviceBaseInfo, sqlStr, snAssetsIc, snAssetsIc, snAssetsIc)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return
}
