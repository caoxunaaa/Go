package DeviceManage

import (
	"database/sql"
	"time"
)

type DeviceMaintenanceItem struct {
	ID        uint   `gorm:"primary_key"`
	Category  string `gorm:"not null"` //保养项目类型
	Name      string `gorm:"not null"` //某个保养类型中具体保养名称
	Period    uint   `gorm:"not null"` //保养间隔
	Threshold uint   `gorm:"not null"` //保养提醒时间
}

type DeviceMaintenanceCurrentInfo struct {
	ID                  uint   `gorm:"primary_key"`
	ItemCategory        string `gorm:"not null"`
	ItemName            string `gorm:"not null"`
	DeviceName          string `gorm:"not null"`
	DeviceSn            string `gorm:"not null"`
	DeviceAssets        string `gorm:"not null"`
	DeviceSort          sql.NullString
	DeviceOwner         string
	LastMaintenanceTime time.Time
	Deadline            time.Time
	StatusOfMaintenance string `gorm:"default:'未绑定'"` //未绑定，正常，待保养，保养超时, 维修或报废
}

type DeviceMaintenanceRecord struct {
	ID              uint `gorm:"primary_key"`
	DeviceName      string
	DeviceSn        string
	DeviceAssets    string
	DeviceSort      string
	ItemCategory    sql.NullString
	ItemName        string
	MaintenanceTime time.Time
	MaintenanceUser string
	Remark          sql.NullString
	FilePath        string
}
