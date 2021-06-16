package OsaStatisticDisplay

import "SuperxonWebSite/Databases"

type OsaYieldWithStation struct {
	InsName string  `db:"insname"`
	OsaPn   string  `db:"osa_pn"`
	Rate    string  `db:"rate"`
	Pass    int     `db:"pass"`
	Fail    int     `db:"fail"`
	Grand   int     `db:"grand"`
	AvgZok  float64 `db:"avg_zok"`
	AvgUok  float64 `db:"avg_uok"`
	AvgEnd  float64 `db:"avg_end"`
	Yield   float64 `db:"yield"`
	WYield  float64 `db:"w_yield"`
}

//获取某段时间osaPn对应的工位良率信息
func GetOsaYieldWithStationByOsaPn(osaPn, startTime, endTime string) ([]OsaYieldWithStation, error) {
	sqlStr := `With OSA AS (select * from (select t.id,
(case when instr(t.pn,' ') <= '15' then substr(t.PN,0,instr(t.pn,' ')-1) else substr(t.PN,0,instr(t.pn,'-')-1) end)  AS OSA_PN,
(case when instr(t.pn,' ') <= '15' then '1G' else substr(PN,16,instr(PN,'G')-15) end)  AS RATE,
substr(t.sn,0,instr(t.sn,' ')-1) as sn,
substr(t.worknum,0,instr(t.worknum,' ')-1) AS OSA_PCH,
substr(t.insname,0,7) AS insname,t.testtime,
substr(t.status,0,instr(t.status,' ')-1) AS status,t.pwrzok,pwruok,t.pwrend,
(case when round(t.pwrzok*t.pwrend,1) > 0 then'1'else'0'end) as END_ZOK,
(case when round(t.pwruok*t.pwrend,1) > 0 then'1'else'0'end) as END_UOK,
dense_rank()over(partition by t.sn order by t.id desc)rr 
from superxon.autodt_transmit_autocouple t where t.testtime BETWEEN TO_DATE('` + startTime + `','YYYY-MM-DD,HH24:MI:SS') 
AND TO_DATE('` + endTime + `','YYYY-MM-DD,HH24:MI:SS'))D
where d.rr =1 AND D.OSA_PN LIKE '` + osaPn + `%'
AND END_ZOK<> '0'AND END_UOK<>'0')

select f.*,round(f.PASS/f.grand*100,2) as Yield,
(case when f.rate = '1G' then round(f.avg_end/f.avg_zok*100,2) else round(f.avg_end/f.avg_uok*100,2) end) as W_yield 
from (select distinct g.insname,g.OSA_PN,G.RATE,
sum(case g.status when 'Pass' then 1 else 0 end)over(partition by g.OSA_PN,g.RATE,g.insname) PASS,
sum(case g.status when 'Pass' then 0 else 1 end)over(partition by g.OSA_PN,g.RATE,g.insname) FAIL,
count(sn)over(partition by g.OSA_PN,g.RATE,g.insname) as grand,
round(avg(g.pwrzok)over(partition by g.OSA_PN,g.RATE,g.insname),2) as avg_zok,
round(avg(g.pwruok)over(partition by g.OSA_PN,g.RATE,g.insname),2) as avg_uok,
round(avg(g.pwrend)over(partition by g.OSA_PN,g.RATE,g.insname),2) as avg_end        
from OSA g where g.rr = 1)f order by f.insname,f.grand desc`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var p OsaYieldWithStation
	var res = make([]OsaYieldWithStation, 0)
	for rows.Next() {
		err = rows.Scan(
			&p.InsName,
			&p.OsaPn,
			&p.Rate,
			&p.Pass,
			&p.Fail,
			&p.Grand,
			&p.AvgZok,
			&p.AvgUok,
			&p.AvgEnd,
			&p.Yield,
			&p.WYield)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}
