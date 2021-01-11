package QaStatisticDisplay

import (
	"SuperxonWebSite/Databases"
	"database/sql"
)

type QueryCondition struct {
	Pn            string //产品型号
	WorkOrderId   string //工单号
	BomId         string //BOM号
	Process       string //工序
	WorkOrderType string //工单类型
	StartTime     string
	EndTime       string
}

type PnSetParam struct {
	PartNumber       string
	Version          string
	DtFlag           string
	TemperFlag       string
	TxPowerMin       sql.NullFloat64
	TxPowerMax       sql.NullFloat64
	TrackingErrorMax sql.NullFloat64
	TxErMin          sql.NullFloat64
	TxErMax          sql.NullFloat64
	TxBiasMin        sql.NullFloat64
	TxBiasMax        sql.NullFloat64
	TxCrossMin       sql.NullFloat64
	TxCrossMax       sql.NullFloat64
	TxMargin         sql.NullFloat64
	Smsr             sql.NullFloat64
	OverLoad         sql.NullFloat64
	Sens             sql.NullFloat64
	SdasSert         sql.NullFloat64
	SddeasSert       sql.NullFloat64
	TxPowPrec        sql.NullFloat64
	RxPowPrec        sql.NullFloat64
	FirmWareVersion  sql.NullString
	FepNumber        sql.NullString
}

func GetPnSetParams(queryCondition *QueryCondition) (pnSetParamList []PnSetParam, err error) {
	sqlStr := `select distinct t.partnumber,t.version,t.dt_flag,t.temper_flag,t.txpowmin,t.txpowmax,t.trackingerrormax,t.txermin,t.txermax,
                t.txbiasmin,t.txbiasmax,t.txcrossmin,t.txcrossmax,t.txmargin,t.smsr,t.OVERLOAD,t.SENS,t.SDASSERT,t.SDDEASSERT,b.tx_pow_prec,b.rx_pow_prec,
                c.firmware_ver,c.fep_number
from  superxon.autodt_spec_ate_new t JOIN SUPERXON.AUTODT_SPEC_MONITOR b ON t.version=b.version AND T.DT_FLAG=B.DT_FLAG AND T.TEMPER_FLAG=B.TEMPER_FLAG
                                     join superxon.autodt_spec_image c on t.partnumber=c.partnumber and t.version=c.version
                                     join superxon.sgd_scdd_trx d on t.partnumber=d.partnumber and t.version=d.version
where t.partnumber LIKE '%` + queryCondition.Pn + `%'
  and d.pch_tc like '%` + queryCondition.WorkOrderId + `%'
  and t.version like '%` + queryCondition.BomId + `%'
  and t.dt_flag like '%` + queryCondition.Process + `%'
order by t.version,t.dt_flag,t.temper_flag ASC`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pnSetParam PnSetParam
	for rows.Next() {
		err = rows.Scan(
			&pnSetParam.PartNumber,
			&pnSetParam.Version,
			&pnSetParam.DtFlag,
			&pnSetParam.TemperFlag,
			&pnSetParam.TxPowerMin,
			&pnSetParam.TxPowerMax,
			&pnSetParam.TrackingErrorMax,
			&pnSetParam.TxErMin,
			&pnSetParam.TxErMax,
			&pnSetParam.TxBiasMin,
			&pnSetParam.TxBiasMax,
			&pnSetParam.TxCrossMin,
			&pnSetParam.TxCrossMax,
			&pnSetParam.TxMargin,
			&pnSetParam.Smsr,
			&pnSetParam.OverLoad,
			&pnSetParam.Sens,
			&pnSetParam.SdasSert,
			&pnSetParam.SddeasSert,
			&pnSetParam.TxPowPrec,
			&pnSetParam.RxPowPrec,
			&pnSetParam.FirmWareVersion,
			&pnSetParam.FepNumber)
		if err != nil {
			return nil, err
		}
		pnSetParamList = append(pnSetParamList, pnSetParam)
	}
	return
}

type PnWorkOrderYield struct {
	WorkOrderId string
	Pn          string
	Seq         string
	Process     string
	Version     sql.NullString
	TotalInput  uint32
	TotalPass   uint32
	TotalFail   uint32
	YieldRate   string
}

func GetPnWorkOrderYields(queryCondition *QueryCondition) (pnWorkOrderYieldList []PnWorkOrderYield, err error) {
	sqlStr := `with TRX as (SELECT distinct a.MANUFACTURE_GROUP,d.LOT_TYPE,
(case when substr(b.softversion,length(b.softversion)-4) like '%验证软件' then substr(b.softversion,0,length(b.softversion)-5)
when substr(b.softversion,length(b.softversion)-1) LIKE '%*_'escape '*' then substr(b.softversion,0,length(b.softversion)-1) else B.SOFTVERSION END) as SVERSION,b.*,
rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz,
rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr,
c."sequence" as SEQ
FROM superxon.autodt_process_log b,(SELECT C.id,c.partnumber,c.manufacture_group,c.tosa_group,c.rosa_group,c.bosa_group,c.pcba1_group,c.bosa_sn,c.modifydate,c.la
FROM (select t.*,dense_rank()over(partition by T.PARTNUMBER,T.BOSA_SN order by T.MODIFYDATE DESC)LA from superxon.autodt_tracking t)C where C.LA=1) a,superxon.workstage c,
(select t.partnumber,t.version,t.pch_tc, (case when  substr(t.pch_lx,0,10) like'TRX试生产产品工单%' then 'TRX正常品'
when substr(t.pch_lx,0,10) like  'TRX量产产品工单%' then 'TRX正常品'  else'TRX改制返工品' END) as LOT_TYPE,t.pch_lx
from superxon.sgd_scdd_trx t) d
where b.sn=a.bosa_sn and b.log_action = c."processname" and a.partnumber =b.pn and d.pch_tc=a.manufacture_group and b.pn=d.partnumber
and ( b.pn LIKE  '%` + queryCondition.Pn + `%' /*and b.log_action like '&工序%'*/ AND D.LOT_TYPE LIKE '` + queryCondition.WorkOrderType + `%')
and b.action_time between to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss')
and to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')
and b.log_action not like 'MODULE_SN')
select e.* from (select distinct h.MANUFACTURE_GROUP as 工单号,h.PN as PN,h.SEQ as 序列,h.log_action as 工序,h.SVERSION,
count(sn)over(partition by h.log_action,h.PN,h.SVERSION,h.MANUFACTURE_GROUP)总输入,
sum(case h.p_value when 'PASS' then 1 else 0 end)over(partition by h.log_action,h.PN,h.SVERSION,h.MANUFACTURE_GROUP)最终良品,
sum(case h.p_value when 'PASS' then 0 else 1 end)over(partition by h.log_action,h.PN,h.SVERSION,h.MANUFACTURE_GROUP)最终不良品,
ROUND((sum(case h.p_value when 'PASS' then 1 else 0 end)over(partition by h.log_action,h.PN,h.SVERSION,h.MANUFACTURE_GROUP)/(count(sn)over(partition by h.log_action,h.PN,h.SVERSION,h.MANUFACTURE_GROUP))*100),2)||'%' 工位良率
from TRX h where h.rr=1)e WHERE E.总输入>0 order by e.工单号,e.序列 ASC`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pnWorkOrderYield PnWorkOrderYield
	for rows.Next() {
		err = rows.Scan(
			&pnWorkOrderYield.WorkOrderId,
			&pnWorkOrderYield.Pn,
			&pnWorkOrderYield.Seq,
			&pnWorkOrderYield.Process,
			&pnWorkOrderYield.Version,
			&pnWorkOrderYield.TotalInput,
			&pnWorkOrderYield.TotalPass,
			&pnWorkOrderYield.TotalFail,
			&pnWorkOrderYield.YieldRate)
		if err != nil {
			return nil, err
		}
		pnWorkOrderYieldList = append(pnWorkOrderYieldList, pnWorkOrderYield)
	}
	return
}
