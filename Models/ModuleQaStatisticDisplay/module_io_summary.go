package ModuleQaStatisticDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
	"reflect"
)

type IOSummaryInfo struct {
	Data    string
	Pn      string
	Process string
	Input   int
	Pass    int
	Bad     int
}

func Get10GLineIOSummaryInfoList(q *QueryCondition) ([]interface{}, error) {
	ioSummaryInfoList := make([]IOSummaryInfo, 0)
	sqlStr := `SELECT x.* FROM (select DISTINCT to_char(ACTION_TIME,'yyyy-mm-dd')日期,Pn AS 型号 ,LOG_ACTION AS 工序
,COUNT(sn)OVER(PARTITION BY pn,LOG_action,to_char(ACTION_TIME,'yyyy-mm-dd'))投入
,SUM(case P_VALUE when 'PASS' then 1 else 0 end)OVER(PARTITION BY pn,LOG_action,to_char(ACTION_TIME,'yyyy-mm-dd') ) 良品
,SUM(case P_VALUE when 'PASS' then 0 else 1 end)OVER(PARTITION BY pn,LOG_action,to_char(ACTION_TIME,'yyyy-mm-dd') ) 不良品
 from (SELECT c.*,rank() over(partition by sn,log_action
                        order by action_time desc)rr
FROM superxon.AutoDT_Process_LOG c
WHERE ACTION_TIME >=to_date('` + q.StartTime + `','yyyy-mm-dd hh24:mi:ss')
and ACTION_TIME <=to_date('` + q.EndTime + `','yyyy-mm-dd hh24:mi:ss')
) a
WHERE a.rr=1 AND a.pn like 'SO%62%'
order by to_char(ACTION_TIME,'yyyy-mm-dd'),pn,LOG_ACTION) x`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pnListTemp := make([]string, 0)
	var ioSummaryInfo IOSummaryInfo
	for rows.Next() {
		err = rows.Scan(
			&ioSummaryInfo.Data,
			&ioSummaryInfo.Pn,
			&ioSummaryInfo.Process,
			&ioSummaryInfo.Input,
			&ioSummaryInfo.Pass,
			&ioSummaryInfo.Bad,
		)
		if err != nil {
			return nil, err
		}
		pnListTemp = append(pnListTemp, ioSummaryInfo.Pn)
		ioSummaryInfoList = append(ioSummaryInfoList, ioSummaryInfo)
	}
	pnList := Utils.RemoveRepeatedElement(pnListTemp)
	r := make([]interface{}, 0)
	for i := 0; i < len(pnList); i++ {
		m := make(map[string]interface{})
		m["型号"] = pnList[i]
		m["调试-投入"] = 0
		m["调试-产出"] = 0
		m["常温-投入"] = 0
		m["常温-产出"] = 0
		m["EE-投入"] = 0
		m["EE-产出"] = 0
		m["入库-投入"] = 0
		m["入库-产出"] = 0
		if m["代码"], err = Utils.GetIOSummaryList(pnList[i]); err != nil {
			m["代码"] = ""
		}
		for _, ioSummaryInfoValue := range ioSummaryInfoList {
			if ioSummaryInfoValue.Pn == pnList[i] {
				switch ioSummaryInfoValue.Process {
				case "TUN_10G":
					m["调试-投入"] = ioSummaryInfoValue.Input
				case "TUNING":
					m["调试-产出"] = ioSummaryInfoValue.Pass
				case "TEST_10G":
					m["常温-投入"] = ioSummaryInfoValue.Input
				case "TESTING":
					m["常温-产出"] = ioSummaryInfoValue.Pass
				case "EEPROM CHECK":
					m["EE-投入"] = ioSummaryInfoValue.Input
					m["EE-产出"] = ioSummaryInfoValue.Pass
				case "FQC_INWH":
					m["入库-投入"] = ioSummaryInfoValue.Input
					m["入库-产出"] = ioSummaryInfoValue.Pass
				}
			}
		}
		m["调试-不良"] = reflect.ValueOf(m["调试-投入"]).Int() - reflect.ValueOf(m["调试-产出"]).Int()
		m["常温-不良"] = reflect.ValueOf(m["常温-投入"]).Int() - reflect.ValueOf(m["常温-产出"]).Int()
		m["EE-不良"] = reflect.ValueOf(m["EE-投入"]).Int() - reflect.ValueOf(m["EE-产出"]).Int()
		m["入库-不良"] = reflect.ValueOf(m["入库-投入"]).Int() - reflect.ValueOf(m["入库-产出"]).Int()
		r = append(r, m)
	}
	return r, nil
}
