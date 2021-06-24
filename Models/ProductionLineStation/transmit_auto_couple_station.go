package ProductionLineStation

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Utils"
	"strings"
	"time"
)

type TransmitAutoCoupleStationOverView struct {
	InsName    string    `db:"insname"` //台位Id
	Input      int       `db:"input"`   //总投入
	Pass       int       `db:"pass"`    //总通过数
	PassRate   float64   // 总通过率(根据总投入和通过数计算)
	StartTime  time.Time `db:"starttime"`  //开线时间
	LatestTime time.Time `db:"latesttime"` //最近一只产品生产的时间
	LatestPn   string    `db:"latestpn"`   //最近一只产品pn
	WorkStatus string    //通过LatestTime和当前时间比较，超过半小时就是 闲置中，否则 工作中
	WorkNum    string    `db:"worknum"` //工单号
	Worker     string    `db:"worker"`  //工人工号
}

//总览所有发射耦合台位 总产出output 总通过数pass 开机时间starttime, 最后一只的时间latesttime，最后一只的pn，worknum,worker
func GetAllTransmitAutoCoupleStationOverView(startTime, endTime string) ([]TransmitAutoCoupleStationOverView, error) {
	var res = make([]TransmitAutoCoupleStationOverView, 0)
	now := time.Now()
	sqlStr := `select a.insname, a.input, a.pass, a.starttime, a.latesttime, b.pn as latestpn,b.worknum,b.worker
from (select t.insname, COUNT(t.insname) as input, SUM(case when t.status='Pass' then 1 else 0 end) as pass, MIN(t.testtime) as starttime, MAX(t.testtime) as latesttime
from superxon.autodt_transmit_autocouple t
where t.testtime between to_date('` + startTime + `', 'yyyy-mm-dd hh24:mi:ss') and to_date('` + endTime + `', 'yyyy-mm-dd hh24:mi:ss')
group by t.insname) a,
superxon.autodt_transmit_autocouple b
where a.insname = b.insname
and a.latesttime = b.testtime`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var t TransmitAutoCoupleStationOverView
	for rows.Next() {
		err = rows.Scan(
			&t.InsName,
			&t.Input,
			&t.Pass,
			&t.StartTime,
			&t.LatestTime,
			&t.LatestPn,
			&t.WorkNum,
			&t.Worker,
		)
		if err != nil {
			return nil, err
		}
		t.InsName = strings.TrimSpace(t.InsName)
		t.LatestPn = strings.TrimSpace(t.LatestPn)
		t.WorkNum = strings.TrimSpace(t.WorkNum)
		t.Worker = strings.TrimSpace(t.Worker)
		t.PassRate = Utils.FloatRound(float64(t.Pass)/float64(t.Input)*100, 2)

		if (now.Sub(t.LatestTime).Minutes() + 8*60) > 30 {
			t.WorkStatus = "闲置中"
		} else {
			t.WorkStatus = "工作中"
		}

		res = append(res, t)
	}
	return res, nil
}

type TransmitAutoCoupleStationDetailInfo struct {
	InsName         string    `db:"insname"`   //台位Id
	Pn              string    `db:"pn"`        //产品pn
	TestTime        time.Time `db:"testtime"`  //测试时间
	Status          string    `db:"status"`    //测试结果
	PwrUOk          int       `db:"pwruok"`    //焊接相关
	PwrEnd          int       `db:"pwrend"`    //焊接相关
	PwrZOk          int       `db:"pwrzok"`    //焊接相关
	PwrZWeld4       int       `db:"pwrzweld4"` //焊接相关
	LapWeldRate     float64   //搭接焊变率
	ThroughWeldRate float64   //穿透焊变率
	WaveLength      int       `db:"wavelength"` //波长
}

//通过时间段和某个台位编号获得 所有生产可用信息
func GetTransmitAutoCoupleStationDetailInfoByInsName(insname, startTime, endTime string) ([]TransmitAutoCoupleStationDetailInfo, error) {
	var res = make([]TransmitAutoCoupleStationDetailInfo, 0)
	sqlStr := `select t.insname,t.pn,t.testtime,t.status,t.pwruok,t.pwrend,t.pwrzok,t.pwrzweld4,t.wavelength from superxon.autodt_transmit_autocouple t
where t.insname = '` + insname + `' and t.testtime between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss') order by t.status`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var t TransmitAutoCoupleStationDetailInfo
	for rows.Next() {
		err = rows.Scan(
			&t.InsName,
			&t.Pn,
			&t.TestTime,
			&t.Status,
			&t.PwrUOk,
			&t.PwrEnd,
			&t.PwrZOk,
			&t.PwrZWeld4,
			&t.WaveLength,
		)
		if err != nil {
			return nil, err
		}
		t.InsName = strings.TrimSpace(t.InsName)
		t.Pn = strings.TrimSpace(t.Pn)
		t.Status = strings.TrimSpace(t.Status)

		weldRate1 := Utils.FloatRound(float64(t.PwrUOk-t.PwrEnd)/float64(t.PwrUOk)*100, 2)
		weldRate2 := Utils.FloatRound(float64(t.PwrZOk-t.PwrZWeld4)/float64(t.PwrZOk)*100, 2)

		//默认波长大于1500为10G，小于1500为1G
		if t.WaveLength > 1500 {
			t.LapWeldRate = weldRate1
			t.ThroughWeldRate = weldRate2
		} else {
			t.LapWeldRate = weldRate2
			t.ThroughWeldRate = weldRate1
		}

		res = append(res, t)
	}
	return res, nil
}

type TransmitAutoCoupleStatistic struct {
	InsName   string  //台位Id
	Pn        string  //产品pn
	Input     int     //总投入
	PassCount int     //通过数
	PassRate  float64 //通过率
}

//通过时间段和某个台位编号 pn分组统计
func GetTransmitAutoCoupleStatisticGroupByPn(insname, startTime, endTime string) ([]TransmitAutoCoupleStatistic, error) {
	var res = make([]TransmitAutoCoupleStatistic, 0)
	sqlStr := `select t.pn, COUNT(t.status) as input, SUM(case when t.status='Pass' then 1 else 0 end) as pass from superxon.autodt_transmit_autocouple t
where t.insname = '` + insname + `' and t.testtime between to_date('` + startTime + `','yyyy-mm-dd hh24:mi:ss') and to_date('` + endTime + `','yyyy-mm-dd hh24:mi:ss') group by t.pn`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var t TransmitAutoCoupleStatistic
	for rows.Next() {
		err = rows.Scan(
			&t.Pn,
			&t.Input,
			&t.PassCount,
		)
		if err != nil {
			return nil, err
		}
		t.InsName = insname
		t.Pn = strings.TrimSpace(t.Pn)
		t.PassRate = Utils.FloatRound(float64(t.PassCount)/float64(t.Input)*100, 2)
		res = append(res, t)
	}
	return res, nil
}
