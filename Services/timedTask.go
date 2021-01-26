package Services

import (
	"SuperxonWebSite/Models/RunDisplay"
	"SuperxonWebSite/Utils"
	"fmt"
	"github.com/robfig/cron"
)

var timedTask *cron.Cron

func InitCron() {
	timedTask = cron.New()

	pnList, _ := Utils.GetCommonPnList()
	fmt.Println(pnList)

	//spec1 := "0 0 3 * * ?" //凌晨3点执行CPKRSSI CPKBASE任务
	//_ = timedTask.AddFunc(spec1, func() { QaStatisticDisplay.CronGetQaCpkInfoList(startTimeStr, endTimeStr) })
	//spec2 := "0 0 */1 * * ?" //每隔1小时执行任务
	//_ = timedTask.AddFunc(spec2, func() { RunDisplay.CronGetOsaList(startTimeStr, endTimeStr) })
	spec3 := "0 0 */6 * * ?" //每隔6小时执行任务
	_ = timedTask.AddFunc(spec3, func() { _, _ = RunDisplay.CronGetProjectPlanList() })
	timedTask.Start()
}

func CloseCron() {
	timedTask.Stop()
}
