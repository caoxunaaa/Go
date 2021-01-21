package Services

import (
	"SuperxonWebSite/Models/RunDisplay"
	"github.com/robfig/cron"
)

var timedTask *cron.Cron

func InitCron() {
	timedTask = cron.New()
	//spec1 := "0 0 */1 * * ?" //每隔1小时执行任务
	//_ = timedTask.AddFunc(spec1, func() { RunDisplay.CronGetModuleList(startTimeStr, endTimeStr) })
	//spec2 := "0 0 */1 * * ?" //每隔1小时执行任务
	//_ = timedTask.AddFunc(spec2, func() { RunDisplay.CronGetOsaList(startTimeStr, endTimeStr) })
	spec3 := "0 0 */6 * * ?" //每隔6小时执行任务
	_ = timedTask.AddFunc(spec3, func() { _, _ = RunDisplay.CronGetProjectPlanList() })
	timedTask.Start()
}

func CloseCron() {
	timedTask.Stop()
}
