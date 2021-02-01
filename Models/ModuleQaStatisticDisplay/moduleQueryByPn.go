package ModuleQaStatisticDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type QaPn struct {
	Pn string
}

func GetQaPnList(queryCondition *QueryCondition) (qaPnList []QaPn, err error) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", queryCondition.StartTime, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", queryCondition.EndTime, time.Local)
	var sqlStr string

	if endTime.After(startTime.AddDate(0, 0, 6)) {
		queryCondition.StartTime, queryCondition.EndTime = Utils.GetAgoAndCurrentTime(Utils.Ago{Years: 0, Months: -1, Days: 0})
		sqlStr = `select distinct t.partnumber from superxon.sgd_scdd_trx t where t.partnumber LIKE 'SO%' and t.pch_tc_date between to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')`
	} else {
		sqlStr = `select distinct t.pn from superxon.autodt_process_log t where t.pn LIKE 'SO%' and t.action_time between to_date('` + queryCondition.StartTime + `', 'yyyy-mm-dd hh24:mi:ss') and to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')`
	}

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
	TotalInput    uint32
	FinalOk       uint32
	FinalBad      uint32
	FinalPassRate string
}

func GetQaStatisticOrderInfoList(queryCondition *QueryCondition) (qaStatisticInfoList []QaStatisticInfo, err error) {
	sqlStr := `with TRX as (SELECT distinct a.MANUFACTURE_GROUP,d.LOT_TYPE,b.*,
        rank()over(partition by b.sn,b.log_action order by b.action_time asc)zz,
        rank()over(partition by b.sn,b.log_action order by b.action_time DESC)rr,
        c."sequence" as SEQ
        FROM superxon.autodt_process_log b,(SELECT aa.* FROM (select t.*,dense_rank()over(partition by T.PARTNUMBER,T.BOSA_SN order by T.MODIFYDATE DESC)LA
        from superxon.autodt_tracking t)aa where aa.LA=1) a,superxon.workstage c,
        (select t.partnumber,t.version,t.pch_tc, (case when  substr(t.pch_lx,0,10) like'TRX试生产产品工单%' then 'TRX正常品'
        when substr(t.pch_lx,0,10) like  'TRX量产产品工单%' then 'TRX正常品'  else'TRX改制返工品' END) as LOT_TYPE,t.pch_lx
        from superxon.sgd_scdd_trx t) d
        where b.sn=a.bosa_sn and b.log_action = c."processname" and a.partnumber =b.pn and D.LOT_TYPE like '%` + queryCondition.WorkOrderType + `%'
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

func RedisGetQaStatisticOrderInfoList(queryCondition *QueryCondition) (qaStatisticInfoList []QaStatisticInfo, err error) {
	key := "ModuleInfo" + queryCondition.Pn + queryCondition.WorkOrderType + queryCondition.StartTime + queryCondition.EndTime
	reBytes, _ := redis.Bytes(Databases.RedisPool.Get().Do("get", key))
	if len(reBytes) != 0 {
		_ = json.Unmarshal(reBytes, &qaStatisticInfoList)
		if len(qaStatisticInfoList) != 0 {
			fmt.Println("使用redis")
			return
		}
	}

	qaStatisticInfoList, _ = GetQaStatisticOrderInfoList(queryCondition)

	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", queryCondition.StartTime, time.Local)
	startTime = startTime.AddDate(0, 0, 7)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", queryCondition.EndTime, time.Local)
	if !startTime.After(endTime) {
		datas, _ := json.Marshal(qaStatisticInfoList)
		_, err = Databases.RedisPool.Get().Do("SET", key, datas)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = Databases.RedisPool.Get().Do("expire", key, 60*60*30)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}

type QaDefectsInfoByPn struct {
	Pn             string
	Sequence       string
	Version        string
	ErrorCode      string
	ErrorCount     uint32
	ErrorRate      string
	ErrorInputRate string
}

func GetQaDefectsOrderInfoListByPn(queryCondition *QueryCondition) (qaDefectsInfoList []QaDefectsInfoByPn, err error) {
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
				and b.pn = '` + queryCondition.Pn + `' and D.LOT_TYPE LIKE '%` + queryCondition.WorkOrderType + `%'
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
	var qaDefectsInfo QaDefectsInfoByPn
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

func RedisGetQaDefectsOrderInfoListByPn(queryCondition *QueryCondition) (qaDefectsInfoList []QaDefectsInfoByPn, err error) {
	key := "ModuleDefects" + queryCondition.Pn + queryCondition.WorkOrderType + queryCondition.StartTime + queryCondition.EndTime
	reBytes, _ := redis.Bytes(Databases.RedisPool.Get().Do("get", key))
	if len(reBytes) != 0 {
		_ = json.Unmarshal(reBytes, &qaDefectsInfoList)
		if len(qaDefectsInfoList) != 0 {
			fmt.Println("使用redis")
			return
		}
	}

	qaDefectsInfoList, _ = GetQaDefectsOrderInfoListByPn(queryCondition)

	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", queryCondition.StartTime, time.Local)
	startTime = startTime.AddDate(0, 0, 7)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", queryCondition.EndTime, time.Local)
	if !startTime.After(endTime) {
		datas, _ := json.Marshal(qaDefectsInfoList)
		_, err = Databases.RedisPool.Get().Do("SET", key, datas)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = Databases.RedisPool.Get().Do("expire", key, 60*60*30)
		if err != nil {
			fmt.Println(err)
			return
		}
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

func GetQaDefectsDetailByPn(queryCondition *QueryCondition) (qaDefectsDetailInfoList []QaDefectsDetailInfo, err error) {
	sqlStr := `select c.manufacture_group,d.PCH_LX,a.pn,a.sn,a.operator,a.log_action,a.action_time,a.parameter,a.p_value,a.comments,a.stationid,a.softversion,a.RUNTIME,a.RUNCOUNT,a.COMMITVER,b.MODULESN,b.version,b.temperature,b.errorcode,b.txaop,b.txer,
				b.txppj,b.txtr,b.txtf,b.marging,b.txcross,b.txoffpower,b.overload,b.sensitivity,b.sdassert,b.sddassert,
				b.pathpenalty,b.peakwavelength,b.sigma,b.bandwidth,b.smsr,b.txi,b.rxi,b.te,b.dac_apc,b.dac_mod,
				b.DAC_CROSS,b.DAC_APD,b.DAC_LOS,b.INFO_APC,b.INFO_MOD,b.INFO_CROSS,b.INFO_APD,b.INFO_LOS,b.A2_TEMPERATRUE,b.A2_VCC,b.A2_IBIAS,
				b.A2_TXMON,b.CASETEMPERATRUE,b.TEC_TEMPERATURE,b.EA_ABSORB,b.INFO_EA,b.DAC_EA,b.CENTER_SENSITIVITY,b.AOP_STABILITY01,b.AOP_STABILITY02,
				b.AOP_STABILITY03,b.AOP_STABILITY04,b.AOP_STABILITY05,b.AOP_STABILITY06,b.AOP_STABILITY07,b.AOP_STABILITY08,b.AOP_STABILITY09,b.AOP_STABILITY10,
				b.AOP_STABILITY_DELTA,b.CURRENT5V,b.SET_VCC
				from (SeLECT distinct x.*,RANK()OVER(partition by x.sn,x.log_action order by x.action_time desc)rr
				from superxon.autodt_process_log x
				WHERE x.pn like '` + queryCondition.Pn + `'
				and ACTION_TIME >=to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss')
				and ACTION_TIME <=to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')
				and x.log_action like '` + queryCondition.Process + `'
				) a join superxon.autodt_results_ate_new b on a.resultsid=b.id
				and b.errorcode='` + queryCondition.ErrorCode + `'
				JOIN (SeLECT distinct t.*,RANK()OVER(partition by t.bosa_sn order by t.ID desc)ee
				from superxon.autodt_tracking t)c ON a.sn=c.bosa_sn and c.ee=1,
				(select t.partnumber,t.version,t.pch_tc, (case when  substr(t.pch_lx,0,10) like'TRX试生产产品工单%' then 'TRX正常品'
				when substr(t.pch_lx,0,10) like  'TRX量产产品工单%' then 'TRX正常品'  else'TRX改制返工品' END) as LOT_TYPE,t.pch_lx
				from superxon.sgd_scdd_trx t) d
				where d.pch_tc=c.manufacture_group and c.partnumber=d.partnumber
				and d.LOT_TYPE like '%` + queryCondition.WorkOrderType + `%'
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
where t.partnumber LIKE '` + queryCondition.Pn + `'
  and d.pch_tc like '` + queryCondition.WorkOrderId + `'
  and t.version like '` + queryCondition.BomId + `%'
  and t.dt_flag like '` + queryCondition.Process + `'
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

type WorkOrderYieldByPn struct {
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

/*
	暂时未使用到的良率函数
*/
func GetWorkOrderYieldsByPn(queryCondition *QueryCondition) (pnWorkOrderYieldList []WorkOrderYieldByPn, err error) {
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
and ( b.pn LIKE  '` + queryCondition.Pn + `' /*and b.log_action like '&工序%'*/ AND D.LOT_TYPE LIKE '%` + queryCondition.WorkOrderType + `%')
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
	var pnWorkOrderYield WorkOrderYieldByPn
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
