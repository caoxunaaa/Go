package OsaRunDisplay

import (
	"SuperxonWebSite/Databases"
)

type OsaWipInfo struct {
	Pn      string
	Seq     uint
	Process string
	Pass    uint32
	Fail    uint32
	Total   uint32
}

func GetOsaWipInfoList(osaQueryCondition *OsaQueryCondition) (osaWipInfoList []OsaWipInfo, err error) {
	sqlStr := `with OSA AS (select * from (SELECT distinct a.LOTNO,b.*,dense_rank()over(partition by b.sn order by b.action_time asc)FP,
dense_rank()over(partition by b.sn order by b.id DESC)LP,
c."sequence" as SEQ FROM superxon.autodt_process_log b,superxon.AUTODT_RESULTS_OSA_TRACKING a,superxon.workstage c where b.sn=a.SN 
and b.pn=a.PN and b.log_action = c."processname"
and (B.sn IN (select distinct D.SN from (SELECT distinct a.LOTNO,b.*,dense_rank()over(partition by b.sn order by b.action_time asc)FP,
dense_rank()over(partition by b.sn order by b.action_time DESC)LP,c."sequence" as SEQ
FROM superxon.autodt_process_log b,superxon.AUTODT_RESULTS_OSA_TRACKING a,superxon.workstage c
where b.sn=a.SN and b.pn=a.PN and b.log_action = c."processname"
and b.action_time between to_date('` + osaQueryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss')and to_date('` + osaQueryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss') 
and (b.pn like '` + osaQueryCondition.Pn + `'))D WHERE D.LP = 1)))E
WHERE E.LP = 1 )

select f.* from (select distinct g.PN as PN,g.SEQ as 序列,g.log_action as 工序,
sum(case g.p_value when 'PASS' then 1 else 0 end)over(partition by g.log_action,g.PN) PASS,
sum(case g.p_value when 'PASS' then 0 else 1 end)over(partition by g.log_action,g.PN) FAIL,
count(sn)over(partition by g.log_action,g.PN)总计        
from OSA g where g.LP = 1 )f order by /*f.PN,f.总计,*/f.序列 desc`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var osaWipInfo OsaWipInfo
	for rows.Next() {
		err = rows.Scan(
			&osaWipInfo.Pn,
			&osaWipInfo.Seq,
			&osaWipInfo.Process,
			&osaWipInfo.Pass,
			&osaWipInfo.Fail,
			&osaWipInfo.Total)
		if err != nil {
			return nil, err
		}
		osaWipInfoList = append(osaWipInfoList, osaWipInfo)
	}
	return
}
