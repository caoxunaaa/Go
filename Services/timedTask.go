package Services

import (
	"SuperxonWebSite/Models/QaStatisticDisplay"
	"SuperxonWebSite/Models/RunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/robfig/cron"
	"strconv"
)

var timedTask *cron.Cron

func InitCron() {
	timedTask = cron.New()

	timedGetCpkInfo()

	timedGetCpkRssi()

	spec3 := "0 0 */6 * * ?" //每隔6小时执行任务
	_ = timedTask.AddFunc(spec3, func() { _, _ = RunDisplay.CronGetProjectPlanList() })
	timedTask.Start()
}

func CloseCron() {
	timedTask.Stop()
}

/*
定时在凌晨1点获取一周，半个月，一个月的CPKInfo存入缓存
*/
func timedGetCpkInfo() {
	pnList, _ := Utils.GetCommonPnList()
	processList, _ := QaStatisticDisplay.GetAllProcessOfTRX()

	oneWeekAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -7})
	halfMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -15})
	oneMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -30})

	for indexPn, pn := range pnList {
		pnTemp := pn
		for indexProcess, process := range processList {
			processTemp := process
			temp1 := &QaStatisticDisplay.QueryCondition{
				Pn:        pnTemp,
				StartTime: oneWeekAgo,
				EndTime:   now,
				Process:   processTemp.Name}
			spec1 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(indexProcess/60+indexPn%60) + " 1 * * ?"
			_ = timedTask.AddFunc(spec1, func() {
				_, _ = QaStatisticDisplay.CronGetQaCpkInfoList(temp1)
			})

			temp2 := &QaStatisticDisplay.QueryCondition{
				Pn:        pnTemp,
				StartTime: halfMonthAgo,
				EndTime:   now,
				Process:   processTemp.Name}
			spec2 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(20+indexProcess/60+indexPn%60) + " 1 * * ?"
			_ = timedTask.AddFunc(spec2, func() {
				_, _ = QaStatisticDisplay.CronGetQaCpkInfoList(temp2)
			})

			temp3 := &QaStatisticDisplay.QueryCondition{
				Pn:        pnTemp,
				StartTime: oneMonthAgo,
				EndTime:   now,
				Process:   processTemp.Name}
			spec3 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(40+indexProcess/60+indexPn%60) + " 1 * * ?"
			_ = timedTask.AddFunc(spec3, func() {
				_, _ = QaStatisticDisplay.CronGetQaCpkInfoList(temp3)
			})
		}
	}
}

/*
定时在凌晨2点获取一周，半个月，一个月的CPKRssi存入缓存
*/
func timedGetCpkRssi() {
	pnList, _ := Utils.GetCommonPnList()
	processList, _ := QaStatisticDisplay.GetAllProcessOfTRX()

	oneWeekAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -7})
	halfMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -15})
	oneMonthAgo, now := Utils.GetAgoAndCurrentTimeZero(Utils.Ago{Years: 0, Months: 0, Days: -30})

	for indexPn, pn := range pnList {
		pnTemp := pn
		for indexProcess, process := range processList {
			processTemp := process
			temp1 := &QaStatisticDisplay.QueryCondition{
				Pn:        pnTemp,
				StartTime: oneWeekAgo,
				EndTime:   now,
				Process:   processTemp.Name}
			spec1 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(indexProcess/60+indexPn%60) + " 2 * * ?"
			_ = timedTask.AddFunc(spec1, func() {
				_, _ = QaStatisticDisplay.CronGetQaCpkRssiList(temp1)
			})

			temp2 := &QaStatisticDisplay.QueryCondition{
				Pn:        pnTemp,
				StartTime: halfMonthAgo,
				EndTime:   now,
				Process:   processTemp.Name}
			spec2 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(20+indexProcess/60+indexPn%60) + " 2 * * ?"
			_ = timedTask.AddFunc(spec2, func() {
				_, _ = QaStatisticDisplay.CronGetQaCpkRssiList(temp2)
			})

			temp3 := &QaStatisticDisplay.QueryCondition{
				Pn:        pnTemp,
				StartTime: oneMonthAgo,
				EndTime:   now,
				Process:   processTemp.Name}
			spec3 := strconv.Itoa(indexProcess%60) + " " + strconv.Itoa(40+indexProcess/60+indexPn%60) + " 2 * * ?"
			_ = timedTask.AddFunc(spec3, func() {
				_, _ = QaStatisticDisplay.CronGetQaCpkRssiList(temp3)
			})
		}
	}
}
