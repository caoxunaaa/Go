package Services

import (
	"SuperxonWebSite/Models/RunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/robfig/cron"
)

var timedTask *cron.Cron

func InitCron() {
	timedTask = cron.New()

	startTimeStr, endTimeStr := Utils.GetAgoAndCurrentTime(Utils.Ago{Days: -10})
	spec1 := "0 0 1 * * ?" //每天凌晨一点执行任务
	_ = timedTask.AddFunc(spec1, func() { RunDisplay.CronGetModuleList(startTimeStr, endTimeStr) })
	spec2 := "0 0 2 * * ?" //每天凌晨两点点执行任务
	_ = timedTask.AddFunc(spec2, func() { RunDisplay.CronGetOsaList(startTimeStr, endTimeStr) })
	spec3 := "0 0 3 * * ?" //每天凌晨三点点执行任务
	_ = timedTask.AddFunc(spec3, func() { _, _ = RunDisplay.CronGetProjectPlanList() })
	timedTask.Start()
}

func CloseCron() {
	timedTask.Stop()
}
