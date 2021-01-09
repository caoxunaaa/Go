package DeviceManage

import (
	"SuperxonWebSite/Databases"
	"database/sql"
	"fmt"
	"time"
)

type DeviceMaintenanceItem struct {
	ID        uint   `gorm:"primary_key" db:"id"`
	Category  string `gorm:"not null" db:"category"`  //保养项目类型
	Name      string `gorm:"not null" db:"name"`      //某个保养类型中具体保养名称
	Period    uint   `gorm:"not null" db:"period"`    //保养间隔
	Threshold uint   `gorm:"not null" db:"threshold"` //保养提醒时间
}

type DeviceMaintenanceCurrentInfo struct {
	ID                  uint           `gorm:"primary_key" db:"id"`
	ItemCategory        string         `gorm:"not null" db:"item_category"`
	ItemName            string         `gorm:"not null" db:"item_name"`
	DeviceName          string         `gorm:"not null" db:"device_name"`
	DeviceSn            string         `gorm:"not null" db:"device_sn"`
	DeviceAssets        string         `gorm:"not null" db:"device_assets"`
	DeviceSort          sql.NullString `db:"device_sort"`
	DeviceOwner         string         `db:"device_owner"`
	LastMaintenanceTime time.Time      `db:"last_maintenance_time"`
	Deadline            time.Time      `db:"deadline"`
	StatusOfMaintenance string         `gorm:"default:'未绑定'" db:"status_of_maintenance"` //未绑定，正常，待保养，保养超时, 维修或报废
}

type DeviceMaintenanceRecord struct {
	ID              uint           `gorm:"primary_key" db:"id"`
	DeviceName      string         `db:"device_name"`
	DeviceSn        string         `db:"device_sn"`
	DeviceAssets    string         `db:"device_assets"`
	DeviceSort      string         `db:"device_sort"`
	ItemCategory    sql.NullString `db:"item_category"`
	ItemName        string         `db:"item_name"`
	MaintenanceTime time.Time      `db:"maintenance_time"`
	MaintenanceUser sql.NullString `db:"maintenance_user"`
	Remark          sql.NullString `db:"remark"`
	FilePath        sql.NullString `db:"file_path"`
}

func GetAllDeviceMaintenanceCategoryList() (deviceMaintenanceItemCategoryList []string, err error) {
	sqlStr := "select DISTINCT category from device_maintenance_items"
	err = Databases.SuperxonDbDevice.Select(&deviceMaintenanceItemCategoryList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func GetDeviceMaintenanceItemOfCategory(category string) (deviceMaintenanceItems []*DeviceMaintenanceItem, err error) {
	sqlStr := "select * from device_maintenance_items where category = ?"
	err = Databases.SuperxonDbDevice.Select(&deviceMaintenanceItems, sqlStr, category)
	if err != nil {
		return nil, err
	}
	return
}

func GetAllDeviceMaintenanceCurrentInfoList() (deviceMaintenanceCurrentInfoList []*DeviceMaintenanceCurrentInfo, err error) {
	sqlStr := "select * from device_maintenance_current_infos ORDER BY (CASE StatusOfMaintenance  WHEN '保养超时' THEN 1 WHEN '待保养' THEN 2 WHEN '未绑定' THEN 3 ELSE 4 END) ASC"
	err = Databases.SuperxonDbDevice.Select(&deviceMaintenanceCurrentInfoList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func GetDeviceMaintenanceCurrentInfo(snAssets string) (deviceMaintenanceCurrentInfo []*DeviceMaintenanceCurrentInfo, err error) {
	sqlStr := "select * from device_maintenance_current_infos where device_sn = ? or device_assets = ?"
	err = Databases.SuperxonDbDevice.Select(&deviceMaintenanceCurrentInfo, sqlStr, snAssets, snAssets)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return
}

func GetAllDeviceMaintenanceRecords(itemName string) (deviceMaintenanceRecords []*DeviceMaintenanceRecord, err error) {
	sqlStr := "select * from device_maintenance_records where item_name like ?"
	err = Databases.SuperxonDbDevice.Select(&deviceMaintenanceRecords, sqlStr, "%"+itemName+"%")
	if err != nil {
		return nil, err
	}
	return
}

func GetDeviceMaintenanceRecords(snAssets string, itemName string) (deviceMaintenanceRecords []*DeviceMaintenanceRecord, err error) {
	sqlStr := "select * from device_maintenance_records where item_name like ? and (device_sn = ? or device_assets = ?)"
	err = Databases.SuperxonDbDevice.Select(&deviceMaintenanceRecords, sqlStr, "%"+itemName+"%", snAssets, snAssets)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return
}
