package Utils

import (
	"time"
)

func GetCurrentAndZeroTime() (zeroTimeStr string, currentTimeStr string) {
	currentTime := time.Now()
	currentTimeStr = currentTime.Format("2006-01-02 15:04:05")
	zeroTimeStr = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), -23, 0, 0, 0, currentTime.Location()).Format("2006-01-02 15:04:05")
	return
}

type Ago struct {
	Years  int
	Months int
	Days   int
}

func GetAgoAndCurrentTime(date Ago) (AgoTimeStr string, currentTimeStr string) {
	currentTime := time.Now()
	currentTimeStr = currentTime.Format("2006-01-02 15:04:05")
	AgoTimeStr = currentTime.AddDate(date.Years, date.Months, date.Days).Format("2006-01-02 15:04:05")
	return
}
