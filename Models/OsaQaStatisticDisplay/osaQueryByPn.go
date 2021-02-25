package OsaQaStatisticDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Models/OsaRunDisplay"
	"SuperxonWebSite/Utils"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

func GetQaOsaPnList(osaQueryCondition *OsaRunDisplay.OsaQueryCondition) (osaPnList []string, err error) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", osaQueryCondition.StartTime, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", osaQueryCondition.EndTime, time.Local)
	var sqlStr string

	if endTime.After(startTime.AddDate(0, 0, 6)) {
		osaQueryCondition.StartTime, osaQueryCondition.EndTime = Utils.GetAgoAndCurrentTime(Utils.Ago{Years: 0, Months: -1, Days: 0})
		sqlStr = `select distinct t.partnumber from superxon.sgd_scdd_trx t where t.partnumber LIKE '0%' and t.pch_tc_date between to_date('` + osaQueryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + osaQueryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')`
	} else {
		sqlStr = `select distinct t.pn from superxon.autodt_process_log t where t.pn LIKE '0%' and t.action_time between to_date('` + osaQueryCondition.StartTime + `', 'yyyy-mm-dd hh24:mi:ss') and to_date('` + osaQueryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')`
	}
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var osaPn string
	for rows.Next() {
		err = rows.Scan(
			&osaPn)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		osaPnList = append(osaPnList, osaPn)
	}
	return
}

func RedisGetQaOsaStatisticInfoListByPn(osaQueryCondition *OsaRunDisplay.OsaQueryCondition) (qaOsaStatisticInfoList []OsaRunDisplay.OsaInfo, err error) {
	key := "OsaInfo" + osaQueryCondition.Pn + osaQueryCondition.StartTime + osaQueryCondition.EndTime
	reBytes, _ := redis.Bytes(Databases.RedisPool.Get().Do("get", key))
	if len(reBytes) != 0 {
		_ = json.Unmarshal(reBytes, &qaOsaStatisticInfoList)
		if len(qaOsaStatisticInfoList) != 0 {
			fmt.Println("使用redis")
			return
		}
	}

	qaOsaStatisticInfoList, _ = OsaRunDisplay.GetOsaInfoList(osaQueryCondition)

	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", osaQueryCondition.StartTime, time.Local)
	startTime = startTime.AddDate(0, 0, 7)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", osaQueryCondition.EndTime, time.Local)
	if !startTime.After(endTime) {
		datas, _ := json.Marshal(qaOsaStatisticInfoList)
		_, err = Databases.RedisPool.Get().Do("SET", key, datas, "NX", "EX", 60*60*23+60*50)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
	return
}

type QaOsaDefectsInfo struct {
	Pn             string
	Process        string
	ErrorCode      string
	ErrorCount     uint
	ErrorInputRate string
}

func GetQaOsaDefectsInfoListByPn(osaQueryCondition *OsaRunDisplay.OsaQueryCondition) (qaOsaDefectsInfoList []QaOsaDefectsInfo, err error) {
	sqlStr := `with OSA AS (select y.errorcode, x.* from (SELECT distinct b.LOTNO,d.LOT_TYPE,
(case when substr(a.softversion,length(a.softversion)-4) like '%验证软件' then substr(a.softversion,0,length(a.softversion)-5)
when substr(a.softversion,length(a.softversion)-1) LIKE '%*_'escape '*' then substr(a.softversion,0,length(a.softversion)-1) else a.SOFTVERSION END) as SVERSION,a.*,
dense_rank()over(partition by a.sn,a.log_action order by a.action_time asc)FP,
dense_rank()over(partition by a.sn,a.log_action order by a.action_time DESC)LP,
c."sequence" as SEQ
FROM superxon.autodt_process_log a,(SELECT e.* FROM (select t.*,dense_rank()over(partition by T.PN,T.SN order by T.TXDATETIME DESC)LTX,
dense_rank()over(partition by T.PN,T.SN order by T.RXDATETIME DESC)LRX from superxon.AUTODT_RESULTS_OSA_TRACKING t)e where e.LTX=1) b,superxon.workstage c,
(select t.partnumber,t.version,t.pch_tc, (case when  substr(t.pch_lx,0,10) like'OSA试生产产品工单' then 'OSA正常品' 
when substr(t.pch_lx,0,10) like  'OSA量产产品工单' then 'OSA正常品'  else'OSA改制返工品' END) as LOT_TYPE,t.pch_lx
from superxon.sgd_scdd_trx t)d
where b.sn=a.sn and a.log_action = c."processname" and a.pn =b.pn and d.pch_tc=b.LOTNO and a.pn=d.partnumber
and ( a.pn LIKE  '` + osaQueryCondition.Pn + `' /*and a.log_action like '&工序%' AND D.LOT_TYPE LIKE '&工单类型%'*/) 
and a.action_time between to_date('` + osaQueryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss') 
and to_date('` + osaQueryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss'))x ,
(select * from (select * from (select (t.id+1) as id,t.pn,t.sn,t.errorcode,t.tc_flag,to_date(t.apd_d_date,'yyyy-mm-dd hh24:mi:ss') as TESTDATE from superxon.autodt_results_opticsdata t)
 where TESTDATE between to_date('` + osaQueryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + osaQueryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')
union all
select * from (select t.id as id,t.pn,t.sn,t.errorcode,t.tc_flag,t.testdate as TESTDATE from superxon.autodt_results_liv t
 where TESTDATE between to_date('` + osaQueryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + osaQueryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss'))))y
where x.resultsid = y.id AND x.sn = y.SN and y.errorcode <> '0')

select d.* from (select distinct g.PN as PN,g.log_action as 工序,/*G.SVERSION,*/g.ERRORCODE,
count(G.sn)over(partition by g.ERRORCODE,g.log_action)不良数量,
ROUND((count(G.sn)over(partition by g.ERRORCODE,g.log_action)/(sum(case g.p_value when 'FAIL' then 1 else 0 end)over(partition by g.pn))*100),2)||'%' 不良比重
from OSA g where g.lp=1)d
order by d.不良数量 DESC`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaOsaDefectsInfo QaOsaDefectsInfo
	for rows.Next() {
		err = rows.Scan(
			&qaOsaDefectsInfo.Pn,
			&qaOsaDefectsInfo.Process,
			&qaOsaDefectsInfo.ErrorCode,
			&qaOsaDefectsInfo.ErrorCount,
			&qaOsaDefectsInfo.ErrorInputRate)
		if err != nil {
			return nil, err
		}
		qaOsaDefectsInfoList = append(qaOsaDefectsInfoList, qaOsaDefectsInfo)
	}
	return
}

func RedisGetQaOsaDefectsInfoListByPn(osaQueryCondition *OsaRunDisplay.OsaQueryCondition) (qaOsaDefectsInfoList []QaOsaDefectsInfo, err error) {
	key := "OsaDefectsInfo" + osaQueryCondition.Pn + osaQueryCondition.StartTime + osaQueryCondition.EndTime
	reBytes, _ := redis.Bytes(Databases.RedisPool.Get().Do("get", key))
	if len(reBytes) != 0 {
		_ = json.Unmarshal(reBytes, &qaOsaDefectsInfoList)
		if len(qaOsaDefectsInfoList) != 0 {
			fmt.Println("使用redis")
			return
		}
	}

	qaOsaDefectsInfoList, _ = GetQaOsaDefectsInfoListByPn(osaQueryCondition)

	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", osaQueryCondition.StartTime, time.Local)
	startTime = startTime.AddDate(0, 0, 7)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", osaQueryCondition.EndTime, time.Local)
	if !startTime.After(endTime) {
		datas, _ := json.Marshal(qaOsaDefectsInfoList)
		_, err = Databases.RedisPool.Get().Do("SET", key, datas, "NX", "EX", 60*60*23+60*50)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}
