package OsaRunDisplay

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

func GetStationStatus() (stationStatusList []StationStatus, err error) {
	startTime, endTime := Utils.GetCurrentAndZeroTime()
	sqlStr := `with osa as (select x.* from (SELECT distinct b.LOTNO,d.LOT_TYPE,
(case when substr(a.softversion,length(a.softversion)-4) like '%验证软件' then substr(a.softversion,0,length(a.softversion)-5)
when substr(a.softversion,length(a.softversion)-1) LIKE '%*_'escape '*' then substr(a.softversion,0,length(a.softversion)-1) else a.SOFTVERSION END) as SVERSION,a.*,
dense_rank()over(partition by a.sn,a.log_action order by a.action_time asc)zz,
dense_rank()over(partition by a.sn,a.log_action order by a.action_time DESC)rr,
c."sequence" as SEQ
FROM superxon.autodt_process_log a,(SELECT e.* FROM (select t.*,dense_rank()over(partition by T.PN,T.SN order by T.TXDATETIME DESC)LTX,
dense_rank()over(partition by T.PN,T.SN order by T.RXDATETIME DESC)LRX from superxon.AUTODT_RESULTS_OSA_TRACKING t)e where e.LTX=1) b,superxon.workstage c,
(select t.partnumber,t.version,t.pch_tc, (case when  substr(t.pch_lx,0,10) like'OSA试生产产品工单' then 'OSA正常品' 
when substr(t.pch_lx,0,10) like  'OSA量产产品工单' then 'OSA正常品' else'OSA改制返工品' END) as LOT_TYPE,t.pch_lx
from superxon.sgd_scdd_trx t)d
where b.sn=a.sn and a.log_action = c."processname" and a.pn =b.pn and d.pch_tc=b.LOTNO and a.pn=d.partnumber
and a.action_time between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') 
and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss'))x)
select d.* from (select distinct g.stationid,g.log_action as 工序,g.pn as 型号,
count(G.sn)over(partition by g.stationid,g.log_action,g.pn) 总生产数,sum(case g.p_value when 'PASS' then 1 else 0 end)over(partition by g.stationid,g.log_action,g.pn) 良品数量,
ROUND((sum(case g.p_value when 'PASS' then 1 else 0 end)over(partition by g.stationid,g.log_action,g.pn)/(count(G.sn)over(partition by g.stationid,g.log_action,g.pn))*100),2)||'%' 工位良率
from OSA g where g.RR=1)d order by d.stationid asc`
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
