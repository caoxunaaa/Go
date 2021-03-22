package Services

import (
	"SuperxonWebSite/Models/DeviceManage"
	"SuperxonWebSite/Models/ModuleQaStatisticDisplay"
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

func InitCron() {
	timedTask = cron.New()

	timedGetCpkInfo(30, "1", 5)

	timedGetCpkRssi(1, "2", 5)

	//timedGetQaStatisticOrderInfo(30, "2", 5)

	//timedGetQaDefectsOrderInfoListByPn(1, "3", 5)

	timedInsertChartDataList(0, "4", 5)

	timedInsertWarningCountChartDataList(30, "4", 5)

	spec3 := "0 50 */1 * * ?" //每隔1小时执行任务
	_ = timedTask.AddFunc(spec3, func() { _, _ = ModuleRunDisplay.RedisGetProjectPlanList() })
	//timedTask.Start()

	spec4 := "0 0 5 * * ?" //每天5点执行保养更新任务
	_ = timedTask.AddFunc(spec4, func() {
		err := DeviceManage.CronUpdateDeviceMainenanceStatus()
		if err != nil {
			return
		}
	})

	spec5 := "0 10 6 * * ?" //每天5点30执行保养更新任务
	_ = timedTask.AddFunc(spec5, func() {
		err := DeviceManage.CronUpdateDeviceBaseMainenanceInfo()
		if err != nil {
			return
		}
	})

	spec6 := "0 0 */1 * * ?" //每隔1小时执行任务
	_ = timedTask.AddFunc(spec6, func() { TimedUpdateStationWarningStatistic() })

	timedTask.Start()
}

func CloseCron() {
	timedTask.Stop()
}

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
			statisticsEachHour := strconv.Itoa(stationWarningList[index].Count) + ",-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1"
			err = ModuleRunDisplay.CreateStationWarningStatistic(&ModuleRunDisplay.StationWarningStatistic{StationId: stationWarningList[index].StationId.String, RecordDate: currentDate, StatisticsEachHour: statisticsEachHour})
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
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

/*
定时在凌晨1点获取一周，半个月，一个月的CPKInfo存入缓存
*/
func timedGetCpkInfo(min int, hour string, interval int) {
	pnList, _ := Utils.GetCommonPnList()
	processList, _ := ModuleQaStatisticDisplay.GetAllProcessOfTRX()
	for indexPn, pn := range pnList {
		pnTemp := pn
		for indexProcess, process := range processList {
			processTemp := process
			spec1 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec1, func() {
				oneWeekAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -7})
				temp1 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:        pnTemp,
					StartTime: oneWeekAgo,
					EndTime:   now,
					Process:   processTemp.Name}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaCpkInfoList(temp1)
			})

			//spec2 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			//_ = timedTask.AddFunc(spec2, func() {
			//	halfMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -15})
			//	temp2 := &ModuleQaStatisticDisplay.QueryCondition{
			//		Pn:        pnTemp,
			//		StartTime: halfMonthAgo,
			//		EndTime:   now,
			//		Process:   processTemp.Name}
			//	_, _ = ModuleQaStatisticDisplay.RedisGetQaCpkInfoList(temp2)
			//})

			//spec3 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval*2+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			//_ = timedTask.AddFunc(spec3, func() {
			//	oneMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -30})
			//	temp3 := &ModuleQaStatisticDisplay.QueryCondition{
			//		Pn:        pnTemp,
			//		StartTime: oneMonthAgo,
			//		EndTime:   now,
			//		Process:   processTemp.Name}
			//	_, _ = ModuleQaStatisticDisplay.RedisGetQaCpkInfoList(temp3)
			//})
		}
	}
}

/*
定时在凌晨2点获取一周，半个月，一个月的CPKRssi存入缓存
*/
func timedGetCpkRssi(min int, hour string, interval int) {
	pnList, _ := Utils.GetCommonPnList()
	processList, _ := ModuleQaStatisticDisplay.GetAllProcessOfTRX()

	for indexPn, pn := range pnList {
		pnTemp := pn
		for indexProcess, process := range processList {
			processTemp := process

			spec1 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec1, func() {
				oneWeekAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -7})
				temp1 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:        pnTemp,
					StartTime: oneWeekAgo,
					EndTime:   now,
					Process:   processTemp.Name}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaCpkRssiList(temp1)
			})

			//spec2 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			//_ = timedTask.AddFunc(spec2, func() {
			//	halfMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -15})
			//	temp2 := &ModuleQaStatisticDisplay.QueryCondition{
			//		Pn:        pnTemp,
			//		StartTime: halfMonthAgo,
			//		EndTime:   now,
			//		Process:   processTemp.Name}
			//	_, _ = ModuleQaStatisticDisplay.RedisGetQaCpkRssiList(temp2)
			//})

			//spec3 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval*2+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			//_ = timedTask.AddFunc(spec3, func() {
			//	oneMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -30})
			//	temp3 := &ModuleQaStatisticDisplay.QueryCondition{
			//		Pn:        pnTemp,
			//		StartTime: oneMonthAgo,
			//		EndTime:   now,
			//		Process:   processTemp.Name}
			//	_, _ = ModuleQaStatisticDisplay.RedisGetQaCpkRssiList(temp3)
			//})
		}
	}
}

/*
定时在凌晨3点获取一周，半个月，一个月的QaStatisticOrderInfo存入缓存
*/
func timedGetQaStatisticOrderInfo(min int, hour string, interval int) {
	pnList, _ := Utils.GetCommonPnList()
	processList := []string{"", "TRX正常品", "TRX改制返工品"}

	for indexPn, pn := range pnList {
		pnTemp := pn
		for indexProcess, process := range processList {
			processTemp := process

			spec1 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec1, func() {
				oneWeekAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -7})
				temp1 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:            pnTemp,
					StartTime:     oneWeekAgo,
					EndTime:       now,
					WorkOrderType: processTemp}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaStatisticOrderInfoList(temp1)
			})

			//spec2 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			//_ = timedTask.AddFunc(spec2, func() {
			//	halfMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -15})
			//	temp2 := &ModuleQaStatisticDisplay.QueryCondition{
			//		Pn:            pnTemp,
			//		StartTime:     halfMonthAgo,
			//		EndTime:       now,
			//		WorkOrderType: processTemp}
			//	_, _ = ModuleQaStatisticDisplay.RedisGetQaStatisticOrderInfoList(temp2)
			//})

			//spec3 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval*2+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			//_ = timedTask.AddFunc(spec3, func() {
			//	oneMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -30})
			//	temp3 := &ModuleQaStatisticDisplay.QueryCondition{
			//		Pn:            pnTemp,
			//		StartTime:     oneMonthAgo,
			//		EndTime:       now,
			//		WorkOrderType: processTemp}
			//	_, _ = ModuleQaStatisticDisplay.RedisGetQaStatisticOrderInfoList(temp3)
			//})
		}
	}
}

/*
定时在凌晨4点获取一周，半个月，一个月的GetQaDefectsOrderInfoListByPn存入缓存
*/
func timedGetQaDefectsOrderInfoListByPn(min int, hour string, interval int) {
	pnList, _ := Utils.GetCommonPnList()
	processList := []string{"", "TRX正常品", "TRX改制返工品"}

	for indexPn, pn := range pnList {
		pnTemp := pn
		for indexProcess, process := range processList {
			processTemp := process
			spec1 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec1, func() {
				oneWeekAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -7})
				temp1 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:        pnTemp,
					StartTime: oneWeekAgo,
					EndTime:   now,
					Process:   processTemp}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaDefectsOrderInfoListByPn(temp1)
			})

			//spec2 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			//_ = timedTask.AddFunc(spec2, func() {
			//	halfMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -15})
			//	temp2 := &ModuleQaStatisticDisplay.QueryCondition{
			//		Pn:            pnTemp,
			//		StartTime:     halfMonthAgo,
			//		EndTime:       now,
			//		WorkOrderType: processTemp}
			//	_, _ = ModuleQaStatisticDisplay.RedisGetQaDefectsOrderInfoListByPn(temp2)
			//})

			//spec3 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval*2+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			//_ = timedTask.AddFunc(spec3, func() {
			//	oneMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -30})
			//	temp3 := &ModuleQaStatisticDisplay.QueryCondition{
			//		Pn:            pnTemp,
			//		StartTime:     oneMonthAgo,
			//		EndTime:       now,
			//		WorkOrderType: processTemp}
			//	_, _ = ModuleQaStatisticDisplay.RedisGetQaDefectsOrderInfoListByPn(temp3)
			//})
		}
	}
}
