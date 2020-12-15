package Utils

import (
	"time"
)

func GetCurrentAndZeroTime() (zeroTimeStr string, currentTimeStr string) {
	currentTime := time.Now()
	currentTimeStr = currentTime.Format("2006-01-02 15:04:05")
	zeroTimeStr = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Format("2006-01-02 15:04:05")
	return
}

func GetTwoMonthAgoAndCurrentTime() (twoMonthAgoTimeStr string, currentTimeStr string) {
	currentTime := time.Now()
	currentTimeStr = currentTime.Format("2006-01-02 15:04:05")
	twoMonthAgoTimeStr = currentTime.AddDate(0, -2, 0).Format("2006-01-02 15:04:05")
	return
}
