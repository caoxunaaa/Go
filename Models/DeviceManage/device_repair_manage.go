package DeviceManage

import (
	"SuperxonWebSite/Databases"
	"database/sql"
	"errors"
	"fmt"
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

type Test struct {
	DeviceName          string    `db:"device_name"`
	DeviceAssets        string    `db:"device_assets"`
	DeviceOwner         string    `db:"device_owner"`
	Deadline            time.Time `db:"deadline"`
	ItemName            string    `db:"item_name"`
	ItemCategory        string    `db:"item_category"`
	DeviceSn            string    `db:"device_sn"`
	DeviceSort          string    `db:"device_sort"`
	LastMaintenanceTime time.Time `db:"last_maintenance_time"`
	Id                  int64     `db:"id"`
	Status              string    `db:"status"`
}

func GetTest() (test []*Test, err error) {
	sqlStr := fmt.Sprintf(`SELECT c.*, b.status FROM device_maintenance_current_infos c, 
        (SELECT a.*, CASE WHEN a.v > a.threshold then '正常' WHEN a.v > 0 then '待保养' ELSE '保养超时' END 'status' 
        FROM (SELECT t.id, t.device_sn, t.deadline, c.threshold, TIMESTAMPDIFF(DAY, CURRENT_DATE, t.deadline) as v 
        FROM device_maintenance_current_infos t, device_maintenance_items c 
        WHERE t.item_category= c.category AND t.item_name= c.name AND t.device_sn = '091070022'
        ORDER BY v) a) b where c.device_sn = b.device_sn AND c.id = b.id`)
	fmt.Println(sqlStr)
	err = Databases.SuperxonDbDevice.Select(&test, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func GetAllDeviceRepairInfoList() (deviceRepairInfoList []*DeviceRepairInfo, err error) {
	sqlStr := "select * from device_repair_infos ORDER BY sn ASC"
	err = Databases.SuperxonDbDevice.Select(&deviceRepairInfoList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func GetDeviceRepairInfo(sn string) (deviceRepairInfo []*DeviceRepairInfo, err error) {
	sqlStr := "select * from device_repair_infos where sn = ? order by send_to_repair_time DESC"
	err = Databases.SuperxonDbDevice.Select(&deviceRepairInfo, sqlStr, sn)
	if err != nil {
		return nil, err
	}
	return
}

func CreateDeviceRepairInfo(deviceRepairInfo *DeviceRepairInfo, repairStatus string) (err error) {
	//创建维修记录成功的时候进行设备信息中维修状态的改变
	//找到设备
	deviceBaseInfo, err := GetDeviceBaseInfo(deviceRepairInfo.Sn)
	if err != nil {
		return err
	}
	if deviceRepairInfo.FinishTime.Time.IsZero() {
		deviceBaseInfo.StatusOfRepair = "维修中"
	} else {
		if repairStatus == "正常" || repairStatus == "报废" {
			deviceBaseInfo.StatusOfRepair = repairStatus //报废或者正常
		} else {
			return errors.New("维修完成状态必须为报废或者正常")
		}
	}
	//创建维修信息
	if deviceRepairInfo.FinishTime.Time.IsZero() {
		sqlStr := "INSERT INTO device_repair_infos(name, sort, sn, assets, repair_category, delegator, repair_factory, send_to_repair_time, is_shelf_life, reason, solution, pr, cost) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err = Databases.SuperxonDbDevice.Exec(sqlStr,
			deviceRepairInfo.Name,
			deviceRepairInfo.Sort.String,
			deviceRepairInfo.Sn,
			deviceRepairInfo.Assets,
			deviceRepairInfo.RepairCategory,
			deviceRepairInfo.Delegator.String,
			deviceRepairInfo.RepairFactory.String,
			deviceRepairInfo.SendToRepairTime,
			deviceRepairInfo.IsShelfLife,
			deviceRepairInfo.Reason.String,
			deviceRepairInfo.Solution.String,
			deviceRepairInfo.PR.String,
			deviceRepairInfo.Cost,
		)
		if err != nil {
			return err
		}
	} else {
		sqlStr := "INSERT INTO device_repair_infos(name, sort, sn, assets, repair_category, delegator, repair_factory, send_to_repair_time, finish_time, is_shelf_life, reason, solution, pr, cost) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err = Databases.SuperxonDbDevice.Exec(sqlStr,
			deviceRepairInfo.Name,
			deviceRepairInfo.Sort.String,
			deviceRepairInfo.Sn,
			deviceRepairInfo.Assets,
			deviceRepairInfo.RepairCategory,
			deviceRepairInfo.Delegator.String,
			deviceRepairInfo.RepairFactory.String,
			deviceRepairInfo.SendToRepairTime,
			deviceRepairInfo.FinishTime.Time,
			deviceRepairInfo.IsShelfLife,
			deviceRepairInfo.Reason.String,
			deviceRepairInfo.Solution.String,
			deviceRepairInfo.PR.String,
			deviceRepairInfo.Cost,
		)
		if err != nil {
			return err
		}
	}

	//更改设备信息中的维修状态
	_, err = UpdateDeviceBaseInfo(deviceBaseInfo, deviceRepairInfo.Sn)
	if err != nil {
		return err
	}
	return
}

func UpdateDeviceRepairInfo(deviceRepairInfo *DeviceRepairInfo, oldId uint, repairStatus string) (length int64, err error) {
	//创建维修记录成功的时候进行设备信息中维修状态的改变
	//找到设备
	deviceBaseInfo, err := GetDeviceBaseInfo(deviceRepairInfo.Sn)
	if err != nil {
		return
	}
	if deviceRepairInfo.FinishTime.Time.IsZero() {
		deviceBaseInfo.StatusOfRepair = "维修中"
	} else {
		if repairStatus == "正常" || repairStatus == "报废" {
			deviceBaseInfo.StatusOfRepair = repairStatus //报废或者正常
		} else {
			return 0, errors.New("维修完成状态必须为报废或者正常")
		}
	}

	if deviceRepairInfo.FinishTime.Time.IsZero() {
		sqlStr := "UPDATE device_repair_infos SET name=?, sort=?, sn=?, assets=?, repair_category=?, delegator=?, repair_factory=?, send_to_repair_time=?, is_shelf_life=?, reason=?, solution=?, pr=?, cost=? WHERE id=?"
		res, err := Databases.SuperxonDbDevice.Exec(sqlStr,
			deviceRepairInfo.Name,
			deviceRepairInfo.Sort.String,
			deviceRepairInfo.Sn,
			deviceRepairInfo.Assets,
			deviceRepairInfo.RepairCategory,
			deviceRepairInfo.Delegator.String,
			deviceRepairInfo.RepairFactory.String,
			deviceRepairInfo.SendToRepairTime,
			deviceRepairInfo.IsShelfLife,
			deviceRepairInfo.Reason.String,
			deviceRepairInfo.Solution.String,
			deviceRepairInfo.PR.String,
			deviceRepairInfo.Cost,
			oldId)
		if err != nil {
			return length, err
		}
		length, err = res.RowsAffected()
	} else {
		sqlStr := "UPDATE device_repair_infos SET name=?, sort=?, sn=?, assets=?, repair_category=?, delegator=?, repair_factory=?, send_to_repair_time=?, finish_time=?, is_shelf_life=?, reason=?, solution=?, pr=?, cost=? WHERE id=?"
		res, err := Databases.SuperxonDbDevice.Exec(sqlStr,
			deviceRepairInfo.Name,
			deviceRepairInfo.Sort.String,
			deviceRepairInfo.Sn,
			deviceRepairInfo.Assets,
			deviceRepairInfo.RepairCategory,
			deviceRepairInfo.Delegator.String,
			deviceRepairInfo.RepairFactory.String,
			deviceRepairInfo.SendToRepairTime,
			deviceRepairInfo.FinishTime.Time,
			deviceRepairInfo.IsShelfLife,
			deviceRepairInfo.Reason.String,
			deviceRepairInfo.Solution.String,
			deviceRepairInfo.PR.String,
			deviceRepairInfo.Cost,
			oldId)
		if err != nil {
			return length, err
		}
		length, err = res.RowsAffected()
	}

	length, err = UpdateDeviceBaseInfo(deviceBaseInfo, deviceRepairInfo.Sn)
	if err != nil {
		return length, err
	}

	return
}
