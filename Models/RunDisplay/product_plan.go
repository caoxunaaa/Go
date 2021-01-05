package RunDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 计划中的产品情况
type UndoneProjectPlanInfo struct {
	Id        int
	Type      string
	Customers string
	Code      string
	Pn        string
	PlanToPay int //计划交付数量
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

func GetProjectPlanList() (projectPlanInfoList []ProjectPlanInfo, err error) {
	reBytes, _ := redis.Bytes(Databases.RedisConn.Do("get", "projectPlanInfoList"))
	_ = json.Unmarshal(reBytes, &projectPlanInfoList)
	if len(projectPlanInfoList) != 0 {
		fmt.Println("使用redis")
		return
	}

	undoneProjectPlanInfoList := make([]UndoneProjectPlanInfo, 0)
	stmt, _ := Databases.SqliteDbEntry.Prepare("SELECT * from ProjectPlanInfo")
	rowsUndone, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rowsUndone.Close()
	var undoneProjectPlanInfo UndoneProjectPlanInfo
	for rowsUndone.Next() {
		_ = rowsUndone.Scan(
			&undoneProjectPlanInfo.Id,
			&undoneProjectPlanInfo.Type,
			&undoneProjectPlanInfo.Customers,
			&undoneProjectPlanInfo.Code,
			&undoneProjectPlanInfo.Pn,
			&undoneProjectPlanInfo.PlanToPay)
		undoneProjectPlanInfoList = append(undoneProjectPlanInfoList, undoneProjectPlanInfo)
	}
	fmt.Println(undoneProjectPlanInfoList)

	doneProjectPlanInfoList := make([]DoneProjectPlanInfo, 0)
	startTime, endTime := Utils.GetCurrentAndZeroDayTime()
	sqlStr := `select model, count(*) from superxon.storagemanage_main a where a.shipmenttime between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss') group by a.model`
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
	fmt.Println(doneProjectPlanInfoList)

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

	fmt.Println(projectPlanInfoList)

	datas, _ := json.Marshal(projectPlanInfoList)
	_, _ = Databases.RedisConn.Do("SET", "projectPlanInfoList", datas)
	_, err = Databases.RedisConn.Do("expire", "projectPlanInfoList", 60*60)
	return
}
