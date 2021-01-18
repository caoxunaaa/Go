package DeviceManage

import (
	"SuperxonWebSite/Databases"
	"database/sql"
	"errors"
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
