package DeviceManage

import (
	"database/sql"
	"time"
)

type DeviceTransmitInfo struct {
	ID                   uint `gorm:"primary_key"`
	DeviceName           string
	DeviceCategoryRoot   string
	DeviceCategoryChild  string
	DeviceSort           sql.NullString
	DeviceSn             string `gorm:"not null"`
	DeviceAssets         string `gorm:"not null"`
	DeviceInternalCoding sql.NullString
	OldOwner             string
	NewOwner             string
	TransmitTime         time.Time
}
