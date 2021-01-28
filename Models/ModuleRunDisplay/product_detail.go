package ModuleRunDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
	"database/sql"
)

type Product struct {
	Name sql.NullString
}

func GetModuleList(startTime string, endTime string) (moduleList []Product, err error) {
	sqlStr := `select distinct t.partnumber from superxon.sgd_scdd_trx t where t.partnumber LIKE 'SO%' and t.pch_tc_date between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss')`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
	defer rows.Close()
	var moduleInfo ProductInfo
	for rows.Next() {
		err = rows.Scan(
			&moduleInfo.Pn,
			&moduleInfo.Sequence,
			&moduleInfo.Process,
			&moduleInfo.TotalInvestment,
			&moduleInfo.OnceOk,
			&moduleInfo.OnceBad,
			&moduleInfo.OncePassRate,
			&moduleInfo.TotalInput,
			&moduleInfo.FinalOk,
			&moduleInfo.FinalBad,
			&moduleInfo.FinalPassRate,
			&moduleInfo.AccTotalTest,
			&moduleInfo.AccOk,
			&moduleInfo.AccBad,
			&moduleInfo.AccPassRate)
		if err != nil {
			return nil, err
		}
		moduleInfoList = append(moduleInfoList, moduleInfo)
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
	defer rows.Close()
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
