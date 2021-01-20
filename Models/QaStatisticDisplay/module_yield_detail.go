package QaStatisticDisplay

import (
	"SuperxonWebSite/Databases"
	"database/sql"
	"fmt"
)

type QaPn struct {
	Pn string
}

func GetQaPnList(queryCondition *QueryCondition) (qaPnList []QaPn, err error) {
	sqlStr := `select distinct t.pn from superxon.autodt_process_log t where t.action_time between to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaPn *QaPn
	qaPn = new(QaPn)
	for rows.Next() {
		err = rows.Scan(
			&qaPn.Pn)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		qaPnList = append(qaPnList, *qaPn)
	}
	return
}

type QaStatisticInfo struct {
	Pn            string
	Sequence      string
	Process       string
	Version       string
	TotalInput    uint32
	FinalOk       uint32
	FinalBad      uint32
	FinalPassRate string
}

func GetQaStatisticInfoList(queryCondition *QueryCondition) (qaStatisticInfoList []QaStatisticInfo, err error) {
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
				and b.action_time between to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss')
				and to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss') and b.pn = '` + queryCondition.Pn + `')
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

func GetQaStatisticOrderInfoList(queryCondition *QueryCondition) (qaStatisticInfoList []QaStatisticInfo, err error) {
	sqlStr := `with TRX as (SELECT distinct a.MANUFACTURE_GROUP,d.LOT_TYPE,
			(case when substr(b.softversion,length(b.softversion)-4) like '%验证软件' then substr(b.softversion,0,length(b.softversion)-5)
			when substr(b.softversion,length(b.softversion)-1) LIKE '%*_'escape '*' then substr(b.softversion,0,length(b.softversion)-1) else B.SOFTVERSION END) as SVERSION,b.*,
			rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz,
			rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr,
			c."sequence" as SEQ
			FROM superxon.autodt_process_log b,(SELECT C.id,c.partnumber,c.manufacture_group,c.tosa_group,c.rosa_group,c.bosa_group,c.pcba1_group,c.bosa_sn,c.modifydate,c.la 
			FROM (select t.*,dense_rank()over(partition by T.PARTNUMBER,T.BOSA_SN order by T.MODIFYDATE DESC)LA from superxon.autodt_tracking t)C where C.LA=1) a,superxon.workstage c,
			(select t.partnumber,t.version,t.pch_tc, (case when  substr(t.pch_lx,0,10) like'TRX试生产产品工单%' then 'TRX正常品'
			when substr(t.pch_lx,0,10) like  'TRX量产产品工单%' then 'TRX正常品'  else 'TRX改制返工品' END) as LOT_TYPE,t.pch_lx
			from superxon.sgd_scdd_trx t) d
			where b.sn=a.bosa_sn and b.log_action = c."processname" and a.partnumber =b.pn and d.pch_tc=a.manufacture_group and b.pn=d.partnumber
			and b.action_time between to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss')
			and to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss') and b.pn like '%` + queryCondition.Pn + `%' AND D.LOT_TYPE LIKE '%` + queryCondition.WorkOrderType + `%')
			select e.* from (select distinct h.PN as PN,h.SEQ as 序列,h.log_action as 工序,h.SVERSION,
			count(sn)over(partition by h.log_action,h.PN,h.SVERSION)总输入,
			sum(case h.p_value when 'PASS' then 1 else 0 end)over(partition by h.log_action,h.PN,h.SVERSION)最终良品,
			sum(case h.p_value when 'PASS' then 0 else 1 end)over(partition by h.log_action,h.PN,h.SVERSION)最终不良品
			from TRX h where h.rr=1)e
			order by e.pn,e.序列 ASC`
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
			&qaStatisticInfo.Version,
			&qaStatisticInfo.TotalInput,
			&qaStatisticInfo.FinalOk,
			&qaStatisticInfo.FinalBad)
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
	Version        string
	ErrorCode      string
	ErrorCount     uint32
	ErrorRate      string
	ErrorInputRate string
}

func GetQaDefectsInfoList(queryCondition *QueryCondition) (qaDefectsInfoList []QaDefectsInfo, err error) {
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
				and b.pn = '` + queryCondition.Pn + `'
				and b.action_time between to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss') 
				and to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss'))x,superxon.autodt_results_ate_new y
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

func GetQaDefectsOrderInfoList(queryCondition *QueryCondition) (qaDefectsInfoList []QaDefectsInfo, err error) {
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
				and b.pn = '` + queryCondition.Pn + `' and D.LOT_TYPE = '` + queryCondition.WorkOrderType + `' 
				and b.action_time between to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss')
				and to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss'))x,superxon.autodt_results_ate_new y
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
			&qaDefectsInfo.Version,
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

type QaDefectsDetailInfo struct {
	WorkOrderId       string
	PchLx             string
	Pn                string
	Sn                string
	Operator          string
	LogAction         string
	ActionTime        string
	Parameter         string
	PValue            string
	Comments          string
	StationId         string
	SoftVersion       string
	RunTime           int
	RunCount          int
	CommitVer         string
	ModuleSn          sql.NullString
	Version           string
	Temperature       float64
	ErrorCode         string
	TxAOP             float64
	TxER              float64
	TxPPJ             float64
	TxTR              float64
	TxTF              float64
	MARGIN            float64
	TxCROSS           float64
	TxOffPOWER        float64
	OVERLOAD          float64
	SENSITIVITY       float64
	SDASSERT          float64
	SDDASSERT         float64
	PathPenalty       float64
	PeakWavelength    float64
	SIGMA             float64
	BandWidth         float64
	SMSR              float64
	TXI               float64
	RXI               float64
	TE                float64
	DacApc            float64
	DacMod            float64
	DacCross          float64
	DacApd            float64
	DacLos            float64
	InfoApc           sql.NullString
	InfoMod           sql.NullString
	InfoCross         sql.NullString
	InfoApd           sql.NullString
	InfoLos           sql.NullString
	A2Temperature     float64
	A2Vcc             float64
	A2Ibias           float64
	A2TxMon           float64
	CaseTemperature   float64
	TecTemperature    float64
	EAAbsorb          float64
	InfoEA            sql.NullString
	DacEA             float64
	CenterSensitivity float64
	AopStability01    float64
	AopStability02    float64
	AopStability03    float64
	AopStability04    float64
	AopStability05    float64
	AopStability06    float64
	AopStability07    float64
	AopStability08    float64
	AopStability09    float64
	AopStability10    float64
	AopStabilityDelta float64
	Current5V         float64
	SetVcc            float64
}

//func GetQaDefectsDetailByPn(queryCondition *QueryCondition) ()
