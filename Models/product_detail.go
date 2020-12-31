package Models

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Product struct {
	Name sql.NullString
}

func GetModuleList(startTime string, endTime string) (moduleList []Product, err error) {
	reBytes, _ := redis.Bytes(Databases.RedisConn.Do("get", "moduleList"))
	_ = json.Unmarshal(reBytes, &moduleList)
	if len(moduleList) != 0 {
		fmt.Println("使用redis")
		return
	}
	sqlStr := `select distinct t.partnumber from superxon.sgd_scdd_trx t where t.pch_tc_date between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss')`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	var pn Product
	for rows.Next() {
		err = rows.Scan(&pn.Name)
		if err != nil {
			return nil, err
		}
		if pn.Name.Valid == true {
			moduleList = append(moduleList, pn)
		}
	}
	datas, _ := json.Marshal(moduleList)
	_, _ = Databases.RedisConn.Do("SET", "moduleList", datas)
	_, err = Databases.RedisConn.Do("expire", "moduleList", 60*60)

	return
}

func GetOsaList(startTime string, endTime string) (osaList []Product, err error) {
	rebytes, _ := redis.Bytes(Databases.RedisConn.Do("get", "osaList"))
	_ = json.Unmarshal(rebytes, &osaList)
	if len(osaList) != 0 {
		fmt.Println("使用redis")
		return
	}
	rows, err := Databases.OracleDB.Query("select distinct t.pn from superxon.autodt_results_liv t where t.testdate between to_date('" + startTime + "','yyyy-mm-dd hh24:mi:ss') and to_date('" + endTime + "','yyyy-mm-dd hh24:mi:ss')")
	if err != nil {
		return nil, err
	}
	var osa Product
	for rows.Next() {
		err = rows.Scan(&osa.Name)
		if err != nil {
			return nil, err
		}
		if osa.Name.Valid == true {
			osaList = append(osaList, osa)
		}
	}
	datas, _ := json.Marshal(osaList)
	_, _ = Databases.RedisConn.Do("SET", "osaList", datas)
	_, err = Databases.RedisConn.Do("expire", "osaList", 60*60)
	return
}

type ProductInfo struct {
	Pn              string
	Sequence        string
	Process         string
	TotalInvestment int
	OnceOk          int
	OnceBad         int
	OncePassRate    string
	TotalInput      int
	FinalOk         int
	FinalBad        int
	FinalPassRate   string
	AccTotalTest    int
	AccOk           int
	AccBad          int
	AccPassRate     string
}

func GetModuleInfoList(product string, startTime string, endTime string) (moduleInfoList []ProductInfo, err error) {
	sqlStr := `with TRX as (SELECT distinct b.*, 
dense_rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz, 
dense_rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr, 
c."sequence" as SEQ FROM superxon.autodt_process_log b,superxon.autodt_tracking a,
superxon.workstage c where b.sn=a.bosa_sn and b.log_action = c."processname" and a.partnumber =b.pn and ( b.pn='` + product + `')
and b.action_time between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss')) 
select d.*,round (d.一次良品/d.总投入*100,2)||'%' 直通率,e.总输入,e.最终良品,e.最终不良品,
round(e.最终良品/e.总输入*100,2)||'%' 最终良率,
f.累计测试总数,f.累计良品,f.累计不良品,round(f.累计良品/f.累计测试总数*100,2)||'%' 累计良率 from (select distinct g.PN as PN,g.SEQ as 序列,g.log_action as 工序, 
count(sn)over(partition by g.log_action,g.PN)总投入, 
sum(case g.p_value when 'PASS' then 1 else 0 end)over(partition by g.log_action,g.PN)一次良品,
sum(case g.p_value when 'PASS' then 0 else 1 end)over(partition by g.log_action,g.PN)一次不良品 from TRX g where g.zz=1)d, 
(select distinct h.PN as PN,h.SEQ as 序列,h.log_action as 工序, 
count(sn)over(partition by h.log_action,h.PN)总输入, 
sum(case h.p_value when 'PASS' then 1 else 0 end)over(partition by h.log_action,h.PN)最终良品, 
sum(case h.p_value when 'PASS' then 0 else 1 end)over(partition by h.log_action,h.PN)最终不良品 from TRX h where h.rr=1)e, 
(select distinct m.PN as PN,m.SEQ as 序列,m.log_action as 工序, 
count(sn)over(partition by m.log_action,m.PN)累计测试总数, 
sum(case m.p_value when 'PASS' then 1 else 0 end)over(partition by m.log_action,m.PN)累计良品, 
sum(case m.p_value when 'PASS' then 0 else 1 end)over(partition by m.log_action,m.PN)累计不良品 from TRX m )f where d.工序=e.工序 and d.PN=e.PN and f.pn=e.pn and d.工序= f.工序 and d.序列=e.序列 and e.序列= f.序列 order by d.序列 ASC`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	var moduleinfo ProductInfo
	for rows.Next() {
		err = rows.Scan(
			&moduleinfo.Pn,
			&moduleinfo.Sequence,
			&moduleinfo.Process,
			&moduleinfo.TotalInvestment,
			&moduleinfo.OnceOk,
			&moduleinfo.OnceBad,
			&moduleinfo.OncePassRate,
			&moduleinfo.TotalInput,
			&moduleinfo.FinalOk,
			&moduleinfo.FinalBad,
			&moduleinfo.FinalPassRate,
			&moduleinfo.AccTotalTest,
			&moduleinfo.AccOk,
			&moduleinfo.AccBad,
			&moduleinfo.AccPassRate)
		if err != nil {
			return nil, err
		}
		moduleInfoList = append(moduleInfoList, moduleinfo)
	}
	return
}

func GetOsaInfoList(product string, startTime string, endTime string) (osaInfoList []ProductInfo, err error) {
	sqlStr := `with OSA as (select * from (select s.lotno,t.PN ,t.sn,t.errorcode,t.tc_flag,t.testdate as 时间, 
dense_rank()over(partition by t.pn,t.sn,t.tc_flag order by t.testdate asc)rr, 
dense_rank()over(partition by t.pn,t.sn,t.tc_flag order by t.testdate desc)zz, 
h."sequence" as SEQ from superxon.autodt_results_liv t,superxon.autodt_results_osa_tracking s,
superxon.workstage h where t.sn=s.sn and t.pn=s.pn and h."processname"=t.tc_flag and (t.pn = '` + product + `') 
and testdate between to_date('` + startTime + `','YYYY-MM-DD HH24:MI:SS') and to_date('` + endTime + `','YYYY-MM-DD HH24:MI:SS') order by testdate desc) 
union all 
select * from (select  p.lotno,o.PN ,o.sn,o.errorcode,o.tc_flag,to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') as 时间, 
dense_rank()over(partition by o.pn,o.sn,o.tc_flag order by to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') asc)rr, 
dense_rank()over(partition by o.pn,o.sn,o.tc_flag order by to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') desc)zz, 
h."sequence" as SEQ from superxon.autodt_results_opticsdata o,superxon.autodt_results_osa_tracking p,superxon.workstage h 
where o.pn=p.pn　and o.sn=p.sn and h."processname"=o.tc_flag and (o.pn = '` + product + `') 
and to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') between to_date('` + startTime + `','YYYY-MM-DD HH24:MI:SS') 
and to_date('` + endTime + `','YYYY-MM-DD HH24:MI:SS') order by to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') desc) 
union all 
select * from (select j.*,h."sequence" as SEQ from (select  y.lotno,x.PN ,x.sn,x.errorcode, 
(case when x.rosa_t_type = '1' THEN 'ROSASENS' ELSE 'Other' END) as TC_FLAG,to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') as 时间, 
dense_rank()over(partition by x.pn,x.sn,x.rosa_t_type order by to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') asc)rr, 
dense_rank()over(partition by x.pn,x.sn,x.rosa_t_type order by to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') desc)zz 
from superxon.autodt_results_rosasens x ,superxon.autodt_results_osa_tracking y 
where x.pn=y.pn and x.sn=y.sn and  ( x.pn = '` + product + `') 
and to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') between to_date('` + startTime + `','YYYY-MM-DD HH24:MI:SS') 
and to_date('` + endTime + `','YYYY-MM-DD HH24:MI:SS') order by to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') desc) j 
join superxon.workstage h on h."processname" = j.TC_FLAG) 
union all 
select * from (select j.*, h."sequence" as SEQ from (select substr(t.worknum,0,8) as LOTNO,substr(t.pn,0,14) as PN,substr(t.sn,0,13) as sn, 
(case when substr(t.status,0,4)='Pass' then 0 else 1 end )as errorcode, 
(case when substr(t.flownum,0,6)='boxnum' then 'TX_COUPLE' else 'TX_COUPLE' end) as TC_FLAG,t.testtime as 时间, 
dense_rank()over(partition by t.pn,t.sn,t.worknum order by testtime asc)rr, 
dense_rank()over(partition by t.pn,t.sn,t.worknum order by testtime desc)zz 
from superxon.autodt_transmit_autocouple t where ( t.pn = '` + product + `') 
and testtime between to_date('` + startTime + `','yyyy-mm-dd hh24-mi-ss') 
and to_date('` + endTime + `','yyyy-mm-dd hh24-mi-ss') order by testtime desc) j 
join superxon.workstage h on h."processname"=j.TC_FLAG) 
union all 
select * from (Select J.* ,h."sequence" as SEQ from (select substr(t.worknum,0,8) as LOTNO,substr(t.pn,0,14) as PN,substr(t.sn,0,13) as sn, 
(case when substr(t.status,0,4)='Pass' then 0 else 1 end )as errorcode, 
(case when substr(t.flownum,0,6)='0001' then 'RX_COUPLE' else 'RX_COUPLE' END) as TC_FLAG,t.testtime as 时间, 
dense_rank()over(partition by t.pn,t.sn,t.worknum order by testtime asc)rr, 
dense_rank()over(partition by t.pn,t.sn,t.worknum order by testtime desc)zz 
from superxon.autodt_recive_autocouple t 
where ( t.pn = '` + product + `')and testtime between to_date('` + startTime + `','yyyy-mm-dd hh24-mi-ss') 
and to_date('` + endTime + `','yyyy-mm-dd hh24-mi-ss') order by testtime desc) j 
join superxon.workstage h on h."processname"=J.TC_FLAG)) 
select b.*,round(b.一次良品/b.总投入*100,2)||'%'直通率,d.总输入,d.最终良品,d.最终不良品, 
round(d.最终良品/d.总输入*100,2)||'%'最终良率,f.累计测试数,f.累计良品,f.累计不良品,round(f.累计良品/f.累计测试数*100,2)||'%' 累计良率 
from (select distinct a.PN AS PN,a.SEQ AS 序列,a.tc_flag AS 工序, 
count(sn)over(partition by a.tc_flag,a.PN)总投入, 
sum(case a.errorcode 
when 0 then 1 else 0 end)over(partition by a.tc_flag,a.PN)一次良品, 
sum(case a.errorcode 
when 0 then 0 else 1 end)over(partition by a.tc_flag,a.PN)一次不良品 
from OSA a where a.rr=1)b, 
(select distinct c.PN AS PN,c.SEQ AS 序列,c.tc_flag AS 工序, 
count(sn)over(partition by c.tc_flag,c.PN)总输入, 
sum(case c.errorcode when 0 then 1 else 0 end)over(partition by c.tc_flag,c.PN)最终良品, 
sum(case c.errorcode when 0 then 0 else 1 end)over(partition by c.tc_flag,c.PN)最终不良品 
from OSA c where c.zz=1)d, 
(select distinct e.PN AS PN,e.SEQ AS 序列,e.tc_flag AS 工序, 
count(sn)over(partition by e.tc_flag,e.PN)累计测试数, 
sum(case e.errorcode when 0 then 1 else 0 end)over(partition by e.tc_flag,e.PN)累计良品, 
sum(case e.errorcode when 0 then 0 else 1 end)over(partition by e.tc_flag,e.PN)累计不良品 
from OSA e )f where b.工序=d.工序 and b.工序=f.工序  and b.序列=d.序列 and d.序列=f.序列 ORDER BY b.序列  asc`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	var osainfo ProductInfo
	for rows.Next() {
		err = rows.Scan(
			&osainfo.Pn,
			&osainfo.Sequence,
			&osainfo.Process,
			&osainfo.TotalInvestment,
			&osainfo.OnceOk,
			&osainfo.OnceBad,
			&osainfo.OncePassRate,
			&osainfo.TotalInput,
			&osainfo.FinalOk,
			&osainfo.FinalBad,
			&osainfo.FinalPassRate,
			&osainfo.AccTotalTest,
			&osainfo.AccOk,
			&osainfo.AccBad,
			&osainfo.AccPassRate)
		if err != nil {
			return nil, err
		}
		osaInfoList = append(osaInfoList, osainfo)
	}
	return
}

func GetYesterdayModuleInfoList(product string) (moduleInfoList []ProductInfo, err error) {
	startTime, endTime := Utils.GetAgoAndCurrentTime(Utils.Ago{Years: 0, Months: 0, Days: -1})
	sqlStr := `with TRX as (SELECT distinct b.*, 
dense_rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz, 
dense_rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr, 
c."sequence" as SEQ FROM superxon.autodt_process_log b,superxon.autodt_tracking a,
superxon.workstage c where b.sn=a.bosa_sn and b.log_action = c."processname" and a.partnumber =b.pn and ( b.pn='` + product + `')
and b.action_time between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss')) 
select d.*,round (d.一次良品/d.总投入*100,2)||'%' 直通率,e.总输入,e.最终良品,e.最终不良品,
round(e.最终良品/e.总输入*100,2)||'%' 最终良率,
f.累计测试总数,f.累计良品,f.累计不良品,round(f.累计良品/f.累计测试总数*100,2)||'%' 累计良率 from (select distinct g.PN as PN,g.SEQ as 序列,g.log_action as 工序, 
count(sn)over(partition by g.log_action,g.PN)总投入, 
sum(case g.p_value when 'PASS' then 1 else 0 end)over(partition by g.log_action,g.PN)一次良品,
sum(case g.p_value when 'PASS' then 0 else 1 end)over(partition by g.log_action,g.PN)一次不良品 from TRX g where g.zz=1)d, 
(select distinct h.PN as PN,h.SEQ as 序列,h.log_action as 工序, 
count(sn)over(partition by h.log_action,h.PN)总输入, 
sum(case h.p_value when 'PASS' then 1 else 0 end)over(partition by h.log_action,h.PN)最终良品, 
sum(case h.p_value when 'PASS' then 0 else 1 end)over(partition by h.log_action,h.PN)最终不良品 from TRX h where h.rr=1)e, 
(select distinct m.PN as PN,m.SEQ as 序列,m.log_action as 工序, 
count(sn)over(partition by m.log_action,m.PN)累计测试总数, 
sum(case m.p_value when 'PASS' then 1 else 0 end)over(partition by m.log_action,m.PN)累计良品, 
sum(case m.p_value when 'PASS' then 0 else 1 end)over(partition by m.log_action,m.PN)累计不良品 from TRX m )f where d.工序=e.工序 and d.PN=e.PN and f.pn=e.pn and d.工序= f.工序 and d.序列=e.序列 and e.序列= f.序列 order by d.序列 ASC`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	var moduleinfo ProductInfo
	for rows.Next() {
		err = rows.Scan(
			&moduleinfo.Pn,
			&moduleinfo.Sequence,
			&moduleinfo.Process,
			&moduleinfo.TotalInvestment,
			&moduleinfo.OnceOk,
			&moduleinfo.OnceBad,
			&moduleinfo.OncePassRate,
			&moduleinfo.TotalInput,
			&moduleinfo.FinalOk,
			&moduleinfo.FinalBad,
			&moduleinfo.FinalPassRate,
			&moduleinfo.AccTotalTest,
			&moduleinfo.AccOk,
			&moduleinfo.AccBad,
			&moduleinfo.AccPassRate)
		if err != nil {
			return nil, err
		}
		moduleInfoList = append(moduleInfoList, moduleinfo)
	}
	return
}

func GetYesterdayOsaInfoList(product string) (osaInfoList []ProductInfo, err error) {
	startTime, endTime := Utils.GetAgoAndCurrentTime(Utils.Ago{Years: 0, Months: 0, Days: -1})
	sqlStr := `with OSA as (select * from (select s.lotno,t.PN ,t.sn,t.errorcode,t.tc_flag,t.testdate as 时间, 
dense_rank()over(partition by t.pn,t.sn,t.tc_flag order by t.testdate asc)rr, 
dense_rank()over(partition by t.pn,t.sn,t.tc_flag order by t.testdate desc)zz, 
h."sequence" as SEQ from superxon.autodt_results_liv t,superxon.autodt_results_osa_tracking s,
superxon.workstage h where t.sn=s.sn and t.pn=s.pn and h."processname"=t.tc_flag and (t.pn = '` + product + `') 
and testdate between to_date('` + startTime + `','YYYY-MM-DD HH24:MI:SS') and to_date('` + endTime + `','YYYY-MM-DD HH24:MI:SS') order by testdate desc) 
union all 
select * from (select  p.lotno,o.PN ,o.sn,o.errorcode,o.tc_flag,to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') as 时间, 
dense_rank()over(partition by o.pn,o.sn,o.tc_flag order by to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') asc)rr, 
dense_rank()over(partition by o.pn,o.sn,o.tc_flag order by to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') desc)zz, 
h."sequence" as SEQ from superxon.autodt_results_opticsdata o,superxon.autodt_results_osa_tracking p,superxon.workstage h 
where o.pn=p.pn　and o.sn=p.sn and h."processname"=o.tc_flag and (o.pn = '` + product + `') 
and to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') between to_date('` + startTime + `','YYYY-MM-DD HH24:MI:SS') 
and to_date('` + endTime + `','YYYY-MM-DD HH24:MI:SS') order by to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') desc) 
union all 
select * from (select j.*,h."sequence" as SEQ from (select  y.lotno,x.PN ,x.sn,x.errorcode, 
(case when x.rosa_t_type = '1' THEN 'ROSASENS' ELSE 'Other' END) as TC_FLAG,to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') as 时间, 
dense_rank()over(partition by x.pn,x.sn,x.rosa_t_type order by to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') asc)rr, 
dense_rank()over(partition by x.pn,x.sn,x.rosa_t_type order by to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') desc)zz 
from superxon.autodt_results_rosasens x ,superxon.autodt_results_osa_tracking y 
where x.pn=y.pn and x.sn=y.sn and  ( x.pn = '` + product + `') 
and to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') between to_date('` + startTime + `','YYYY-MM-DD HH24:MI:SS') 
and to_date('` + endTime + `','YYYY-MM-DD HH24:MI:SS') order by to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') desc) j 
join superxon.workstage h on h."processname" = j.TC_FLAG) 
union all 
select * from (select j.*, h."sequence" as SEQ from (select substr(t.worknum,0,8) as LOTNO,substr(t.pn,0,14) as PN,substr(t.sn,0,13) as sn, 
(case when substr(t.status,0,4)='Pass' then 0 else 1 end )as errorcode, 
(case when substr(t.flownum,0,6)='boxnum' then 'TX_COUPLE' else 'TX_COUPLE' end) as TC_FLAG,t.testtime as 时间, 
dense_rank()over(partition by t.pn,t.sn,t.worknum order by testtime asc)rr, 
dense_rank()over(partition by t.pn,t.sn,t.worknum order by testtime desc)zz 
from superxon.autodt_transmit_autocouple t where ( t.pn = '` + product + `') 
and testtime between to_date('` + startTime + `','yyyy-mm-dd hh24-mi-ss') 
and to_date('` + endTime + `','yyyy-mm-dd hh24-mi-ss') order by testtime desc) j 
join superxon.workstage h on h."processname"=j.TC_FLAG) 
union all 
select * from (Select J.* ,h."sequence" as SEQ from (select substr(t.worknum,0,8) as LOTNO,substr(t.pn,0,14) as PN,substr(t.sn,0,13) as sn, 
(case when substr(t.status,0,4)='Pass' then 0 else 1 end )as errorcode, 
(case when substr(t.flownum,0,6)='0001' then 'RX_COUPLE' else 'RX_COUPLE' END) as TC_FLAG,t.testtime as 时间, 
dense_rank()over(partition by t.pn,t.sn,t.worknum order by testtime asc)rr, 
dense_rank()over(partition by t.pn,t.sn,t.worknum order by testtime desc)zz 
from superxon.autodt_recive_autocouple t 
where ( t.pn = '` + product + `')and testtime between to_date('` + startTime + `','yyyy-mm-dd hh24-mi-ss') 
and to_date('` + endTime + `','yyyy-mm-dd hh24-mi-ss') order by testtime desc) j 
join superxon.workstage h on h."processname"=J.TC_FLAG)) 
select b.*,round(b.一次良品/b.总投入*100,2)||'%'直通率,d.总输入,d.最终良品,d.最终不良品, 
round(d.最终良品/d.总输入*100,2)||'%'最终良率,f.累计测试数,f.累计良品,f.累计不良品,round(f.累计良品/f.累计测试数*100,2)||'%' 累计良率 
from (select distinct a.PN AS PN,a.SEQ AS 序列,a.tc_flag AS 工序, 
count(sn)over(partition by a.tc_flag,a.PN)总投入, 
sum(case a.errorcode 
when 0 then 1 else 0 end)over(partition by a.tc_flag,a.PN)一次良品, 
sum(case a.errorcode 
when 0 then 0 else 1 end)over(partition by a.tc_flag,a.PN)一次不良品 
from OSA a where a.rr=1)b, 
(select distinct c.PN AS PN,c.SEQ AS 序列,c.tc_flag AS 工序, 
count(sn)over(partition by c.tc_flag,c.PN)总输入, 
sum(case c.errorcode when 0 then 1 else 0 end)over(partition by c.tc_flag,c.PN)最终良品, 
sum(case c.errorcode when 0 then 0 else 1 end)over(partition by c.tc_flag,c.PN)最终不良品 
from OSA c where c.zz=1)d, 
(select distinct e.PN AS PN,e.SEQ AS 序列,e.tc_flag AS 工序, 
count(sn)over(partition by e.tc_flag,e.PN)累计测试数, 
sum(case e.errorcode when 0 then 1 else 0 end)over(partition by e.tc_flag,e.PN)累计良品, 
sum(case e.errorcode when 0 then 0 else 1 end)over(partition by e.tc_flag,e.PN)累计不良品 
from OSA e )f where b.工序=d.工序 and b.工序=f.工序  and b.序列=d.序列 and d.序列=f.序列 ORDER BY b.序列  asc`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	var osainfo ProductInfo
	for rows.Next() {
		err = rows.Scan(
			&osainfo.Pn,
			&osainfo.Sequence,
			&osainfo.Process,
			&osainfo.TotalInvestment,
			&osainfo.OnceOk,
			&osainfo.OnceBad,
			&osainfo.OncePassRate,
			&osainfo.TotalInput,
			&osainfo.FinalOk,
			&osainfo.FinalBad,
			&osainfo.FinalPassRate,
			&osainfo.AccTotalTest,
			&osainfo.AccOk,
			&osainfo.AccBad,
			&osainfo.AccPassRate)
		if err != nil {
			return nil, err
		}
		osaInfoList = append(osaInfoList, osainfo)
	}
	return
}
