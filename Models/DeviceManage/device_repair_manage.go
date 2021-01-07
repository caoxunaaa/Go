package DeviceManage

import (
	"database/sql"
	"time"
)

type DeviceRepairInfo struct {
	ID               uint `gorm:"primary_key"`
	Name             string
	Sort             sql.NullString
	Sn               string         `gorm:"not null"`
	Assets           string         `gorm:"not null"`
	RepairCategory   string         //内部修理, 外部维修
	Delegator        sql.NullString //委托人
	RepairFactory    sql.NullString //维修厂家
	SendToRepairTime time.Time
	FinishTime       time.Time
	IsShelfLife      bool
	Reason           sql.NullString
	Solution         sql.NullString
	PR               sql.NullString
	Cost             uint32
}
