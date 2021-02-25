package ModuleQaStatisticDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

type QaCpkInfo struct {
	TxAop    float64
	TxER     float64
	A2Ibias  float64
	EaAbsorb float64
	Sigma    float64
	Smsr     float64
}

type QaCpkInfoResult struct {
	TxAop    []float64
	TxER     []float64
	A2Ibias  []float64
	EaAbsorb []float64
	Sigma    []float64
	Smsr     []float64
}

type Process struct {
	Name string
}

func GetAllProcessOfTRX() (processList []Process, err error) {
	sqlStr := `select distinct t."processname" from superxon.workstage t where  t."workshop" like '%TRX%'`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var process Process
	for rows.Next() {
		err = rows.Scan(
			&process.Name,
		)
		if err != nil {
			return nil, err
		}
		processList = append(processList, process)
	}
	return
}

func GetQaCpkInfoList(queryCondition *QueryCondition) (result map[string]map[string]uint, err error) {
	var qaCpkInfoList []QaCpkInfo
	sqlStr := `select b.txaop,B.TXER,B.A2_IBIAS,B.EA_ABSORB,b.sigma,B.SMSR
from (SeLECT distinct x.*,RANK()OVER(partition by x.sn,x.log_action order by x.action_time desc)rr
from superxon.autodt_process_log x
WHERE x.pn like '` + queryCondition.Pn + `'
and x.log_action like '` + queryCondition.Process + `'
and ACTION_TIME >=to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss')
and ACTION_TIME <=to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')
) a  JOIN (SeLECT distinct y.*,RANK()OVER(partition by y.bosa_sn order by y.ID desc)ee
from superxon.autodt_tracking y)c ON a.sn=c.bosa_sn and c.ee=1
 join superxon.autodt_results_ate_new b on a.resultsid=b.id
 where a.rr=1 AND C.EE=1`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaCpkInfo QaCpkInfo
	for rows.Next() {
		err = rows.Scan(
			&qaCpkInfo.TxAop,
			&qaCpkInfo.TxER,
			&qaCpkInfo.A2Ibias,
			&qaCpkInfo.EaAbsorb,
			&qaCpkInfo.Sigma,
			&qaCpkInfo.Smsr,
		)
		if err != nil {
			return nil, err
		}
		qaCpkInfoList = append(qaCpkInfoList, qaCpkInfo)
	}
	result, err = GetQaCpkResult(qaCpkInfoList...)
	return
}

func RedisGetQaCpkInfoList(queryCondition *QueryCondition) (result map[string]map[string]uint, err error) {
	//result = make(map[string]map[string]uint)
	key := "CpkBase" + queryCondition.Pn + queryCondition.Process + queryCondition.StartTime + queryCondition.EndTime
	reBytes, _ := redis.Bytes(Databases.RedisPool.Get().Do("get", key))
	if len(reBytes) != 0 {
		_ = json.Unmarshal(reBytes, &result)
		if len(result) != 0 {
			fmt.Println("使用redis")
			return
		}
	}

	result, _ = GetQaCpkInfoList(queryCondition)

	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", queryCondition.StartTime, time.Local)
	startTime = startTime.AddDate(0, 0, 7)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", queryCondition.EndTime, time.Local)
	if !startTime.After(endTime) {
		datas, _ := json.Marshal(result)
		_, err = Databases.RedisPool.Get().Do("SET", key, datas, "NX", "EX", 60*60*23+60*50)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}

func CronGetQaCpkInfoList(queryCondition *QueryCondition) (result map[string]map[string]uint, err error) {
	result = make(map[string]map[string]uint)
	key := "CpkBase" + queryCondition.Pn + queryCondition.Process + queryCondition.StartTime + queryCondition.EndTime
	fmt.Println(key + "存入redis")
	result, _ = GetQaCpkInfoList(queryCondition)
	datas, _ := json.Marshal(result)
	_, _ = Databases.RedisPool.Get().Do("SET", key, datas)
	_, err = Databases.RedisPool.Get().Do("expire", key, 60*60*30)
	return
}

type QaCpkRssi struct {
	CP1 float64
	CP2 float64
	CP3 float64
	CP4 float64
	CP5 float64
	CP6 float64
	CP7 float64
	CP8 float64
}

type QaCpkRssiResult struct {
	CP1 []float64
	CP2 []float64
	CP3 []float64
	CP4 []float64
	CP5 []float64
	CP6 []float64
	CP7 []float64
	CP8 []float64
}

func GetQaCpkRssiList(queryCondition *QueryCondition) (result map[string]map[string]uint, err error) {
	var qaCpkRssiList []QaCpkRssi
	sqlStr := `select d.calipoint1,d.calipoint2,d.calipoint3,
d.calipoint4,d.calipoint5,d.calipoint6,d.calipoint7,d.calipoint8
from (SeLECT distinct x.*,RANK()OVER(partition by x.sn,x.log_action order by x.action_time desc)rr
from superxon.autodt_process_log x
WHERE x.pn like '` + queryCondition.Pn + `'
and x.log_action like '` + queryCondition.Process + `%'
and ACTION_TIME >=to_date('` + queryCondition.StartTime + `','yyyy-mm-dd hh24:mi:ss')
and ACTION_TIME <=to_date('` + queryCondition.EndTime + `','yyyy-mm-dd hh24:mi:ss')
) a  JOIN (SeLECT distinct y.*,RANK()OVER(partition by y.bosa_sn order by y.ID desc)ee
from superxon.autodt_tracking y)c ON a.sn=c.bosa_sn and c.ee=1
 join superxon.autodt_results_monitor d on a.sn=d.opticssn and a.ACTION_TIME=d.testdate
 where a.rr=1 AND C.EE=1 `
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qaCpkRssi QaCpkRssi
	for rows.Next() {
		err = rows.Scan(
			&qaCpkRssi.CP1,
			&qaCpkRssi.CP2,
			&qaCpkRssi.CP3,
			&qaCpkRssi.CP4,
			&qaCpkRssi.CP5,
			&qaCpkRssi.CP6,
			&qaCpkRssi.CP7,
			&qaCpkRssi.CP8,
		)
		if err != nil {
			return nil, err
		}
		qaCpkRssiList = append(qaCpkRssiList, qaCpkRssi)
	}
	result, err = GetQaCpkRssiResult(qaCpkRssiList...)
	return
}

func RedisGetQaCpkRssiList(queryCondition *QueryCondition) (result map[string]map[string]uint, err error) {
	result = make(map[string]map[string]uint)
	key := "CpkRssi" + queryCondition.Pn + queryCondition.Process + queryCondition.StartTime + queryCondition.EndTime
	reBytes, _ := redis.Bytes(Databases.RedisPool.Get().Do("get", key))
	if len(reBytes) != 0 {
		_ = json.Unmarshal(reBytes, &result)
		if len(result) != 0 {
			fmt.Println("使用redis")
			return
		}
	}

	result, _ = GetQaCpkRssiList(queryCondition)

	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", queryCondition.StartTime, time.Local)
	startTime = startTime.AddDate(0, 0, 7)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", queryCondition.EndTime, time.Local)
	if !startTime.After(endTime) {
		datas, _ := json.Marshal(result)
		_, err = Databases.RedisPool.Get().Do("SET", key, datas, "NX", "EX", 60*60*23+60*50)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	return
}

func CronGetQaCpkRssiList(queryCondition *QueryCondition) (result map[string]map[string]uint, err error) {
	result = make(map[string]map[string]uint)
	key := "CpkRssi" + queryCondition.Pn + queryCondition.Process + queryCondition.StartTime + queryCondition.EndTime
	fmt.Println(key + "存入redis")
	result, _ = GetQaCpkRssiList(queryCondition)
	datas, _ := json.Marshal(result)
	_, _ = Databases.RedisPool.Get().Do("SET", key, datas)
	_, err = Databases.RedisPool.Get().Do("expire", key, 60*60*30)
	return
}

//获得CpkInfo，QaCpkRssi的分布函数
func GetQaCpkResult(qaCpkInfoList ...QaCpkInfo) (result map[string]map[string]uint, err error) {
	result = make(map[string]map[string]uint)
	result["TxAop"] = make(map[string]uint)
	result["TxER"] = make(map[string]uint)
	result["TxAop"] = make(map[string]uint)
	result["A2Ibias"] = make(map[string]uint)
	result["EaAbsorb"] = make(map[string]uint)
	result["Sigma"] = make(map[string]uint)
	result["Smsr"] = make(map[string]uint)

	length := len(qaCpkInfoList)
	var qaCpkInfoResult QaCpkInfoResult
	qaCpkInfoResult.TxAop = make([]float64, length)
	qaCpkInfoResult.TxER = make([]float64, length)
	qaCpkInfoResult.A2Ibias = make([]float64, length)
	qaCpkInfoResult.EaAbsorb = make([]float64, length)
	qaCpkInfoResult.Sigma = make([]float64, length)
	qaCpkInfoResult.Smsr = make([]float64, length)
	for index, qaCpkInfo := range qaCpkInfoList {
		if qaCpkInfo.TxAop > 0 {
			qaCpkInfoResult.TxAop[index] = qaCpkInfo.TxAop
		}
		if qaCpkInfo.TxER > 0 {
			qaCpkInfoResult.TxER[index] = qaCpkInfo.TxER
		}
		if qaCpkInfo.A2Ibias > 0 {
			qaCpkInfoResult.A2Ibias[index] = qaCpkInfo.A2Ibias
		}
		if qaCpkInfo.EaAbsorb > 0 {
			qaCpkInfoResult.EaAbsorb[index] = qaCpkInfo.EaAbsorb
		}
		if qaCpkInfo.Sigma > 0 {
			qaCpkInfoResult.Sigma[index] = qaCpkInfo.Sigma
		}
		if qaCpkInfo.Smsr > 0 {
			qaCpkInfoResult.Smsr[index] = qaCpkInfo.Smsr
		}
	}
	//startT := time.Now()
	c := make(chan bool, 6)
	defer close(c)
	go CpkDataHandle(Utils.RemoveZero(qaCpkInfoResult.TxAop), 0.5, result["TxAop"], c)
	go CpkDataHandle(Utils.RemoveZero(qaCpkInfoResult.TxER), 0.5, result["TxER"], c)
	go CpkDataHandle(Utils.RemoveZero(qaCpkInfoResult.A2Ibias), 5.0, result["A2Ibias"], c)
	go CpkDataHandle(Utils.RemoveZero(qaCpkInfoResult.EaAbsorb), 0.5, result["EaAbsorb"], c)
	go CpkDataHandle(Utils.RemoveZero(qaCpkInfoResult.Sigma), 2.0, result["Sigma"], c)
	go CpkDataHandle(Utils.RemoveZero(qaCpkInfoResult.Smsr), 2.0, result["Smsr"], c)
	for i := 0; i < 6; i++ {
		<-c
	}
	//fmt.Println(time.Since(startT))
	return
}

func CpkDataHandle(slice []float64, segmentInterval float64, dst map[string]uint, c chan bool) {
	sliceMax, sliceMin := Utils.MaxAndMin(segmentInterval, slice...)
	AxisSlice := make([]float64, 0)
	segment := (sliceMax - sliceMin) / segmentInterval //分成多少段
	for i := 0; i < int(segment); i++ {
		AxisSlice = append(AxisSlice, sliceMin+float64(i)*segmentInterval)
	}
	for _, value := range slice {
		for indexAxis, _ := range AxisSlice {
			if indexAxis < (len(AxisSlice) - 1) {
				if value > AxisSlice[indexAxis] && value < AxisSlice[indexAxis+1] {
					dst[strconv.FormatFloat(AxisSlice[indexAxis], 'f', 1, 64)+"_"+strconv.FormatFloat(AxisSlice[indexAxis+1], 'f', 1, 64)] += 1
					break
				}
			}
		}
	}
	c <- true
}

func GetQaCpkRssiResult(qaCpkRssiList ...QaCpkRssi) (result map[string]map[string]uint, err error) {
	result = make(map[string]map[string]uint)
	result["CP1"] = make(map[string]uint)
	result["CP2"] = make(map[string]uint)
	result["CP3"] = make(map[string]uint)
	result["CP4"] = make(map[string]uint)
	result["CP5"] = make(map[string]uint)
	result["CP6"] = make(map[string]uint)
	result["CP7"] = make(map[string]uint)
	result["CP8"] = make(map[string]uint)

	length := len(qaCpkRssiList)
	var QaCpkRssiResult QaCpkRssiResult
	QaCpkRssiResult.CP1 = make([]float64, length)
	QaCpkRssiResult.CP2 = make([]float64, length)
	QaCpkRssiResult.CP3 = make([]float64, length)
	QaCpkRssiResult.CP4 = make([]float64, length)
	QaCpkRssiResult.CP5 = make([]float64, length)
	QaCpkRssiResult.CP6 = make([]float64, length)
	QaCpkRssiResult.CP7 = make([]float64, length)
	QaCpkRssiResult.CP8 = make([]float64, length)

	for index, qaCpkRssi := range qaCpkRssiList {
		if qaCpkRssi.CP1 < 0 {
			QaCpkRssiResult.CP1[index] = qaCpkRssi.CP1
		}
		if qaCpkRssi.CP2 < 0 {
			QaCpkRssiResult.CP2[index] = qaCpkRssi.CP2
		}
		if qaCpkRssi.CP3 < 0 {
			QaCpkRssiResult.CP3[index] = qaCpkRssi.CP3
		}
		if qaCpkRssi.CP4 < 0 {
			QaCpkRssiResult.CP4[index] = qaCpkRssi.CP4
		}
		if qaCpkRssi.CP5 < 0 {
			QaCpkRssiResult.CP5[index] = qaCpkRssi.CP5
		}
		if qaCpkRssi.CP6 < 0 {
			QaCpkRssiResult.CP6[index] = qaCpkRssi.CP6
		}
		if qaCpkRssi.CP7 < 0 {
			QaCpkRssiResult.CP7[index] = qaCpkRssi.CP7
		}
		if qaCpkRssi.CP8 < 0 {
			QaCpkRssiResult.CP8[index] = qaCpkRssi.CP8
		}

	}
	c := make(chan bool, 8)
	defer close(c)
	go CpkRssiDataHandle(Utils.RemoveZero(QaCpkRssiResult.CP1), 0.5, result["CP1"], c)
	go CpkRssiDataHandle(Utils.RemoveZero(QaCpkRssiResult.CP2), 0.5, result["CP2"], c)
	go CpkRssiDataHandle(Utils.RemoveZero(QaCpkRssiResult.CP3), 0.5, result["CP3"], c)
	go CpkRssiDataHandle(Utils.RemoveZero(QaCpkRssiResult.CP4), 0.5, result["CP4"], c)
	go CpkRssiDataHandle(Utils.RemoveZero(QaCpkRssiResult.CP5), 0.5, result["CP5"], c)
	go CpkRssiDataHandle(Utils.RemoveZero(QaCpkRssiResult.CP6), 0.5, result["CP6"], c)
	go CpkRssiDataHandle(Utils.RemoveZero(QaCpkRssiResult.CP7), 0.5, result["CP7"], c)
	go CpkRssiDataHandle(Utils.RemoveZero(QaCpkRssiResult.CP8), 0.5, result["CP8"], c)
	for i := 0; i < 8; i++ {
		<-c
	}
	return
}

func CpkRssiDataHandle(slice []float64, segmentInterval float64, dst map[string]uint, c chan bool) {
	sliceMax, sliceMin := Utils.NegativeMaxAndMin(segmentInterval, slice...)
	AxisSlice := make([]float64, 0)
	segment := (sliceMax - sliceMin) / segmentInterval //分成多少段
	for i := 0; i < int(segment); i++ {
		AxisSlice = append(AxisSlice, sliceMax-float64(i)*segmentInterval)
	}
	for _, value := range slice {
		for indexAxis, _ := range AxisSlice {
			if indexAxis < (len(AxisSlice) - 1) {
				if value < AxisSlice[indexAxis] && value > AxisSlice[indexAxis+1] {
					dst[strconv.FormatFloat(AxisSlice[indexAxis], 'f', 1, 64)+"_"+strconv.FormatFloat(AxisSlice[indexAxis+1], 'f', 1, 64)] += 1
					break
				}
			}
		}
	}
	c <- true
}

/*  result map[string]map[string]uint 示例
{
    "A2Ibias": {
        "100.0": 1,
        "50.0": 4,
        "60.0": 140,
        "70.0": 166,
        "80.0": 69,
        "90.0": 7
    },
    "EaAbsorb": {},
    "Sigma": {
        "0.0": 372
    },
    "Smsr": {},
    "TxAop": {
        "2.5": 17,
        "3.0": 168,
        "3.5": 176,
        "4.0": 25,
        "4.5": 1
    },
    "TxER": {
        "7.0": 9,
        "7.5": 144,
        "8.0": 211,
        "8.5": 21
    }
}
*/
