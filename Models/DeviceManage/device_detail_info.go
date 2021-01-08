package DeviceManage

import (
	"SuperxonWebSite/Databases"
	"database/sql"
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

func GetAllDeviceBaseInfoList() (deviceBaseInfoList []DeviceBaseInfo, err error) {
	sqlStr := "select * from device_base_infos where id > ?"
	err = Databases.SuperxonDbDevice.Select(&deviceBaseInfoList, sqlStr, 0)
	if err != nil {
		return nil, err
	}
	return
}
