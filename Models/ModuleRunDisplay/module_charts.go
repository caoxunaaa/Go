// @Title  module_charts.go
// @Description  处理获取各类图表的数据
// @Author  曹迅 (时间 2021/02/04  16:00)
// @Update  曹迅 (时间 2021/02/04  16:00)
package ModuleRunDisplay

import (
	"SuperxonWebSite/Databases"
	"time"
)

type ChartQueryCondition struct {
	Pn        string //产品型号
	StartTime string //开始时间
	EndTime   string //结束时间
}

type PnChartData struct {
	ID       uint      `gorm:"primary_key" db:"id"`
	Pn       string    `db:"pn"`
	DateTime time.Time `db:"date_time"`
	PassRate string    `db:"pass_rate"`
}

//todo 获取某个PN某段时间总良率,没有的情况写入0%！！
func GetPnTotalPassRate(chartQueryCondition *ChartQueryCondition) (totalPassRate string, err error) {
	return
}

//获取某个PN某段时间的良率（间隔为1天）
func GetPnChartDataList(chartQueryCondition *ChartQueryCondition) (pnChartDataList []*PnChartData, err error) {
	sqlStr := "SELECT * FROM pn_chart_data WHERE pn=? and date_time between ? and ? order by date_time"
	err = Databases.SuperxonDbDevice.Select(&pnChartDataList, sqlStr, chartQueryCondition.Pn, chartQueryCondition.StartTime, chartQueryCondition.EndTime)
	if err != nil {
		return nil, err
	}
	return
}

//每天PN的总量率数据插入到数据库中
func CreatePnChartData(pnChartData *PnChartData) (err error) {
	sqlStr := "INSERT INTO pn_chart_data(pn, date_time, pass_rate) values (?, ?, ?)"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr,
		pnChartData.Pn,
		pnChartData.DateTime,
		pnChartData.PassRate)
	if err != nil {
		return err
	}
	return
}
