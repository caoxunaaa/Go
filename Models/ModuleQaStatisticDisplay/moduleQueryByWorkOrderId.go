package ModuleQaStatisticDisplay

import (
	"SuperxonWebSite/Databases"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type QueryCondition struct {
	Pn            string //产品型号
	WorkOrderId   string //工单号
	BomId         string //BOM号
	Process       string //工序
	WorkOrderType string //工单类型
	ErrorCode     string //错误代码
	StartTime     string
	EndTime       string
}

type QaWorkOrderId struct {
	WorkOrderId   string
	Pn            string
	WorkOrderType string
}

//获取PN和工单类型下的时间段的工单号
func GetWorkOrderIds(queryCondition *QueryCondition, isFinish string) (qaWorkOrderIdList []QaWorkOrderId, err error) {
	var sqlStr string
	if isFinish == "no" {
		sqlStr = `SELECT d.pch_tc,d.partnumber,d.LOT_TYPE from (select t.*,
				(case when  substr(t.pch_lx,0,10) like'TRX试生产产品工单%' then 'TRX正常品'
				when substr(t.pch_lx,0,10) like  'TRX量产产品工单%' then 'TRX正常品'  else'TRX改制返工品' END) as LOT_TYPE
				from superxon.sgd_scdd_trx t) d
				WHERE d.partnumber LIKE '` + queryCondition.Pn + `' /*and b.log_action like '&工序%'*/
				AND D.LOT_TYPE LIKE '%` + queryCondition.WorkOrderType + `%'
				AND D.STATE NOT LIKE '结案'`
	} else {
		sqlStr = `SELECT d.pch_tc,d.partnumber,d.LOT_TYPE from (select t.*,
				(case when  substr(t.pch_lx,0,10) like'TRX试生产产品工单%' then 'TRX正常品'
				when substr(t.pch_lx,0,10) like  'TRX量产产品工单%' then 'TRX正常品'  else'TRX改制返工品' END) as LOT_TYPE
				from superxon.sgd_scdd_trx t) d
				WHERE d.partnumber LIKE '` + queryCondition.Pn + `' /*and b.log_action like '&工序%'*/
				AND D.LOT_TYPE LIKE '%` + queryCondition.WorkOrderType + `%'
				AND D.STATE LIKE '结案'
				and D.PCH_TC_DATE>=to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss')
				and D.PCH_TC_DATE<=to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')`
	}

	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaWorkOrderId QaWorkOrderId
	for rows.Next() {
		err = rows.Scan(
			&qaWorkOrderId.WorkOrderId,
			&qaWorkOrderId.Pn,
			&qaWorkOrderId.WorkOrderType)
		if err != nil {
			return nil, err
		}
		qaWorkOrderIdList = append(qaWorkOrderIdList, qaWorkOrderId)
	}
	return
}

type WorkOrderYieldByWorkOrderId struct {
	WorkOrderId string
	Pn          string
	Seq         string
	Process     string
	OrderTime   string
	TotalInput  uint32
	TotalPass   uint32
	TotalFail   uint32
	YieldRate   string
}

func GetWorkOrderYieldsByWorkOrderId(queryCondition *QueryCondition) (workOrderYieldList []WorkOrderYieldByWorkOrderId, err error) {
	sqlStr := `with TRX as (SELECT distinct t.MANUFACTURE_GROUP,b.*,
				rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz,
				rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr, 
				c."sequence" as SEQ,d.PCH_TC_DATE AS 下单时间
				FROM superxon.autodt_process_log b, superxon.autodt_tracking t,superxon.workstage c,superxon.sgd_scdd_trx d
				where b.sn=t.bosa_sn and b.log_action = c."processname" AND t.manufacture_group=D.PCH_TC
				and b.action_time>=d.PCH_TC_DATE
				--and b.action_time<=to_date('&T2','yyyy-mm-dd hh24:mi:ss') 
				and t.manufacture_group = '` + queryCondition.WorkOrderId + `'
				and b.pn like '%SO%'
				and b.log_action not in('TC1_IN','TC1_OUT','MODULE_SN','TRACKING SN'))
				select e.*,round(e.最终良品/e.总投入*100,2)||'%' 最终良率 from
				(select distinct h.manufacture_group as 工单号,h.PN as PN,h.SEQ as 序列,h.log_action as 工序,h.下单时间,
				count(sn)over(partition by h.log_action,h.PN)总投入,
				sum(case h.p_value when 'PASS' then 1 else 0 end)over(partition by h.log_action,h.PN)最终良品,
				sum(case h.p_value when 'PASS' then 0 else 1 end)over(partition by h.log_action,h.PN)最终不良品          
				from TRX h where h.rr=1)e
				order by e.pn,e.序列 ASC`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var workOrderYield WorkOrderYieldByWorkOrderId
	for rows.Next() {
		err = rows.Scan(
			&workOrderYield.WorkOrderId,
			&workOrderYield.Pn,
			&workOrderYield.Seq,
			&workOrderYield.Process,
			&workOrderYield.OrderTime,
			&workOrderYield.TotalInput,
			&workOrderYield.TotalPass,
			&workOrderYield.TotalFail,
			&workOrderYield.YieldRate)
		if err != nil {
			return nil, err
		}
		workOrderYieldList = append(workOrderYieldList, workOrderYield)
	}
	return
}

type QaDefectsInfoByWorkOrderId struct {
	WorkOrderId    string
	Pn             string
	Process        string
	ErrorCode      string
	ErrorCount     uint32
	ErrorInputRate string
}

//------------------------------------------------------------

/*
func GetProjectPlanList() (projectPlanInfoList []ProjectPlanInfo, err error)
从数据库中查询工单号对应的不良代码分布并返回
*/
func GetQaDefectsInfoByWorkOrderIdList(queryCondition *QueryCondition) (qaDefectsInfoByWorkOrderIdList []QaDefectsInfoByWorkOrderId, err error) {
	sqlStr := `with TRX as (select y.errorcode,x.* from (SELECT distinct a.MANUFACTURE_GROUP,b.*,
				rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz,
				rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr, 
				c."sequence" as SEQ,d.PCH_TC_DATE AS 下单时间
				FROM superxon.autodt_process_log b, superxon.autodt_tracking A,superxon.workstage c,superxon.sgd_scdd_trx d
				where b.sn=a.bosa_sn and b.log_action = c."processname" AND A.manufacture_group=D.PCH_TC
				and b.action_time>=d.PCH_TC_DATE
				--and b.action_time<=to_date('&T2','yyyy-mm-dd hh24:mi:ss') 
				and a.manufacture_group = '` + queryCondition.WorkOrderId + `'
				and b.pn like '%SO%'
				and b.log_action not in('TC1_IN','TC1_OUT','MODULE_SN','TRACKING SN')
				)x, superxon.autodt_results_ate_new y
				where x.sn=y.opticssn and x.resultsid =y.id and y.errorcode <> '0')
				
				select d.* from (select distinct G.manufacture_group as 工单号,g.PN as PN,g.log_action as 工序,g.ERRORCODE,
				count(G.sn)over(partition by G.PN,g.ERRORCODE,g.log_action)不良数量,
				ROUND((count(G.sn)over(partition by g.ERRORCODE,g.log_action)/(sum(case g.p_value when 'FAIL' then 1 else 0 end)over(partition by g.PN))*100),2)||'%' 不良比重
				from TRX g where g.RR=1)d`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaDefectsInfoByWorkOrderId QaDefectsInfoByWorkOrderId
	for rows.Next() {
		err = rows.Scan(
			&qaDefectsInfoByWorkOrderId.WorkOrderId,
			&qaDefectsInfoByWorkOrderId.Pn,
			&qaDefectsInfoByWorkOrderId.Process,
			&qaDefectsInfoByWorkOrderId.ErrorCode,
			&qaDefectsInfoByWorkOrderId.ErrorCount,
			&qaDefectsInfoByWorkOrderId.ErrorInputRate)
		if err != nil {
			return nil, err
		}
		qaDefectsInfoByWorkOrderIdList = append(qaDefectsInfoByWorkOrderIdList, qaDefectsInfoByWorkOrderId)
	}
	return
}

/*
func RedisGetQaDefectsInfoByWorkOrderIdList(queryCondition *QueryCondition) (qaDefectsInfoByWorkOrderIdList []QaDefectsInfoByWorkOrderId, err error)
获取redis缓存中的qaDefectsInfoByWorkOrderIdList，如果没有就重新在数据库中查询
*/
func RedisGetQaDefectsInfoByWorkOrderIdList(queryCondition *QueryCondition) (qaDefectsInfoByWorkOrderIdList []QaDefectsInfoByWorkOrderId, err error) {
	reBytes, _ := redis.Bytes(Databases.RedisConn.Do("get", "qaDefectsInfoByWorkOrderIdList"+queryCondition.StartTime+queryCondition.EndTime))
	_ = json.Unmarshal(reBytes, &qaDefectsInfoByWorkOrderIdList)
	fmt.Println(len(qaDefectsInfoByWorkOrderIdList), qaDefectsInfoByWorkOrderIdList)
	if len(qaDefectsInfoByWorkOrderIdList) != 0 {
		fmt.Println("使用redis")
		return
	}
	qaDefectsInfoByWorkOrderIdList, _ = GetQaDefectsInfoByWorkOrderIdList(queryCondition)
	return
}

/*
func CronGetQaDefectsInfoByWorkOrderIdList(queryCondition *QueryCondition) (qaDefectsInfoByWorkOrderIdList []QaDefectsInfoByWorkOrderId, err error)
定时从数据库中查询工单号对应的不良代码分布结果，并存入redis缓存中
*/
func CronGetQaDefectsInfoByWorkOrderIdList(queryCondition *QueryCondition) (qaDefectsInfoByWorkOrderIdList []QaDefectsInfoByWorkOrderId, err error) {
	qaDefectsInfoByWorkOrderIdList, _ = GetQaDefectsInfoByWorkOrderIdList(queryCondition)
	fmt.Println("projectPlanInfoList定时任务使用redis")
	datas, _ := json.Marshal(qaDefectsInfoByWorkOrderIdList)
	_, _ = Databases.RedisConn.Do("SET", "qaDefectsInfoByWorkOrderIdList"+queryCondition.StartTime+queryCondition.EndTime, datas)
	_, err = Databases.RedisConn.Do("expire", "qaDefectsInfoByWorkOrderIdList"+queryCondition.StartTime+queryCondition.EndTime, 60*60*30)
	return
}

//------------------------------------------------------------

/*
func GetQaDefectsDetailByWorkOrderId(queryCondition *QueryCondition) (qaDefectsDetailInfoList []QaDefectsDetailInfo, err error)
从数据库中查询工单号对应的不良代码明细并返回
*/
func GetQaDefectsDetailByWorkOrderId(queryCondition *QueryCondition) (qaDefectsDetailInfoList []QaDefectsDetailInfo, err error) {
	sqlStr := `select c.manufacture_group,t.PCH_LX,a.pn,a.sn,a.operator,a.log_action,a.action_time,a.parameter,a.p_value,a.comments,a.stationid,a.softversion,a.RUNTIME,a.RUNCOUNT,a.COMMITVER,b.MODULESN,b.version,b.temperature,b.errorcode,b.txaop,b.txer,
				b.txppj,b.txtr,b.txtf,b.marging,b.txcross,b.txoffpower,b.overload,b.sensitivity,b.sdassert,b.sddassert,
				b.pathpenalty,b.peakwavelength,b.sigma,b.bandwidth,b.smsr,b.txi,b.rxi,b.te,b.dac_apc,b.dac_mod,
				b.DAC_CROSS,b.DAC_APD,b.DAC_LOS,b.INFO_APC,b.INFO_MOD,b.INFO_CROSS,b.INFO_APD,b.INFO_LOS,b.A2_TEMPERATRUE,b.A2_VCC,b.A2_IBIAS,
				b.A2_TXMON,b.CASETEMPERATRUE,b.TEC_TEMPERATURE,b.EA_ABSORB,b.INFO_EA,b.DAC_EA,b.CENTER_SENSITIVITY,b.AOP_STABILITY01,b.AOP_STABILITY02,
				b.AOP_STABILITY03,b.AOP_STABILITY04,b.AOP_STABILITY05,b.AOP_STABILITY06,b.AOP_STABILITY07,b.AOP_STABILITY08,b.AOP_STABILITY09,b.AOP_STABILITY10,
				b.AOP_STABILITY_DELTA,b.CURRENT5V,b.SET_VCC
				from (SeLECT distinct x.*,RANK()OVER(partition by x.sn,x.log_action order by x.action_time desc)rr
				from superxon.autodt_process_log x
				WHERE x.log_action like '` + queryCondition.Process + `'
				) a join superxon.autodt_results_ate_new b on a.resultsid=b.id
				and b.errorcode='` + queryCondition.ErrorCode + `'
				JOIN (SeLECT distinct t.*,RANK()OVER(partition by t.bosa_sn order by t.ID desc)ee
				from superxon.autodt_tracking t)c ON a.sn=c.bosa_sn and c.ee=1
				join superxon.sgd_scdd_trx T on c.manufacture_group=T.PCH_TC
				where  a.action_time>=t.PCH_TC_DATE
				and c.manufacture_group like '` + queryCondition.WorkOrderId + `'
				and a.rr=1 AND C.EE=1`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaDefectsDetailInfo QaDefectsDetailInfo
	for rows.Next() {
		err = rows.Scan(
			&qaDefectsDetailInfo.WorkOrderId,
			&qaDefectsDetailInfo.PchLx,
			&qaDefectsDetailInfo.Pn,
			&qaDefectsDetailInfo.Sn,
			&qaDefectsDetailInfo.Operator,
			&qaDefectsDetailInfo.LogAction,
			&qaDefectsDetailInfo.ActionTime,
			&qaDefectsDetailInfo.Parameter,
			&qaDefectsDetailInfo.PValue,
			&qaDefectsDetailInfo.Comments,
			&qaDefectsDetailInfo.StationId,
			&qaDefectsDetailInfo.SoftVersion,
			&qaDefectsDetailInfo.RunTime,
			&qaDefectsDetailInfo.RunCount,
			&qaDefectsDetailInfo.CommitVer,
			&qaDefectsDetailInfo.ModuleSn,
			&qaDefectsDetailInfo.Version,
			&qaDefectsDetailInfo.Temperature,
			&qaDefectsDetailInfo.ErrorCode,
			&qaDefectsDetailInfo.TxAOP,
			&qaDefectsDetailInfo.TxER,
			&qaDefectsDetailInfo.TxPPJ,
			&qaDefectsDetailInfo.TxTR,
			&qaDefectsDetailInfo.TxTF,
			&qaDefectsDetailInfo.MARGIN,
			&qaDefectsDetailInfo.TxCROSS,
			&qaDefectsDetailInfo.TxOffPOWER,
			&qaDefectsDetailInfo.OVERLOAD,
			&qaDefectsDetailInfo.SENSITIVITY,
			&qaDefectsDetailInfo.SDASSERT,
			&qaDefectsDetailInfo.SDDASSERT,
			&qaDefectsDetailInfo.PathPenalty,
			&qaDefectsDetailInfo.PeakWavelength,
			&qaDefectsDetailInfo.SIGMA,
			&qaDefectsDetailInfo.BandWidth,
			&qaDefectsDetailInfo.SMSR,
			&qaDefectsDetailInfo.TXI,
			&qaDefectsDetailInfo.RXI,
			&qaDefectsDetailInfo.TE,
			&qaDefectsDetailInfo.DacApc,
			&qaDefectsDetailInfo.DacMod,
			&qaDefectsDetailInfo.DacCross,
			&qaDefectsDetailInfo.DacApd,
			&qaDefectsDetailInfo.DacLos,
			&qaDefectsDetailInfo.InfoApc,
			&qaDefectsDetailInfo.InfoMod,
			&qaDefectsDetailInfo.InfoCross,
			&qaDefectsDetailInfo.InfoApd,
			&qaDefectsDetailInfo.InfoLos,
			&qaDefectsDetailInfo.A2Temperature,
			&qaDefectsDetailInfo.A2Vcc,
			&qaDefectsDetailInfo.A2Ibias,
			&qaDefectsDetailInfo.A2TxMon,
			&qaDefectsDetailInfo.CaseTemperature,
			&qaDefectsDetailInfo.TecTemperature,
			&qaDefectsDetailInfo.EAAbsorb,
			&qaDefectsDetailInfo.InfoEA,
			&qaDefectsDetailInfo.DacEA,
			&qaDefectsDetailInfo.CenterSensitivity,
			&qaDefectsDetailInfo.AopStability01,
			&qaDefectsDetailInfo.AopStability02,
			&qaDefectsDetailInfo.AopStability03,
			&qaDefectsDetailInfo.AopStability04,
			&qaDefectsDetailInfo.AopStability05,
			&qaDefectsDetailInfo.AopStability06,
			&qaDefectsDetailInfo.AopStability07,
			&qaDefectsDetailInfo.AopStability08,
			&qaDefectsDetailInfo.AopStability09,
			&qaDefectsDetailInfo.AopStability10,
			&qaDefectsDetailInfo.AopStabilityDelta,
			&qaDefectsDetailInfo.Current5V,
			&qaDefectsDetailInfo.SetVcc,
		)
		if err != nil {
			return nil, err
		}
		qaDefectsDetailInfoList = append(qaDefectsDetailInfoList, qaDefectsDetailInfo)
	}
	return
}

/*
func RedisGetQaDefectsDetailByWorkOrderId(queryCondition *QueryCondition) (qaDefectsInfoByWorkOrderIdList []QaDefectsInfoByWorkOrderId, err error)
获取redis缓存中的qaDefectsInfoByWorkOrderIdList，如果没有就重新在数据库中查询
*/
func RedisGetQaDefectsDetailByWorkOrderId(queryCondition *QueryCondition) (qaDefectsDetailInfoList []QaDefectsDetailInfo, err error) {
	reBytes, _ := redis.Bytes(Databases.RedisConn.Do("get", "qaDefectsDetailInfoList"+queryCondition.StartTime+queryCondition.EndTime))
	_ = json.Unmarshal(reBytes, &qaDefectsDetailInfoList)
	fmt.Println(len(qaDefectsDetailInfoList), qaDefectsDetailInfoList)
	if len(qaDefectsDetailInfoList) != 0 {
		fmt.Println("使用redis")
		return
	}
	qaDefectsDetailInfoList, _ = GetQaDefectsDetailByWorkOrderId(queryCondition)
	return
}

/*
func CronGetQaDefectsDetailByWorkOrderId(queryCondition *QueryCondition) (qaDefectsInfoByWorkOrderIdList []QaDefectsInfoByWorkOrderId, err error)
定时从数据库中查询工单号对应的不良代码分布结果，并存入redis缓存中
*/
func CronGetQaDefectsDetailByWorkOrderId(queryCondition *QueryCondition) (qaDefectsDetailInfoList []QaDefectsDetailInfo, err error) {
	qaDefectsDetailInfoList, _ = GetQaDefectsDetailByWorkOrderId(queryCondition)
	fmt.Println("projectPlanInfoList定时任务使用redis")
	datas, _ := json.Marshal(qaDefectsDetailInfoList)
	_, _ = Databases.RedisConn.Do("SET", "qaDefectsDetailInfoList"+queryCondition.StartTime+queryCondition.EndTime, datas)
	_, err = Databases.RedisConn.Do("expire", "qaDefectsDetailInfoList"+queryCondition.StartTime+queryCondition.EndTime, 60*60*30)
	return
}
