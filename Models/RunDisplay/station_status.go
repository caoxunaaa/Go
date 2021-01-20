package RunDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
)

type StationStatus struct {
	StationId string
	LogAction string
	Pn        string
	TotalNum  uint16
	TotalPass uint16
	PassRate  string
}

// @title GetStationStatus
// @description 获取设备工位的状态信息和当天生产产品的对应良率
// @auth xun.cao
// @return []StationStatus
func GetStationStatus() (stationStatusList []StationStatus, err error) {
	startTime, endTime := Utils.GetCurrentAndZeroTime()
	//lotType := "TRX正常品"
	sqlStr := `with TRX AS (select y.errorcode,x.* from (SELECT distinct a.MANUFACTURE_GROUP,d.LOT_TYPE,
(case when substr(b.softversion,length(b.softversion)-4) like '%验证软件' then substr(b.softversion,0,length(b.softversion)-5) 
when substr(b.softversion,length(b.softversion)-1) LIKE '%*_'escape '*' then substr(b.softversion,0,length(b.softversion)-1) else B.SOFTVERSION END) as SVERSION,b.*,
dense_rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz,
dense_rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr,
c."sequence" as SEQ 
FROM superxon.autodt_process_log b,(SELECT C.id,c.partnumber,c.manufacture_group,c.tosa_group,c.rosa_group,c.bosa_group,c.pcba1_group,c.bosa_sn,c.modifydate,c.la 
FROM (select t.*,dense_rank()over(partition by T.PARTNUMBER,T.BOSA_SN order by T.MODIFYDATE DESC)LA from superxon.autodt_tracking t)C where C.LA=1) a,superxon.workstage c,
(select t.partnumber,t.version,t.pch_tc, (case when  substr(t.pch_lx,0,10) like'TRX试生产产品工单' then 'TRX正常品' 
when substr(t.pch_lx,0,10) like  'TRX量产产品工单' then 'TRX正常品'  else'TRX改制返工品' END) as LOT_TYPE,t.pch_lx
from superxon.sgd_scdd_trx t)d
where b.sn=a.bosa_sn and b.log_action = c."processname" and a.partnumber =b.pn and d.pch_tc=a.manufacture_group and b.pn=d.partnumber
and (/* b.pn LIKE  '&PN%' and*/ /*b.log_action like '&工序'*/ /*and b.stationid like '&Station%'*/ D.LOT_TYPE LIKE '%TRX%') 
and b.action_time between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') 
and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss'))x,superxon.autodt_results_ate_new y
where x.sn=y.opticssn and x.resultsid =y.id and x.rr = '1'/*and y.errorcode <> '0'*/)

select d.* from (select distinct g.stationid,g.log_action as 工序,g.pn as 型号,
count(G.sn)over(partition by g.stationid,g.log_action,g.pn) 总生产数,sum(case g.p_value when 'PASS' then 1 else 0 end)over(partition by g.stationid,g.log_action,g.pn) 良品数量,
ROUND((sum(case g.p_value when 'PASS' then 1 else 0 end)over(partition by g.stationid,g.log_action,g.pn)/（count(G.sn)over(partition by g.stationid,g.log_action,g.pn))*100),2)||'%' 工位良率
from TRX g where g.RR=1)d order by d.stationid asc`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stationStatus StationStatus
	for rows.Next() {
		err = rows.Scan(
			&stationStatus.StationId,
			&stationStatus.LogAction,
			&stationStatus.Pn,
			&stationStatus.TotalNum,
			&stationStatus.TotalPass,
			&stationStatus.PassRate)
		if err != nil {
			return nil, err
		}
		stationStatusList = append(stationStatusList, stationStatus)
	}
	return
}
