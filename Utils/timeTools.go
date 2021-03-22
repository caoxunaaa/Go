package Utils

import (
	"time"
)

//当天0点到当前的时间
func GetCurrentAndZeroTime() (zeroTimeStr string, currentTimeStr string) {
	currentTime := time.Now()
	currentTimeStr = currentTime.Format("2006-01-02 15:04:05")
	zeroTimeStr = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Format("2006-01-02 15:04:05")
	return
}

//当月第一天0点到当前时间
func GetCurrentAndZeroDayTime() (zeroTimeStr string, currentTimeStr string) {
	currentTime := time.Now()
	currentTimeStr = currentTime.Format("2006-01-02 15:04:05")
	zeroTimeStr = time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, currentTime.Location()).Format("2006-01-02 15:04:05")
	return
}

type Ago struct {
	Years  int
	Months int
	Days   int
}

//Ago之前的当前时分秒到当前时间
func GetAgoAndCurrentTime(date Ago) (AgoTimeStr string, currentTimeStr string) {
	currentTime := time.Now()
	currentTimeStr = currentTime.Format("2006-01-02 15:04:05")
	AgoTimeStr = currentTime.AddDate(date.Years, date.Months, date.Days).Format("2006-01-02 15:04:05")
	return
}

//Ago之前的0点到当天的0点
func GetAgoAndCurrentTimeZero(date Ago) (AgoTimeStr string, currentTimeStr string) {
	currentTime := time.Now()
	//currentTimeStr = currentTime.Format("2006-01-02 15:04:05")
	currentTimeStr = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Format("2006-01-02 15:04:05")
	AgoTime := currentTime.AddDate(date.Years, date.Months, date.Days)
	AgoTimeStr = time.Date(AgoTime.Year(), AgoTime.Month(), AgoTime.Day(), 0, 0, 0, 0, AgoTime.Location()).Format("2006-01-02 15:04:05")
	return
}

//获取当前日期和时段
func GetCurrentDateAndHour() (date string, hour int) {
	currentTime := time.Now()
	date = currentTime.Format("2006-01-02")
	hour = currentTime.Hour()
	return
}

//获取当前日期和时段
func GetCurrentTimeAndOneHourAgo() (oneHourAgoStr string, currentTimeStr string) {
	currentTime := time.Now()
	currentTimeStr = currentTime.Format("2006-01-02 15:04:05")
	oneHourAgoStr = currentTime.Add(-time.Hour * 1).Format("2006-01-02 15:04:05")
	return
}
