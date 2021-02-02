package Services

import (
	"SuperxonWebSite/Models/ModuleQaStatisticDisplay"
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/robfig/cron"
	"strconv"
)

var timedTask *cron.Cron

func InitCron() {
	timedTask = cron.New()

	timedGetCpkInfo(30, "1", 5)

	timedGetCpkRssi(1, "2", 5)

	timedGetQaStatisticOrderInfo(30, "2", 5)

	timedGetQaDefectsOrderInfoListByPn(01, "3", 5)

	spec3 := "0 0 */6 * * ?" //每隔6小时执行任务
	_ = timedTask.AddFunc(spec3, func() { _, _ = ModuleRunDisplay.RedisGetProjectPlanList() })
	timedTask.Start()
}

func CloseCron() {
	timedTask.Stop()
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

			spec2 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec2, func() {
				halfMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -15})
				temp2 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:        pnTemp,
					StartTime: halfMonthAgo,
					EndTime:   now,
					Process:   processTemp.Name}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaCpkInfoList(temp2)
			})

			spec3 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval*2+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec3, func() {
				oneMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -30})
				temp3 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:        pnTemp,
					StartTime: oneMonthAgo,
					EndTime:   now,
					Process:   processTemp.Name}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaCpkInfoList(temp3)
			})
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

			spec2 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec2, func() {
				halfMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -15})
				temp2 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:        pnTemp,
					StartTime: halfMonthAgo,
					EndTime:   now,
					Process:   processTemp.Name}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaCpkRssiList(temp2)
			})

			spec3 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval*2+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec3, func() {
				oneMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -30})
				temp3 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:        pnTemp,
					StartTime: oneMonthAgo,
					EndTime:   now,
					Process:   processTemp.Name}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaCpkRssiList(temp3)
			})
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

			spec2 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec2, func() {
				halfMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -15})
				temp2 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:            pnTemp,
					StartTime:     halfMonthAgo,
					EndTime:       now,
					WorkOrderType: processTemp}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaStatisticOrderInfoList(temp2)
			})

			spec3 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval*2+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec3, func() {
				oneMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -30})
				temp3 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:            pnTemp,
					StartTime:     oneMonthAgo,
					EndTime:       now,
					WorkOrderType: processTemp}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaStatisticOrderInfoList(temp3)
			})
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

			spec2 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec2, func() {
				halfMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -15})
				temp2 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:            pnTemp,
					StartTime:     halfMonthAgo,
					EndTime:       now,
					WorkOrderType: processTemp}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaDefectsOrderInfoListByPn(temp2)
			})

			spec3 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(min+interval*2+indexProcess/60+indexPn%60) + " " + hour + " * * ?"
			_ = timedTask.AddFunc(spec3, func() {
				oneMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -30})
				temp3 := &ModuleQaStatisticDisplay.QueryCondition{
					Pn:            pnTemp,
					StartTime:     oneMonthAgo,
					EndTime:       now,
					WorkOrderType: processTemp}
				_, _ = ModuleQaStatisticDisplay.RedisGetQaDefectsOrderInfoListByPn(temp3)
			})
		}
	}
}
