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
