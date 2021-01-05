package QaStatisticDisplay

import "SuperxonWebSite/Databases"

type QaPn struct {
	Pn string
}

func GetQaPnList(startTime string, endTime string) (qaPnList []QaPn, err error) {
	sqlStr := `select distinct t.pn from superxon.autodt_process_log t where t.action_time between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss')`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaPn QaPn
	for rows.Next() {
		err = rows.Scan(
			&qaPn.Pn)
		if err != nil {
			return nil, err
		}
		qaPnList = append(qaPnList, qaPn)
	}
	return
}

type QaStatisticInfo struct {
	Pn            string
	Sequence      string
	Process       string
	TotalInput    uint32
	FinalOk       uint32
	FinalBad      uint32
	FinalPassRate string
}

func GetQaStatisticInfoList(pn string, startTime string, endTime string) (qaStatisticInfoList []QaStatisticInfo, err error) {
	sqlStr := `with TRX as (SELECT distinct a.MANUFACTURE_GROUP,d.LOT_TYPE,b.*,
				rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz,
				rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr,
				c."sequence" as SEQ
				FROM superxon.autodt_process_log b,(SELECT aa.* FROM (select t.*,dense_rank()over(partition by T.PARTNUMBER,T.BOSA_SN order by T.MODIFYDATE DESC)LA
				from superxon.autodt_tracking t)aa where aa.LA=1) a,superxon.workstage c,
				(select t.partnumber,t.version,t.pch_tc, (case when  substr(t.pch_lx,0,10) like'TRX试生产产品工单%' then 'TRX正常品'
				when substr(t.pch_lx,0,10) like  'TRX量产产品工单%' then 'TRX正常品'  else'TRX改制返工品' END) as LOT_TYPE,t.pch_lx
				from superxon.sgd_scdd_trx t) d
				where b.sn=a.bosa_sn and b.log_action = c."processname" and a.partnumber =b.pn 
				and d.pch_tc=a.manufacture_group and b.pn=d.partnumber
				and b.action_time between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss')
				and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss') and b.pn = '` + pn + `')
				select e.pn,e.序列,e.工序,e.总输入,e.最终良品,e.最终不良品,round(e.最终良品/e.总输入*100,2)||'%' 最终良率
				from
				(select distinct h.PN as PN,h.SEQ as 序列,h.log_action as 工序 ,
				count(sn)over(partition by h.log_action,h.PN)总输入,
				sum(case h.p_value when 'PASS' then 1 else 0 end)over(partition by h.log_action,h.PN)最终良品,
				sum(case h.p_value when 'PASS' then 0 else 1 end)over(partition by h.log_action,h.PN)最终不良品
				from TRX h where h.rr=1)e order by e.pn,e.序列 ASC`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaStatisticInfo QaStatisticInfo
	for rows.Next() {
		err = rows.Scan(
			&qaStatisticInfo.Pn,
			&qaStatisticInfo.Sequence,
			&qaStatisticInfo.Process,
			&qaStatisticInfo.TotalInput,
			&qaStatisticInfo.FinalOk,
			&qaStatisticInfo.FinalBad,
			&qaStatisticInfo.FinalPassRate)
		if err != nil {
			return nil, err
		}
		qaStatisticInfoList = append(qaStatisticInfoList, qaStatisticInfo)
	}
	return
}

func GetQaStatisticOrderInfoList(pn string, startTime string, endTime string, order string) (qaStatisticInfoList []QaStatisticInfo, err error) {
	sqlStr := ``
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaStatisticInfo QaStatisticInfo
	for rows.Next() {
		err = rows.Scan(
			&qaStatisticInfo.Pn,
			&qaStatisticInfo.Sequence,
			&qaStatisticInfo.Process,
			&qaStatisticInfo.TotalInput,
			&qaStatisticInfo.FinalOk,
			&qaStatisticInfo.FinalBad,
			&qaStatisticInfo.FinalPassRate)
		if err != nil {
			return nil, err
		}
		qaStatisticInfoList = append(qaStatisticInfoList, qaStatisticInfo)
	}
	return
}

type QaDefectsInfo struct {
	Pn             string
	Sequence       string
	ErrorCode      string
	ErrorCount     uint32
	ErrorRate      string
	ErrorInputRate string
}

func GetQaDefectsInfoList(pn string, startTime string, endTime string) (qaDefectsInfoList []QaDefectsInfo, err error) {
	sqlStr := `with TRX AS (select y.errorcode,x.* from (SELECT distinct a.MANUFACTURE_GROUP,d.LOT_TYPE,b.*,
				rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz,
				rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr,
				c."sequence" as SEQ
				FROM superxon.autodt_process_log b,(SELECT C.id,c.partnumber,c.manufacture_group,c.tosa_group,c.rosa_group,c.bosa_group,c.pcba1_group,c.bosa_sn,c.modifydate,c.la 
				FROM (select t.*,dense_rank()over(partition by T.PARTNUMBER,T.BOSA_SN order by T.MODIFYDATE DESC)LA from superxon.autodt_tracking t)C where C.LA=1) a,superxon.workstage c,
				(select t.partnumber,t.version,t.pch_tc, (case when  substr(t.pch_lx,0,10) like'TRX试生产产品工单%' then 'TRX正常品' 
				when substr(t.pch_lx,0,10) like  'TRX量产产品工单%' then 'TRX正常品'  else'TRX改制返工品' END) as LOT_TYPE,t.pch_lx
				from superxon.sgd_scdd_trx t)d
				where b.sn=a.bosa_sn and b.log_action = c."processname" and a.partnumber =b.pn and d.pch_tc=a.manufacture_group and b.pn=d.partnumber
				and b.pn = '` + pn + `'
				and b.action_time between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') 
				and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss'))x,superxon.autodt_results_ate_new y
				where x.sn=y.opticssn and x.resultsid =y.id and y.errorcode <> '0')              
				select d.* from (select distinct g.PN as PN,g.log_action as 工序,g.ERRORCODE,
				count(G.sn)over(partition by g.log_action,g.ERRORCODE)不良数量,
				ROUND((count(G.sn)over(partition by g.ERRORCODE)/(sum(case g.p_value when 'FAIL' then 1 else 0 end)over(partition by g.PN))*100),2)||'%' 不良比重
				from TRX g where g.rr=1)d`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaDefectsInfo QaDefectsInfo
	for rows.Next() {
		err = rows.Scan(
			&qaDefectsInfo.Pn,
			&qaDefectsInfo.Sequence,
			&qaDefectsInfo.ErrorCode,
			&qaDefectsInfo.ErrorCount,
			&qaDefectsInfo.ErrorRate)
		if err != nil {
			return nil, err
		}
		qaDefectsInfoList = append(qaDefectsInfoList, qaDefectsInfo)
	}
	return
}

func GetQaDefectsOrderInfoList(pn string, startTime string, endTime string, order string) (qaDefectsInfoList []QaDefectsInfo, err error) {
	sqlStr := `with TRX AS (select y.errorcode,x.* from (SELECT distinct a.MANUFACTURE_GROUP,d.LOT_TYPE,
				(case when substr(b.softversion,length(b.softversion)-4) like '%验证软件' then substr(b.softversion,0,length(b.softversion)-5)
				when substr(b.softversion,length(b.softversion)-1) LIKE '%*_'escape '*' then substr(b.softversion,0,length(b.softversion)-1) else B.SOFTVERSION END) as SVERSION,b.*,
				dense_rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz,
				dense_rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr,
				c."sequence" as SEQ
				FROM superxon.autodt_process_log b,(SELECT C.id,c.partnumber,c.manufacture_group,c.tosa_group,c.rosa_group,c.bosa_group,c.pcba1_group,c.bosa_sn,c.modifydate,c.la 
				FROM (select t.*,dense_rank()over(partition by T.PARTNUMBER,T.BOSA_SN order by T.MODIFYDATE DESC)LA from superxon.autodt_tracking t)C where C.LA=1) a,superxon.workstage c,
				(select t.partnumber,t.version,t.pch_tc, (case when  substr(t.pch_lx,0,10) like'TRX试生产产品工单%' then 'TRX正常品'
				when substr(t.pch_lx,0,10) like  'TRX量产产品工单%' then 'TRX正常品'  else'TRX改制返工品' END) as LOT_TYPE,t.pch_lx
				from superxon.sgd_scdd_trx t)d
				where b.sn=a.bosa_sn and b.log_action = c."processname" and a.partnumber =b.pn and d.pch_tc=a.manufacture_group and b.pn=d.partnumber
				and b.pn = '` + pn + `' and D.LOT_TYPE = '` + order + `' 
				and b.action_time between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss')
				and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss'))x,superxon.autodt_results_ate_new y
				where x.sn=y.opticssn and x.resultsid =y.id and y.errorcode <> '0')
										
				select d.* from (select distinct g.PN as PN,g.log_action as 工序,G.SVERSION,g.ERRORCODE,
				count(G.sn)over(partition by g.ERRORCODE,G.SVERSION,g.log_action)不良数量,
				ROUND(（count(G.sn)over(partition by g.ERRORCODE)/(sum(case g.p_value when 'FAIL' then 1 else 0 end)over(partition by g.PN))*100),2)||'%' 不良比重
				from TRX g where g.RR=1)d
				order by d.不良数量 DESC`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaDefectsInfo QaDefectsInfo
	for rows.Next() {
		err = rows.Scan(
			&qaDefectsInfo.Pn,
			&qaDefectsInfo.Sequence,
			&qaDefectsInfo.ErrorCode,
			&qaDefectsInfo.ErrorCount,
			&qaDefectsInfo.ErrorRate)
		if err != nil {
			return nil, err
		}
		qaDefectsInfoList = append(qaDefectsInfoList, qaDefectsInfo)
	}
	return
}
