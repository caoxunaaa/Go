package DeviceManage

import (
	"SuperxonWebSite/Databases"
	"database/sql"
	"time"
)

type DeviceTransmitInfo struct {
	ID                   uint           `gorm:"primary_key" db:"id"`
	DeviceName           string         `db:"device_name"`
	DeviceCategoryRoot   string         `db:"device_category_root"`
	DeviceCategoryChild  string         `db:"device_category_child"`
	DeviceSort           sql.NullString `db:"device_sort"`
	DeviceSn             string         `gorm:"not null" db:"device_sn"`
	DeviceAssets         string         `gorm:"not null" db:"device_assets"`
	DeviceInternalCoding sql.NullString `db:"device_internal_coding"`
	OldOwner             string         `db:"old_owner"`
	NewOwner             string         `db:"new_owner"`
	TransmitTime         time.Time      `db:"transmit_time"`
}

func GetAllDeviceTransmitInfoList() (deviceTransmitInfoList []*DeviceTransmitInfo, err error) {
	sqlStr := "select * from device_transmit_infos"
	err = Databases.SuperxonDbDevice.Select(&deviceTransmitInfoList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func GetDeviceTransmitInfo(sn string) (deviceTransmitInfo *DeviceTransmitInfo, err error) {
	deviceTransmitInfo = new(DeviceTransmitInfo)
	sqlStr := "select * from device_transmit_infos where sn = ?"
	err = Databases.SuperxonDbDevice.Get(deviceTransmitInfo, sqlStr, sn)
	if err != nil {
		return nil, err
	}
	return
}

func GetDeviceTransmitInfoById(id int) (deviceTransmitInfo *DeviceTransmitInfo, err error) {
	deviceTransmitInfo = new(DeviceTransmitInfo)
	sqlStr := "select * from device_transmit_infos where id = ?"
	err = Databases.SuperxonDbDevice.Get(deviceTransmitInfo, sqlStr, id)
	if err != nil {
		return nil, err
	}
	return
}

func CreateDeviceTransmitInfo(deviceTransmitInfo *DeviceTransmitInfo) (err error) {
	sqlStr := "INSERT INTO device_transmit_infos(device_name, device_category_root, device_category_child, device_sort, device_sn, device_assets, device_internal_coding, old_owner, new_owner) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr,
		deviceTransmitInfo.DeviceName,
		deviceTransmitInfo.DeviceCategoryRoot,
		deviceTransmitInfo.DeviceCategoryChild,
		deviceTransmitInfo.DeviceSort.String,
		deviceTransmitInfo.DeviceSn,
		deviceTransmitInfo.DeviceAssets,
		deviceTransmitInfo.DeviceInternalCoding.String,
		deviceTransmitInfo.OldOwner,
		deviceTransmitInfo.NewOwner)
	if err != nil {
		return err
	}
	// 设备转移的同时设备基本信息中的所属者也发生变化
	//找到设备
	deviceBaseInfo, err := GetDeviceBaseInfo(deviceTransmitInfo.DeviceSn)
	if err != nil {
		return err
	}
	//设备SN
	sn := deviceBaseInfo.Sn
	//设备所属者改为转移之后的NewOwner
	deviceBaseInfo.Owner = deviceTransmitInfo.NewOwner
	UpdateDeviceBaseInfo(deviceBaseInfo, sn)
	return
}

func DeleteDeviceTransmitInfo(id int) (length int64, err error) {
	deviceTransmitInfo, err := GetDeviceTransmitInfoById(id)
	if err != nil {
		return length, err
	}
	sqlStr := "DELETE FROM device_transmit_infos where id = ?"
	res, err := Databases.SuperxonDbDevice.Exec(sqlStr, id)
	if err != nil {
		return length, err
	}
	length, err = res.RowsAffected()
	if length > 0 {
		// 设备转移的同时设备基本信息中的所属者也发生变化
		//找到id为id的设备转移记录
		//通过DeviceSn找到设备信息
		deviceBaseInfo, _ := GetDeviceBaseInfo(deviceTransmitInfo.DeviceSn)
		//设备所属者改为转移之后的NewOwner
		deviceBaseInfo.Owner = deviceTransmitInfo.OldOwner
		UpdateDeviceBaseInfo(deviceBaseInfo, deviceTransmitInfo.DeviceSn)
	}
	return
}
