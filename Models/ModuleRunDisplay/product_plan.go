package ModuleRunDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
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
}

/*
func GetProjectPlanList() (projectPlanInfoList []ProjectPlanInfo, err error)
从数据库中查询产品计划结果并返回
*/
func GetProjectPlanList() (projectPlanInfoList []ProjectPlanInfo, err error) {
	undoneProjectPlanInfoList := make([]UndoneProjectPlanInfo, 0)
	sqlStr := `SELECT * from project_plan_infos`
	err = Databases.SuperxonDbDevice.Select(&undoneProjectPlanInfoList, sqlStr)
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
		for _, valueUndone := range undoneProjectPlanInfoList {
			for _, valueDone := range doneProjectPlanInfoList {
				if valueUndone.Pn == valueDone.Pn {
					projectPlanInfoList = append(projectPlanInfoList, ProjectPlanInfo{valueUndone, valueDone})
				}
			}
		}
	} else {
		for _, valueUndone := range undoneProjectPlanInfoList {
			projectPlanInfoList = append(projectPlanInfoList, ProjectPlanInfo{valueUndone, DoneProjectPlanInfo{Pn: valueUndone.Pn, DoneToPay: 0}})
		}
	}
	return
}

/*
func RedisGetProjectPlanList() (projectPlanInfoList []ProjectPlanInfo, err error)
获取redis缓存中的projectPlanInfoList，如果没有就重新在数据库中查询
*/
func RedisGetProjectPlanList() (projectPlanInfoList []ProjectPlanInfo, err error) {
	key := "projectPlanInfoList"
	reBytes, _ := redis.Bytes(Databases.RedisPool.Get().Do("get", key))
	if len(reBytes) != 0 {
		_ = json.Unmarshal(reBytes, &projectPlanInfoList)
		if len(projectPlanInfoList) != 0 {
			fmt.Println("使用redis")
			return
		}
	}

	projectPlanInfoList, err = GetProjectPlanList()

	datas, _ := json.Marshal(projectPlanInfoList)
	_, err = Databases.RedisPool.Get().Do("SET", key, datas)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = Databases.RedisPool.Get().Do("expire", key, 60*60*30)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}
