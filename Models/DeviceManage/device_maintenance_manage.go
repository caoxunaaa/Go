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
	LastMaintenanceTime sql.NullTime   `db:"last_maintenance_time"`
	Deadline            sql.NullTime   `db:"deadline"`
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

//保养计划表
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

func CreateDeviceMaintenanceItem(deviceMaintenanceItem *DeviceMaintenanceItem) (err error) {
	//创建一个保养类型的保养项目，同时需要去判断此保养类型是否已经正在使用中
	sqlStr := "INSERT INTO device_maintenance_items(category, name, period, threshold) values (?, ?, ?, ?)"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr,
		deviceMaintenanceItem.Category,
		deviceMaintenanceItem.Name,
		deviceMaintenanceItem.Period,
		deviceMaintenanceItem.Threshold)
	if err != nil {
		return err
	}

	var deviceMaintenanceCurrentInfoList []*DeviceMaintenanceCurrentInfo //找到所有绑定了此种保养类型的Sn
	sqlStr = "select a.* from device_maintenance_current_infos a,(select DISTINCT device_sn from device_maintenance_current_infos where item_category = ?) b WHERE a.device_sn = b.device_sn"
	err = Databases.SuperxonDbDevice.Select(&deviceMaintenanceCurrentInfoList, sqlStr, deviceMaintenanceItem.Category)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if len(deviceMaintenanceCurrentInfoList) > 0 {
		for _, deviceMaintenanceCurrentInfo := range deviceMaintenanceCurrentInfoList {
			deviceMaintenanceCurrentInfo.ItemName = deviceMaintenanceItem.Name
			_ = CreateDeviceMaintenanceCurrentInfo(deviceMaintenanceCurrentInfo)
		}
	}
	return
}

func UpdateDeviceMaintenanceItem(deviceMaintenanceItem *DeviceMaintenanceItem, id uint) (err error) {
	//更新一个保养类型的保养项目，同时需要去判断此保养类型是否已经正在使用中
	sqlStr := "UPDATE device_maintenance_items SET category=?, name=?, period=?, threshold=? WHERE id=?"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr,
		deviceMaintenanceItem.Category,
		deviceMaintenanceItem.Name,
		deviceMaintenanceItem.Period,
		deviceMaintenanceItem.Threshold,
		id)
	if err != nil {
		return err
	}

	var deviceMaintenanceCurrentInfoList []*DeviceMaintenanceCurrentInfo //找到所有绑定了此种保养类型的Sn
	sqlStr = "select a.* from device_maintenance_current_infos a,(select DISTINCT device_sn from device_maintenance_current_infos where item_category = ?) b WHERE a.device_sn = b.device_sn and a.item_name = ?"
	err = Databases.SuperxonDbDevice.Select(&deviceMaintenanceCurrentInfoList, sqlStr, deviceMaintenanceItem.Category, deviceMaintenanceItem.Name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//1、保养类型已经有被使用
	if len(deviceMaintenanceCurrentInfoList) > 0 {
		for _, deviceMaintenanceCurrentInfo := range deviceMaintenanceCurrentInfoList {
			deviceMaintenanceCurrentInfo.ItemName = deviceMaintenanceItem.Name
			if !deviceMaintenanceCurrentInfo.LastMaintenanceTime.Time.IsZero() { //原本添加时间的情况
				deviceMaintenanceCurrentInfo.Deadline.Time = deviceMaintenanceCurrentInfo.LastMaintenanceTime.Time.AddDate(0, 0, int(deviceMaintenanceItem.Period))
				if deviceMaintenanceCurrentInfo.Deadline.Time.After(time.Now()) {
					deviceMaintenanceCurrentInfo.StatusOfMaintenance = "正常"
				} else if deviceMaintenanceCurrentInfo.LastMaintenanceTime.Time.AddDate(0, 0, int(deviceMaintenanceItem.Period-deviceMaintenanceItem.Threshold)).After(time.Now()) {
					deviceMaintenanceCurrentInfo.StatusOfMaintenance = "待保养"
				} else {
					deviceMaintenanceCurrentInfo.StatusOfMaintenance = "保养超时"
				}

			}
			_, _ = UpdateDeviceMaintenanceCurrentInfo(deviceMaintenanceCurrentInfo, deviceMaintenanceCurrentInfo.LastMaintenanceTime.Time.IsZero())
		}
	}
	return
}

func DeleteDeviceMaintenanceItem(deviceMaintenanceItem *DeviceMaintenanceItem, id uint) (err error) {
	sqlStr := "DELETE FROM device_maintenance_items where id = ?"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr, id)
	if err != nil {
		return
	}

	var deviceMaintenanceCurrentInfoList []*DeviceMaintenanceCurrentInfo //找到所有绑定了此种保养类型的所有Sn
	sqlStr = "select a.* from device_maintenance_current_infos a,(select DISTINCT device_sn from device_maintenance_current_infos where item_category = ?) b WHERE a.device_sn = b.device_sn and a.item_name = ?"
	err = Databases.SuperxonDbDevice.Select(&deviceMaintenanceCurrentInfoList, sqlStr, deviceMaintenanceItem.Category, deviceMaintenanceItem.Name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if len(deviceMaintenanceCurrentInfoList) > 0 {
		for _, deviceMaintenanceCurrentInfo := range deviceMaintenanceCurrentInfoList {
			_, _ = DeleteDeviceMaintenanceCurrentInfo(deviceMaintenanceCurrentInfo.ID)
		}
	}
	return
}

//绑定保养项目
func BindDeviceMaintenanceItem(deviceSn string, deviceMaintenanceItems []*DeviceMaintenanceItem) (err error) {
	deviceBaseInfo, err := GetDeviceBaseInfo(deviceSn)
	if err != nil {
		return
	}
	deviceMaintenanceCurrentInfo := DeviceMaintenanceCurrentInfo{
		DeviceName:          deviceBaseInfo.Name,
		DeviceSn:            deviceBaseInfo.Sn,
		DeviceAssets:        deviceBaseInfo.Assets,
		DeviceSort:          sql.NullString{String: deviceBaseInfo.Sort.String},
		DeviceOwner:         deviceBaseInfo.Owner,
		StatusOfMaintenance: "保养超时"}
	for _, deviceMaintenanceItem := range deviceMaintenanceItems {
		deviceMaintenanceCurrentInfo.ItemCategory = deviceMaintenanceItem.Category
		deviceMaintenanceCurrentInfo.ItemName = deviceMaintenanceItem.Name
		_ = CreateDeviceMaintenanceCurrentInfo(&deviceMaintenanceCurrentInfo)
	}
	deviceBaseInfo.StatusOfMaintenance = "保养超时"
	UpdateDeviceBaseInfo(deviceBaseInfo, deviceBaseInfo.Sn)
	return
}

//解绑保养项目
func UnBindDeviceMaintenanceItem(deviceSn string) (err error) {
	deviceMaintenanceCurrentInfos, err := GetDeviceMaintenanceCurrentInfo(deviceSn)
	for _, deviceMaintenanceCurrentInfo := range deviceMaintenanceCurrentInfos {
		_, err = DeleteDeviceMaintenanceCurrentInfo(deviceMaintenanceCurrentInfo.ID)
	}
	if err != nil {
		return
	}
	deviceBaseInfo, err := GetDeviceBaseInfo(deviceSn)
	if err != nil {
		return
	}
	deviceBaseInfo.StatusOfMaintenance = "未绑定"
	UpdateDeviceBaseInfo(deviceBaseInfo, deviceBaseInfo.Sn)
	return
}

//保养当前信息
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

func CreateDeviceMaintenanceCurrentInfo(deviceMaintenanceCurrentInfo *DeviceMaintenanceCurrentInfo) (err error) {
	sqlStr := "INSERT INTO device_maintenance_current_infos(item_category, item_name, device_name, device_sn, device_assets, device_sort, device_owner) values (?, ?, ?, ?, ?, ?, ?)"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr,
		deviceMaintenanceCurrentInfo.ItemCategory,
		deviceMaintenanceCurrentInfo.ItemName,
		deviceMaintenanceCurrentInfo.DeviceName,
		deviceMaintenanceCurrentInfo.DeviceSn,
		deviceMaintenanceCurrentInfo.DeviceAssets,
		deviceMaintenanceCurrentInfo.DeviceSort.String,
		deviceMaintenanceCurrentInfo.DeviceOwner)
	if err != nil {
		return err
	}
	return
}

func UpdateDeviceMaintenanceCurrentInfo(deviceMaintenanceCurrentInfoList *DeviceMaintenanceCurrentInfo, NoTime bool) (length int64, err error) {
	if NoTime {
		sqlStr := "UPDATE device_maintenance_current_infos SET item_category=?, item_name=?, device_name=?, device_sn=?, device_assets=?, device_sort=?, device_owner=?, status_of_maintenance=? WHERE id=?"
		res, err := Databases.SuperxonDbDevice.Exec(sqlStr,
			deviceMaintenanceCurrentInfoList.ItemCategory,
			deviceMaintenanceCurrentInfoList.ItemName,
			deviceMaintenanceCurrentInfoList.DeviceName,
			deviceMaintenanceCurrentInfoList.DeviceSn,
			deviceMaintenanceCurrentInfoList.DeviceAssets,
			deviceMaintenanceCurrentInfoList.DeviceSort,
			deviceMaintenanceCurrentInfoList.DeviceOwner,
			deviceMaintenanceCurrentInfoList.StatusOfMaintenance,
			deviceMaintenanceCurrentInfoList.ID)
		if err != nil {
			return length, err
		}
		length, err = res.RowsAffected()
	} else {
		sqlStr := "UPDATE device_maintenance_current_infos SET item_category=?, item_name=?, device_name=?, device_sn=?, device_assets=?, device_sort=?, device_owner=?, last_maintenance_time=?, deadline=?, status_of_maintenance=? WHERE id=?"
		res, err := Databases.SuperxonDbDevice.Exec(sqlStr,
			deviceMaintenanceCurrentInfoList.ItemCategory,
			deviceMaintenanceCurrentInfoList.ItemName,
			deviceMaintenanceCurrentInfoList.DeviceName,
			deviceMaintenanceCurrentInfoList.DeviceSn,
			deviceMaintenanceCurrentInfoList.DeviceAssets,
			deviceMaintenanceCurrentInfoList.DeviceSort,
			deviceMaintenanceCurrentInfoList.DeviceOwner,
			deviceMaintenanceCurrentInfoList.LastMaintenanceTime,
			deviceMaintenanceCurrentInfoList.Deadline,
			deviceMaintenanceCurrentInfoList.StatusOfMaintenance,
			deviceMaintenanceCurrentInfoList.ID)
		if err != nil {
			return length, err
		}
		length, err = res.RowsAffected()
	}

	return
}

func DeleteDeviceMaintenanceCurrentInfo(id uint) (length int64, err error) {
	sqlStr := "DELETE FROM device_maintenance_current_infos where id = ?"
	res, err := Databases.SuperxonDbDevice.Exec(sqlStr, id)
	if err != nil {
		return length, err
	}
	length, err = res.RowsAffected()
	return
}

//保养记录
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
