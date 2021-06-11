package ModuleRunDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
	"strconv"
)

// 计划中的产品情况
type UndoneProjectPlanInfo struct {
	Id        int    `gorm:"primary_key" db:"id"`
	Type      string `db:"type"`
	Customers string `db:"customers"`
	Code      string `gorm:"unique;not null" db:"code"`
	Pn        string `gorm:"unique;not null" db:"pn"`
	PlanToPay int    `db:"plan_to_pay"` //计划交付数量
}

//已经完成的产品
type DoneProjectPlanInfo struct {
	Pn        string
	DoneToPay int //已经交付数量
}

type ProjectPlanInfo struct {
	UndoneProjectPlanInfo
	DoneProjectPlanInfo
	DoneRate string
}

func GetUndoneProjectPlanInfoList() (undoneProjectPlanInfoList []*UndoneProjectPlanInfo, err error) {
	sqlStr := "SELECT * from project_plan_infos"
	err = Databases.SuperxonProductionLineProductStatisticDevice.Select(&undoneProjectPlanInfoList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func CreateUndoneProjectPlanInfo(undoneProjectPlanInfo *UndoneProjectPlanInfo) (err error) {
	sqlStr := "INSERT INTO project_plan_infos(type, customers, code, pn, plan_to_pay) values (?, ?, ?, ?, ?)"
	_, err = Databases.SuperxonProductionLineProductStatisticDevice.Exec(sqlStr,
		undoneProjectPlanInfo.Type,
		undoneProjectPlanInfo.Customers,
		undoneProjectPlanInfo.Code,
		undoneProjectPlanInfo.Pn,
		undoneProjectPlanInfo.PlanToPay)
	if err != nil {
		return err
	}
	return
}

func DeleteUndoneProjectPlanInfo(id int) (err error) {
	sqlStr := "DELETE FROM project_plan_infos where id = ?"
	_, err = Databases.SuperxonProductionLineProductStatisticDevice.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	return
}

func UpdateUndoneProjectPlanInfo(undoneProjectPlanInfo *UndoneProjectPlanInfo, id int) (err error) {
	sqlStr := "UPDATE project_plan_infos SET type=?, customers=?, code=?, pn=?, plan_to_pay=? WHERE id=?"
	_, err = Databases.SuperxonProductionLineProductStatisticDevice.Exec(sqlStr,
		undoneProjectPlanInfo.Type,
		undoneProjectPlanInfo.Customers,
		undoneProjectPlanInfo.Code,
		undoneProjectPlanInfo.Pn,
		undoneProjectPlanInfo.PlanToPay,
		id)
	if err != nil {
		return err
	}
	return
}

/*
func GetProjectPlanList() (projectPlanInfoList []ProjectPlanInfo, err error)
从数据库中查询产品计划结果并返回
*/
func GetProjectPlanList() (projectPlanInfoList []ProjectPlanInfo, err error) {
	undoneProjectPlanInfoList := make([]UndoneProjectPlanInfo, 0)
	sqlStr := `SELECT * from project_plan_infos`
	err = Databases.SuperxonProductionLineProductStatisticDevice.Select(&undoneProjectPlanInfoList, sqlStr)
	if err != nil {
		return nil, err
	}

	doneProjectPlanInfoList := make([]DoneProjectPlanInfo, 0)
	startTime, endTime := Utils.GetCurrentAndZeroDayTime()
	sqlStr = `select model, count(*) from superxon.storagemanage_main a where a.shipmenttime between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss') group by a.model`
	rowsDone, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rowsDone.Close()
	var doneProjectPlanInfo DoneProjectPlanInfo
	for rowsDone.Next() {
		_ = rowsDone.Scan(
			&doneProjectPlanInfo.Pn,
			&doneProjectPlanInfo.DoneToPay)
		doneProjectPlanInfoList = append(doneProjectPlanInfoList, doneProjectPlanInfo)
	}

	if len(doneProjectPlanInfoList) != 0 {
		var doneCount int
		for _, valueUndone := range undoneProjectPlanInfoList {
			undone := valueUndone
			for _, valueDone := range doneProjectPlanInfoList {
				done := valueDone
				doneCount = 0
				if undone.Pn == done.Pn {
					doneCount = done.DoneToPay
					break
				}
			}
			projectPlanInfoList = append(projectPlanInfoList, ProjectPlanInfo{valueUndone, DoneProjectPlanInfo{Pn: valueUndone.Pn, DoneToPay: doneCount}, strconv.Itoa(doneCount*100/valueUndone.PlanToPay) + "%"})
		}
	} else {
		for _, valueUndone := range undoneProjectPlanInfoList {
			projectPlanInfoList = append(projectPlanInfoList, ProjectPlanInfo{valueUndone, DoneProjectPlanInfo{Pn: valueUndone.Pn, DoneToPay: 0}, "0%"})
		}
	}
	return
}
