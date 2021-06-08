package Services

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Models/WaringDisplay"
	"SuperxonWebSite/Utils"
	"fmt"
	"github.com/robfig/cron"
	"strconv"
	"strings"
	"time"
)

var timedTask *cron.Cron

func InitCron() error {
	timedTask = cron.New()
	var spec string

	timedInsertChartDataList(0, "4", 5)

	timedInsertWarningCountChartDataList(30, "4", 5)

	spec = "0 50 */1 * * ?" //每隔1小时执行任务
	err := timedTask.AddFunc(spec, func() { _, _ = ModuleRunDisplay.GetProjectPlanList() })
	if err != nil {
		return err
	}

	spec = "0 55 */1 * * ?" //每隔1小时执行更新工位告警任务
	err = timedTask.AddFunc(spec, func() { TimedUpdateStationWarningStatistic() })
	if err != nil {
		return err
	}

	spec = "0 30 8 * * ?"
	err = timedTask.AddFunc(spec, func() {
		KafkaPushModuleWarningInfoByClock(8)
	})
	if err != nil {
		return err
	}

	spec = "0 0 13 * * ?"
	err = timedTask.AddFunc(spec, func() {
		KafkaPushModuleWarningInfoByClock(13)
	})
	if err != nil {
		return err
	}

	spec = "0 0 17 * * ?"
	err = timedTask.AddFunc(spec, func() {
		KafkaPushModuleWarningInfoByClock(17)
	})
	if err != nil {
		return err
	}

	timedTask.Start()
	return nil
}

func CloseCron() {
	timedTask.Stop()
}

//定时发送邮箱

//定时新增每天总量率的作图数据
func timedInsertChartDataList(min int, hour string, interval int) {
	pnList, _ := Utils.GetChartsPnList()
	for indexPn, pn := range pnList {
		pnTemp := pn
		spec := strconv.Itoa(indexPn%60) + " " + strconv.Itoa(min+indexPn/60) + " " + hour + " * * ?"
		_ = timedTask.AddFunc(spec, func() {
			yesterday, today := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -1})
			temp := &WaringDisplay.ChartQueryCondition{
				Pn:        pnTemp,
				StartTime: yesterday,
				EndTime:   today}
			//todo 获取某个PN某段时间的总良率 PassRate
			totalPassRate, _ := WaringDisplay.GetPnTotalPassRate(temp)
			NeedDateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", yesterday, time.Local)
			pnChartData := &WaringDisplay.PnPassRateChartData{Pn: pnTemp, DateTime: NeedDateTime, PassRate: totalPassRate}
			_ = WaringDisplay.CreatePnChartData(pnChartData)
		})
	}
}

//定时新增每天总量率的作图数据
func timedInsertWarningCountChartDataList(min int, hour string, interval int) {
	warningClassificationList, _ := Utils.GetWarningClassificationList()
	for indexWarningClassification, warningClassification := range warningClassificationList {
		warningClassificationTemp := warningClassification
		spec := strconv.Itoa(indexWarningClassification%60) + " " + strconv.Itoa(min+indexWarningClassification/60) + " " + hour + " * * ?"
		_ = timedTask.AddFunc(spec, func() {
			yesterday, today := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -1})
			temp := &WaringDisplay.ChartQueryCondition{
				Classification: warningClassificationTemp,
				StartTime:      yesterday,
				EndTime:        today}
			//todo 获取某个PN某段时间的总良率 PassRate
			warningCount, warningTotal, _ := WaringDisplay.GetWaningCount(temp)
			NeedDateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", yesterday, time.Local)
			warningCountChartData := &WaringDisplay.WarningCountChartData{DateTime: NeedDateTime, Classification: warningClassificationTemp, Count: warningCount, Total: warningTotal}
			_ = WaringDisplay.CreateWarningCountChartData(warningCountChartData)
		})
	}
}

//定时更新工位告警的个数
func TimedUpdateStationWarningStatistic() {
	currentDate, hour := Utils.GetCurrentDateAndHour()
	oneHourAgoStr, currentTimeStr := Utils.GetCurrentTimeAndOneHourAgo()
	stationWarningList, err := ModuleRunDisplay.GetStationWarningFlag(&ModuleRunDisplay.QueryCondition{StartTime: oneHourAgoStr, EndTime: currentTimeStr})
	if err != nil {
		fmt.Println(err)
		return
	}
	//如果当天的没有则先Create,否则Update
	for index := 0; index < len(stationWarningList); index++ {
		stationWarningStatisticList, err := ModuleRunDisplay.GetStationWarningStatisticFindOne(&ModuleRunDisplay.QueryCondition{StationId: stationWarningList[index].StationId.String, StartTime: currentDate})
		fmt.Println(stationWarningStatisticList, err)
		if len(stationWarningStatisticList) <= 0 {
			fmt.Println(err)
			var statisticsEachHour string
			if stationWarningList[index].Count > 0 {
				statisticsEachHour = "1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1"
			} else {
				statisticsEachHour = "-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1"
			}
			err = ModuleRunDisplay.CreateStationWarningStatistic(&ModuleRunDisplay.StationWarningStatistic{StationId: stationWarningList[index].StationId.String, RecordDate: currentDate, StatisticsEachHour: statisticsEachHour})
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			if hour > 0 {
				fmt.Println("更新")
				stationWarningStatisticTemp := stationWarningStatisticList[0]
				//更新statisticsEachHour字符串
				count := stationWarningList[index].Count
				temp := stationWarningStatisticTemp.StatisticsEachHour
				stHourList := strings.Split(temp, ",")
				stHourList[hour] = strconv.Itoa(count)
				stationWarningStatisticTemp.StatisticsEachHour = strings.Join(stHourList, ",")
				fmt.Println(stationWarningStatisticTemp)
				err = ModuleRunDisplay.UpdateStationWarningStatistic(stationWarningStatisticTemp)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
