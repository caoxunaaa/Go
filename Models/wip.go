package Models

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
)

type WipInfo struct {
	Pn      string
	Process string
	Pass    uint32
	Fail    uint32
	Total   uint32
}

func GetWipInfoList(pn string) (wipInfoList []WipInfo, err error) {
	startTime, endTime := Utils.GetCurrentAndZeroDayTime()
	sqlStr := `with TRX AS (select * from (SELECT distinct a.MANUFACTURE_GROUP,b.*,dense_rank()over(partition by b.sn order by b.action_time asc)zz,
dense_rank()over(partition by b.sn order by b.id DESC)rr,
c."sequence" as SEQ FROM superxon.autodt_process_log b,superxon.autodt_tracking a,superxon.workstage c where b.sn=a.bosa_sn 
and b.pn=a.partnumber and b.log_action = c."processname"
and (B.sn IN (select distinct D.SN from (SELECT distinct a.MANUFACTURE_GROUP,b.*,dense_rank()over(partition by b.sn order by b.action_time asc)zz,
dense_rank()over(partition by b.sn order by b.action_time DESC)rr,c."sequence" as SEQ
FROM superxon.autodt_process_log b,superxon.autodt_tracking a,superxon.workstage c
where b.sn=a.bosa_sn and b.pn=a.partnumber and b.log_action = c."processname"
and b.action_time between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss')and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss') 
and (b.pn = '` + pn + `'/*and b.log_action = 'TRACKING SN'*/))D WHERE D.RR = 1)))E
WHERE E.RR = 1 /*AND E.log_action <> 'FQC_INWH'*/)
select f.* from (select distinct g.PN as PN,/*g.SEQ as 序列,*/g.log_action as 工序,
sum(case g.p_value when 'PASS' then 1 else 0 end)over(partition by g.log_action,g.PN) PASS,
sum(case g.p_value when 'PASS' then 0 else 1 end)over(partition by g.log_action,g.PN) FAIL,
count(sn)over(partition by g.log_action,g.PN)总计        
from TRX g where g.rr = 1 /*and g.log_action = 'FQC_INWH'*/)f order by f.PN,f.总计 desc
`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	var wipInfo WipInfo
	for rows.Next() {
		err = rows.Scan(&wipInfo.Pn, &wipInfo.Process, &wipInfo.Pass, &wipInfo.Fail, &wipInfo.Total)
		if err != nil {
			return nil, err
		}
		wipInfoList = append(wipInfoList, wipInfo)
	}
	return
}
