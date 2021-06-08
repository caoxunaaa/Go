package OsaQaStatisticDisplay

import "SuperxonWebSite/Databases"

type OsaTransmitResp struct {
	InsName string
	Pn      string
	Pass    int
	Fail    int
	Total   int
}

func GetOsaTransmitter(startTime, endTime string) ([]*OsaTransmitResp, error) {
	sqlStr := `With OSA AS (select * from (select t.id,substr(t.PN,0,instr(t.pn,' ')-1) AS PN,
substr(t.sn,0,instr(t.sn,' ')-1) as sn,
substr(t.worknum,0,instr(t.worknum,' ')-1) AS worknum,
substr(t.insname,0,7) AS insname,t.testtime,
substr(t.status,0,instr(t.status,' ')-1) AS status,t.pwrzok,pwruok,t.pwrend,
(case when round(t.pwrzok*t.pwrend,1) > 0 then'1'else'0'end) as XX,
dense_rank()over(partition by t.sn order by t.id desc)rr 
from superxon.autodt_transmit_autocouple t where t.testtime BETWEEN TO_DATE('` + startTime + `','YYYY-MM-DD,HH24:MI:SS') AND TO_DATE('` + endTime + `','YYYY-MM-DD,HH24:MI:SS'))D)

select f.* from (select distinct g.insname,g.PN,
sum(case g.status when 'Pass' then 1 else 0 end)over(partition by g.PN,g.insname) PASS,
sum(case g.status when 'Pass' then 0 else 1 end)over(partition by g.PN,g.insname) FAIL,
count(sn)over(partition by g.PN,g.insname)总计        
from OSA g where g.rr = 1)f order by f.insname,f.总计 desc`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resp = make([]*OsaTransmitResp, 0)
	for rows.Next() {
		var res = new(OsaTransmitResp)
		err = rows.Scan(
			&res.InsName,
			&res.Pn,
			&res.Pass,
			&res.Fail,
			&res.Total)
		if err != nil {
			return nil, err
		}
		resp = append(resp, res)
	}
	return resp, nil
}
