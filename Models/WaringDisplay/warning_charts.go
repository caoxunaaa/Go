// @Title  warning_charts.go
// @Description  处理获取各类警告图表的数据
// @Author  曹迅 (时间 2021/02/04  16:00)
// @Update  曹迅 (时间 2021/02/04  16:00)
package WaringDisplay

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Models/OsaRunDisplay"
	"strconv"
	"strings"
	"time"
)

type ChartQueryCondition struct {
	Pn             string //产品型号
	StartTime      string //开始时间
	EndTime        string //结束时间
	Classification string //对警告目标的分类,比如OSA直通率，模块直通率，OSA工位良率，模块工位良率等等
}

//PN良率作图数据（时间间隔为1天）
type PnPassRateChartData struct {
	ID       uint      `gorm:"primary_key" db:"id"`
	Pn       string    `db:"pn"`
	DateTime time.Time `db:"date_time"`
	PassRate string    `db:"pass_rate"`
}

//todo 获取某个PN某段时间总良率,没有的情况写入0%！！
func GetPnTotalPassRate(chartQueryCondition *ChartQueryCondition) (totalPassRate string, err error) {
	totalPassRate = "0%"
	return
}

//获取某个PN某段时间的良率作图数据（间隔为1天）
func GetPnChartDataList(chartQueryCondition *ChartQueryCondition) (pnChartDataList []*PnPassRateChartData, err error) {
	sqlStr := "SELECT * FROM production_yield_daily WHERE pn=? and date_time between ? and ? order by date_time"
	err = Databases.SuperxonProductionLineProductStatisticDevice.Select(&pnChartDataList, sqlStr, chartQueryCondition.Pn, chartQueryCondition.StartTime, chartQueryCondition.EndTime)
	if err != nil {
		return nil, err
	}
	return
}

//每天PN的总量率数据插入到数据库中
func CreatePnChartData(pnChartData *PnPassRateChartData) (err error) {
	sqlStr := "INSERT INTO production_yield_daily(pn, date_time, pass_rate) values (?, ?, ?)"
	_, err = Databases.SuperxonProductionLineProductStatisticDevice.Exec(sqlStr,
		pnChartData.Pn,
		pnChartData.DateTime,
		pnChartData.PassRate)
	if err != nil {
		return err
	}
	return
}

//警告次数作图数据（时间间隔为1天）
type WarningCountChartData struct {
	ID             uint      `gorm:"primary_key" db:"id"`
	DateTime       time.Time `db:"date_time"`
	Classification string    `db:"classification"` //对警告目标的分类,比如OSA直通率，模块直通率，OSA工位良率，模块工位良率等等
	Count          int       `db:"count"`
	Total          int       `db:"total"`
}

//获取警告次数，低于90%警告
func GetWaningCount(chartQueryCondition *ChartQueryCondition) (warningCount, warningTotal int, err error) {
	warningCount = 0
	warningTotal = 0
	switch chartQueryCondition.Classification {
	case "模块直通率":
		moduleInfoList, err1 := ModuleRunDisplay.GetAllModuleInfoList(chartQueryCondition.StartTime, chartQueryCondition.EndTime)
		if err1 != nil {
			return
		}
		for _, moduleInfo := range moduleInfoList {
			temp1 := strings.Split(moduleInfo.OncePassRate, "%")
			temp2, _ := strconv.Atoi(temp1[0])
			if temp2 < 90 {
				warningCount++
			}
			warningTotal++
		}
	case "OSA直通率":
		osaInfoList, err1 := OsaRunDisplay.GetAllOsaInfoList(&OsaRunDisplay.OsaQueryCondition{StartTime: chartQueryCondition.StartTime, EndTime: chartQueryCondition.EndTime})
		if err1 != nil {
			return
		}
		for _, osaInfo := range osaInfoList {
			temp1 := strings.Split(osaInfo.OncePassRate, "%")
			temp2, _ := strconv.Atoi(temp1[0])
			if temp2 < 85 {
				warningCount++
			}
			warningTotal++
		}
	case "模块工位良率":
		stationStatusList, err1 := ModuleRunDisplay.GetStationStatus(chartQueryCondition.StartTime, chartQueryCondition.EndTime)
		if err1 != nil {
			return
		}
		for _, stationStatus := range stationStatusList {
			temp1 := strings.Split(stationStatus.PassRate, "%")
			temp2, _ := strconv.Atoi(temp1[0])
			if temp2 < 90 {
				warningCount++
			}
			warningTotal++
		}
	case "OSA工位良率":
		stationStatusList, err1 := OsaRunDisplay.GetStationStatus(chartQueryCondition.StartTime, chartQueryCondition.EndTime)
		if err1 != nil {
			return
		}
		for _, stationStatus := range stationStatusList {
			temp1 := strings.Split(stationStatus.PassRate, "%")
			temp2, _ := strconv.Atoi(temp1[0])
			if temp2 < 85 {
				warningCount++
			}
			warningTotal++
		}
	}
	return
}

//获取某段时间某个分类的警告作图数据（间隔为1天）
func GetWarningCountChartDataList(chartQueryCondition *ChartQueryCondition) (warningCountChartData []*WarningCountChartData, err error) {
	sqlStr := "SELECT * FROM statistic_warning_count_daily WHERE classification=? and date_time between ? and ? order by date_time"
	err = Databases.SuperxonProductionLineProductStatisticDevice.Select(&warningCountChartData, sqlStr, chartQueryCondition.Classification, chartQueryCondition.StartTime, chartQueryCondition.EndTime)
	if err != nil {
		return nil, err
	}
	return
}

//插入每天告警作图数据到数据库中
func CreateWarningCountChartData(warningCountChartData *WarningCountChartData) (err error) {
	sqlStr := "INSERT INTO statistic_warning_count_daily(date_time, classification, count, total) values (?, ?, ?, ?)"
	_, err = Databases.SuperxonProductionLineProductStatisticDevice.Exec(sqlStr,
		warningCountChartData.DateTime,
		warningCountChartData.Classification,
		warningCountChartData.Count,
		warningCountChartData.Total)
	if err != nil {
		return err
	}
	return
}
