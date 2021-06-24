package OsaStatisticDisplay

import (
	"SuperxonWebSite/Databases"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type OsaOpticsDefectInfo struct {
	LotNo     string  `db:"lotno"`
	Id        string  `db:"id"`
	Pn        string  `db:"pn"`
	Sn        string  `db:"sn"`
	ErrorCode string  `db:"errorcode"`
	TcFlag    string  `db:"tc_flag"`
	Time      string  `db:"time"`
	Result1G  float64 `db:"result_1g"`
	Result10G float64 `db:"result_10g"`
	Zz        int     `db:"zz"`
	StationId string  `db:"stationid"`
	InsName   string  `db:"insname"`
}

//获取某段时间osaPn对应的TC1的收端失败信息
func GetOsaOpticsDefectInfoOfRxByPn(osaPn, startTime, endTime string) ([]OsaOpticsDefectInfo, error) {
	sqlStr := `select rxdata.*, rxcouple.insname from (select RX.*,p.stationid from (select i.lotno,(case when o.id<>'0' then o.id+1 else o.id end) as rxid,
o.PN,o.sn,o.errorcode,o.tc_flag,to_date(o.apd_d_date,'yyyy-mm-dd hh24:mi:ss') as 时间,o.apd_t_iop,o.apd_t_io_2g,
dense_rank()over(partition by o.pn,o.sn,o.tc_flag order by to_date(o.apd_D_date,'yyyy-mm-dd hh24:mi:ss') desc)zz
from superxon.autodt_results_opticsdata o join superxon.autodt_results_osa_tracking i on (o.pn=i.pn　and o.sn=i.sn)
where (o.pn like '` + osaPn + `%') AND to_date(o.apd_D_date,'yyyy-mm-dd hh24:mi:ss') between to_date('` + startTime + `','YYYY-MM-DD HH24:MI:SS')
and to_date('` + endTime + `','YYYY-MM-DD HH24:MI:SS') order by to_date(o.apd_D_date,'yyyy-mm-dd hh24:mi:ss') desc)RX JOIN superxon.autodt_process_log p on ( rxid=p.resultsid)
where RX.zz=1 and RX.errorcode <> '0')rxdata join 
(select * from(select t.*,dense_rank()over(partition by t.sn order by t.testtime desc)zz from superxon.autodt_recive_autocouple t) where zz=1)rxcouple on rxdata.sn=rxcouple.sn
order by errorcode`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var p OsaOpticsDefectInfo
	var res = make([]OsaOpticsDefectInfo, 0)
	for rows.Next() {
		err = rows.Scan(
			&p.LotNo,
			&p.Id,
			&p.Pn,
			&p.Sn,
			&p.ErrorCode,
			&p.TcFlag,
			&p.Time,
			&p.Result1G,
			&p.Result10G,
			&p.Zz,
			&p.StationId,
			&p.InsName)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

//获取某段时间osaPn对应的TC1的发端失败信息
func GetOsaOpticsDefectInfoOfTxByPn(osaPn, startTime, endTime string) ([]OsaOpticsDefectInfo, error) {
	sqlStr := `select LIV.*,Trans.insname from (select q.*,p.stationid from (select s.lotno,t.id,t.PN,t.sn,t.errorcode,t.tc_flag,t.testdate as 时间,t.po_io,t.po_io_10g,
dense_rank()over(partition by t.pn,t.sn,t.tc_flag order by t.testdate desc)zz
from superxon.autodt_results_liv t join superxon.autodt_results_osa_tracking s on (t.sn=s.sn and t.pn=s.pn)
where (t.pn like '` + osaPn + `%')and testdate between to_date('` + startTime + `','YYYY-MM-DD HH24:MI:SS') 
and to_date('` + endTime + `','YYYY-MM-DD HH24:MI:SS') order by testdate desc)q join superxon.autodt_process_log p on (q.id=p.resultsid and q.sn=p.sn)
where q.zz=1 and q.errorcode <> '0')LIV JOIN (select * from(select distinct substr(t.sn,0,instr(t.sn,' ')-1) as sn,substr(t.insname,0,7) AS insname,t.testtime,
dense_rank()over(partition by substr(t.sn,0,instr(t.sn,' ')-1) order by t.testtime desc)zz 
from superxon.autodt_transmit_autocouple t )Trans where Trans.zz =1)Trans on Trans.sn=LIV.SN`
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var p OsaOpticsDefectInfo
	var res = make([]OsaOpticsDefectInfo, 0)
	for rows.Next() {
		err = rows.Scan(
			&p.LotNo,
			&p.Id,
			&p.Pn,
			&p.Sn,
			&p.ErrorCode,
			&p.TcFlag,
			&p.Time,
			&p.Result1G,
			&p.Result10G,
			&p.Zz,
			&p.StationId,
			&p.InsName)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

//数据存入Redis
func RedisGetOsaOpticsDefectInfoOfRxByPn(osaPn, startTime, endTime string) ([]OsaOpticsDefectInfo, error) {
	key := "OsaOpticsDefectInfoOfRxByPn:" + osaPn + startTime + endTime
	var res = make([]OsaOpticsDefectInfo, 0)
	reBytes, err := redis.Bytes(Databases.RedisPool.Get().Do("get", key))
	if len(reBytes) > 0 {
		err = json.Unmarshal(reBytes, &res)
		if err != nil {
			return nil, err
		}
		fmt.Println("使用redis")
		return res, nil
	}

	res, err = GetOsaOpticsDefectInfoOfRxByPn(osaPn, startTime, endTime)
	if err != nil {
		return nil, err
	}

	datas, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	_, err = Databases.RedisPool.Get().Do("SET", key, datas, "NX", "EX", 60*30)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func RedisGetOsaOpticsDefectInfoOfTxByPn(osaPn, startTime, endTime string) ([]OsaOpticsDefectInfo, error) {
	key := "OsaOpticsDefectInfoOfTxByPn:" + osaPn + startTime + endTime
	var res = make([]OsaOpticsDefectInfo, 0)
	reBytes, err := redis.Bytes(Databases.RedisPool.Get().Do("get", key))
	if len(reBytes) > 0 {
		err = json.Unmarshal(reBytes, &res)
		if err != nil {
			return nil, err
		}
		fmt.Println("使用redis")
		return res, nil
	}
	res, err = GetOsaOpticsDefectInfoOfTxByPn(osaPn, startTime, endTime)
	if err != nil {
		return nil, err
	}

	datas, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	_, err = Databases.RedisPool.Get().Do("SET", key, datas, "NX", "EX", 60*30)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//统计收发端数据的错误代码分布
func GetErrorCodeDistribution(datas []OsaOpticsDefectInfo) (map[string]uint, error) {
	if len(datas) <= 0 {
		return nil, errors.New("没有OSA测试相关数据")
	}
	var res = make(map[string]uint)
	for _, data := range datas {
		res[data.ErrorCode] += 1
	}
	return res, nil
}

//统计某个错误码的散点分布图,分1G和10G
func GetErrorCodeDetailByErrorCode(datas []OsaOpticsDefectInfo, errorCode string) ([]float64, []float64, error) {
	if len(datas) <= 0 {
		return nil, nil, errors.New("没有OSA测试相关数据")
	}
	var res1g = make([]float64, 0)
	var res10g = make([]float64, 0)
	for _, data := range datas {
		if data.ErrorCode == errorCode {
			res1g = append(res1g, data.Result1G)
			res10g = append(res10g, data.Result10G)
		}
	}
	return res1g, res10g, nil
}

//统计某个错误码在各测试工位上的分布
func GetStationIdDistributionByErrorCode(datas []OsaOpticsDefectInfo, errorCode string) (map[string]uint, error) {
	if len(datas) <= 0 {
		return nil, errors.New("没有OSA测试相关数据")
	}
	var res = make(map[string]uint)
	for _, data := range datas {
		if data.ErrorCode == errorCode {
			res[data.StationId] += 1
		}
	}
	return res, nil
}

//统计某个错误码在各耦合工位上的分布
func GetInsNameDistributionByErrorCode(datas []OsaOpticsDefectInfo, errorCode string) (map[string]uint, error) {
	if len(datas) <= 0 {
		return nil, errors.New("没有OSA测试相关数据")
	}
	var res = make(map[string]uint)
	for _, data := range datas {
		if data.ErrorCode == errorCode {
			res[data.InsName] += 1
		}
	}
	return res, nil
}
