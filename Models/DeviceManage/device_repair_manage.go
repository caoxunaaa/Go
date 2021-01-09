package DeviceManage

import (
	"SuperxonWebSite/Databases"
	"database/sql"
	"time"
)

type DeviceRepairInfo struct {
	ID               uint           `gorm:"primary_key" db:"id"`
	Name             string         `db:"name"`
	Sort             sql.NullString `db:"sort"`
	Sn               string         `gorm:"not null" db:"sn"`
	Assets           string         `gorm:"not null" db:"assets"`
	RepairCategory   string         `db:"repair_category"` //内部修理, 外部维修
	Delegator        sql.NullString `db:"delegator"`       //委托人
	RepairFactory    sql.NullString `db:"repair_factory"`  //维修厂家
	SendToRepairTime time.Time      `db:"send_to_repair_time"`
	FinishTime       sql.NullTime   `db:"finish_time"`
	IsShelfLife      bool           `db:"is_shelf_life"`
	Reason           sql.NullString `db:"reason"`
	Solution         sql.NullString `db:"solution"`
	PR               sql.NullString `db:"pr"`
	Cost             uint32         `db:"cost"`
}

func GetAllDeviceRepairInfoList() (deviceRepairInfoList []*DeviceRepairInfo, err error) {
	sqlStr := "select * from device_repair_infos ORDER BY sn ASC"
	err = Databases.SuperxonDbDevice.Select(&deviceRepairInfoList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func GetDeviceRepairInfo(sn string) (deviceRepairInfo *DeviceRepairInfo, err error) {
	deviceRepairInfo = new(DeviceRepairInfo)
	sqlStr := "select * from device_repair_infos where sn = ?"
	err = Databases.SuperxonDbDevice.Get(deviceRepairInfo, sqlStr, sn)
	if err != nil {
		return nil, err
	}
	return
}
