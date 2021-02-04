package OsaQaStatisticDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Models/OsaRunDisplay"
)

type QaOsaStatisticInfoByWorkOrderId struct {
	WorkOrderId     string
	Process         string
	Seq             string
	Pn              string
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

func GetQaOsaStatisticInfoListByWorkOrderId(osaQueryCondition *OsaRunDisplay.OsaQueryCondition) (qaOsaStatisticInfoByWorkOrderIdList []QaOsaStatisticInfoByWorkOrderId, err error) {
	sqlStr := `with OSA as (select * from (select s.lotno,t.PN,(case when length(t.pn) <= '14' then 'BOSA' else (substr(T.pn,16,instr(T.pn,'G')))end) as 系列,
t.sn,t.errorcode,t.tc_flag,t.testdate as 时间,
dense_rank()over(partition by t.pn,t.sn,t.tc_flag order by t.testdate asc)rr,
dense_rank()over(partition by t.pn,t.sn,t.tc_flag order by t.testdate desc)zz,
h."sequence" as SEQ
from superxon.autodt_results_liv t,superxon.autodt_results_osa_tracking s,superxon.workstage h
where t.sn=s.sn and t.pn=s.pn and h."processname"=t.tc_flag
and (  s.lotno = '` + osaQueryCondition.WorkOrderId + `')
/*and testdate between to_date('&T1','YYYY-MM-DD HH24:MI:SS')
and to_date('&T2','YYYY-MM-DD HH24:MI:SS')*/ order by testdate desc)
union all
select * from (select  p.lotno,o.PN ,(case when length(o.pn) <= '14' then 'BOSA' else (substr(o.pn,16,instr(o.pn,'G')))end) as 系列,
o.sn,o.errorcode,o.tc_flag,to_date(o.apd_t_date,'yyyy-mm-dd hh24:mi:ss') as 时间,
dense_rank()over(partition by o.pn,o.sn,o.tc_flag order by to_date(o.apd_t_date,'yyyy-mm-dd hh24:mi:ss') asc)rr,
dense_rank()over(partition by o.pn,o.sn,o.tc_flag order by to_date(o.apd_t_date,'yyyy-mm-dd hh24:mi:ss') desc)zz,
h."sequence" as SEQ
from superxon.autodt_results_opticsdata o,superxon.autodt_results_osa_tracking p,superxon.workstage h
where o.pn=p.pn　and o.sn=p.sn and h."processname"=o.tc_flag
and   (  p.lotno = '` + osaQueryCondition.WorkOrderId + `')
/*and to_date(o.apd_t_date,'yyyy-mm-dd hh24:mi:ss') between to_date('&T1','YYYY-MM-DD HH24:MI:SS')
and to_date('&T2','YYYY-MM-DD HH24:MI:SS')*/ order by to_date(o.apd_t_date,'yyyy-mm-dd hh24:mi:ss') desc)
union all
select * from (select j.*,h."sequence" as SEQ from (select  y.lotno,x.PN ,
(case when length(x.pn) <= '14' then 'ROSA' else (substr(x.pn,16,instr(x.pn,'G')))end) as 系列,x.sn,x.errorcode,
(case when x.rosa_t_type = '1' THEN 'ROSASENS' ELSE 'Other' END) as TC_FLAG,to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') as 时间,
dense_rank()over(partition by x.pn,x.sn,x.rosa_t_type order by to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') asc)rr,
dense_rank()over(partition by x.pn,x.sn,x.rosa_t_type order by to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') desc)zz
from superxon.autodt_results_rosasens x ,superxon.autodt_results_osa_tracking y
where x.pn=y.pn　and x.sn=y.sn and  (  y.lotno = '` + osaQueryCondition.WorkOrderId + `')
/*and to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') between to_date('&T1','YYYY-MM-DD HH24:MI:SS')
and to_date('&T2','YYYY-MM-DD HH24:MI:SS')*/ order by to_date(x.rosa_t_date,'yyyy-mm-dd hh24:mi:ss') desc) j
join superxon.workstage h on h."processname" = j.TC_FLAG)
union all
select * from (select j.LOTNO,substr(J.PN,0,instr(j.pn,'-')-1) AS PN,(case when length(j.pn) <= '14' then 'TOSA' else (substr(j.pn,16,instr(j.pn,'G')))end) as 系列,
J.SN,J.ERRORCODE,J.TC_FLAG,J.时间,j.rr,j.zz, h."sequence" as SEQ from
(select distinct substr(t.sn,0,instr(t.sn,' ')-1) as sn,substr(t.worknum,0,instr(t.worknum,' ')-1) as LOTNO,
substr(t.pn,0,instr(t.pn,' ')-1) as PN,
(case when substr(t.status,0,4)='Pass' then 0 else 1 end )as errorcode,
(case when substr(t.flownum,0,6)='boxnum'then'TX_COUPLE'else'TX_COUPLE'end) as TC_FLAG,t.testtime as 时间,
dense_rank()over(partition by pn,t.sn,t.worknum order by testtime asc)rr,
dense_rank()over(partition by pn,t.sn,t.worknum order by testtime desc)zz
from superxon.autodt_transmit_autocouple t where (  T.WORKNUM = '` + osaQueryCondition.WorkOrderId + `')
/*and testtime between to_date('&t1','yyyy-mm-dd hh24-mi-ss')
and to_date('&t2','yyyy-mm-dd hh24-mi-ss')*/ order by testtime desc) j join superxon.workstage h on h."processname"=j.TC_FLAG)
union all
select * from (Select  j.LOTNO,substr(J.PN,0,instr(j.pn,'-')-1) AS PN,(case when length(j.pn) <= '14' then 'ROSA' else (substr(j.pn,16,instr(j.pn,'G')))end) as 系列,
J.SN,J.ERRORCODE,J.TC_FLAG,J.时间,j.rr,j.zz,h."sequence" as SEQ from
(select distinct t.sn as sn,t.worknum as LOTNO,
(case when Length(upper(t.pn)) > '19' then substr(UPPER(T.PN),0,instr(UPPER(T.PN),'G')) else UPPER(T.PN) end)as PN,
(case when substr(t.status,0,4)='PASS' then 0 else 1 end )as errorcode,
(case when substr(t.flownum,0,6)='0001' then 'RX_COUPLE' else'RX_COUPLE' END) as TC_FLAG,t.testtime as 时间,
dense_rank()over(partition by pn,sn,t.worknum order by testtime asc)rr,
dense_rank()over(partition by pn,sn,t.worknum order by testtime desc)zz
from superxon.autodt_recive_autocouple t
where ( T.WORKNUM = '` + osaQueryCondition.WorkOrderId + `')and T.SN <> 'NULL'/*AND testtime between to_date('&t1','yyyy-mm-dd hh24-mi-ss')
and to_date('&t2','yyyy-mm-dd hh24-mi-ss')*/ order by testtime desc) j join superxon.workstage h on h."processname"=J.TC_FLAG))

select b.*,round(b.一次良品/b.总投入*100,2)||'%'直通率,d.总输入,d.最终良品,d.最终不良品,
round(d.最终良品/d.总输入*100,2)||'%'最终良率,f.累计测试数,f.累计良品,f.累计不良品,round(f.累计良品/f.累计测试数*100,2)||'%'累计良率
from (select distinct a.lotno as 批次,a.tc_flag AS 工序,a.SEQ AS 序列,a.PN AS PN,
count(sn)over(partition by a.lotno,a.tc_flag,a.PN)总投入,
sum(case a.errorcode
when 0 then 1 else 0 end)over(partition by  a.lotno,a.tc_flag,a.PN)一次良品,
sum(case a.errorcode
when 0 then 0 else 1 end)over(partition by  a.lotno,a.tc_flag,a.PN)一次不良品
from OSA a where a.rr=1)b,
(select distinct c.lotno as 批次,c.tc_flag AS 工序,c.SEQ AS 序列,c.PN AS PN,
count(sn)over(partition by c.lotno,c.tc_flag,c.PN)总输入,
sum(case c.errorcode
when 0 then 1 else 0 end)over(partition by  c.lotno,c.tc_flag,c.PN)最终良品,
sum(case c.errorcode
when 0 then 0 else 1 end)over(partition by  c.lotno,c.tc_flag,c.PN)最终不良品
from OSA c where c.zz=1)d,
(select distinct e.lotno as 批次,e.tc_flag AS 工序,e.SEQ AS 序列,e.PN AS PN,
count(sn)over(partition by e.lotno,e.tc_flag,e.PN）累计测试数,
sum(case e.errorcode
when 0 then 1 else 0 end)over(partition by  e.lotno,e.tc_flag,e.PN)累计良品,
sum(case e.errorcode
when 0 then 0 else 1 end)over(partition by  e.lotno,e.tc_flag,e.PN)累计不良品
from OSA e )f
where b.工序=d.工序 and b.PN =d.PN and b.工序=f.工序 and b.PN=f.PN and b.序列=d.序列 and d.序列=f.序列 ORDER BY b.序列 asc`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaOsaStatisticInfoByWorkOrderId QaOsaStatisticInfoByWorkOrderId
	for rows.Next() {
		err = rows.Scan(
			&qaOsaStatisticInfoByWorkOrderId.WorkOrderId,
			&qaOsaStatisticInfoByWorkOrderId.Process,
			&qaOsaStatisticInfoByWorkOrderId.Seq,
			&qaOsaStatisticInfoByWorkOrderId.Pn,
			&qaOsaStatisticInfoByWorkOrderId.TotalInvestment,
			&qaOsaStatisticInfoByWorkOrderId.OnceOk,
			&qaOsaStatisticInfoByWorkOrderId.OnceBad,
			&qaOsaStatisticInfoByWorkOrderId.OncePassRate,
			&qaOsaStatisticInfoByWorkOrderId.TotalInput,
			&qaOsaStatisticInfoByWorkOrderId.FinalOk,
			&qaOsaStatisticInfoByWorkOrderId.FinalBad,
			&qaOsaStatisticInfoByWorkOrderId.FinalPassRate,
			&qaOsaStatisticInfoByWorkOrderId.AccTotalTest,
			&qaOsaStatisticInfoByWorkOrderId.AccOk,
			&qaOsaStatisticInfoByWorkOrderId.AccBad,
			&qaOsaStatisticInfoByWorkOrderId.AccPassRate)
		if err != nil {
			return nil, err
		}
		qaOsaStatisticInfoByWorkOrderIdList = append(qaOsaStatisticInfoByWorkOrderIdList, qaOsaStatisticInfoByWorkOrderId)
	}
	return
}
