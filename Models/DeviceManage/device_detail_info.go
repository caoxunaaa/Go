package DeviceManage

import (
	"SuperxonWebSite/Databases"
	"database/sql"
	"fmt"
)

type DeviceBaseInfo struct {
	ID                  uint           `gorm:"primary_key" db:"id"`
	StorageTime         string         `db:"storage_time"`
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

func DeleteDeviceBaseInfo(deviceSn string) (length int64, err error) {
	sqlStr := "DELETE FROM device_base_infos where sn = ?"
	res, err := Databases.SuperxonDbDevice.Exec(sqlStr, deviceSn)
	if err != nil {
		return length, err
	}
	length, err = res.RowsAffected()
	return
}

func CreateDeviceBaseInfo(deviceBaseInfo *DeviceBaseInfo) (err error) {
	sqlStr := "INSERT INTO device_base_infos(name, sort, sn, assets, category_root, category_child, owner, internal_coding, calibration_type, supplier, storage_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr,
		deviceBaseInfo.Name,
		deviceBaseInfo.Sort.String,
		deviceBaseInfo.Sn,
		deviceBaseInfo.Assets,
		deviceBaseInfo.CategoryRoot,
		deviceBaseInfo.CategoryChild,
		deviceBaseInfo.Owner,
		deviceBaseInfo.InternalCoding.String,
		deviceBaseInfo.CalibrationType,
		deviceBaseInfo.Supplier.String,
		deviceBaseInfo.StorageTime)
	if err != nil {
		return err
	}
	return
}

func UpdateDeviceBaseInfo(deviceBaseInfo *DeviceBaseInfo, oldSn string) (length int64, err error) {
	sqlStr := "UPDATE device_base_infos SET name=?, sort=?, sn=?, assets=?, category_root=?, category_child=?, owner=?, internal_coding=?, calibration_type=?, supplier=?, storage_time=?, status_of_repair=?, status_of_maintenance=? WHERE sn=?"
	res, err := Databases.SuperxonDbDevice.Exec(sqlStr,
		deviceBaseInfo.Name,
		deviceBaseInfo.Sort,
		deviceBaseInfo.Sn,
		deviceBaseInfo.Assets,
		deviceBaseInfo.CategoryRoot,
		deviceBaseInfo.CategoryChild,
		deviceBaseInfo.Owner,
		deviceBaseInfo.InternalCoding.String,
		deviceBaseInfo.CalibrationType,
		deviceBaseInfo.Supplier.String,
		deviceBaseInfo.StorageTime,
		deviceBaseInfo.StatusOfRepair,
		deviceBaseInfo.StatusOfMaintenance,
		oldSn)
	if err != nil {
		return length, err
	}
	length, err = res.RowsAffected()
	return
}
