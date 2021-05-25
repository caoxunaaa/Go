package runModuleDisplayBroad

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStationStatusHandler(c *gin.Context) {
	stationStatusList, err := ModuleRunDisplay.GetStationStatus(Utils.GetCurrentAndZeroTime())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stationStatusList)
	}
}

func GetStationWarningStatisticHandler(c *gin.Context) {
	currentDate, _ := Utils.GetCurrentDateAndHour()
	yesterdayDate, _ := Utils.GetYesterdayDateAndHour()
	stationWarningStatistic, err := ModuleRunDisplay.GetStationWarningStatisticFindAll(&ModuleRunDisplay.QueryCondition{StartTime: currentDate})
	stationWarningStatisticYesterday, err := ModuleRunDisplay.GetStationWarningStatisticFindAll(&ModuleRunDisplay.QueryCondition{StartTime: yesterdayDate})
	stationWarningStatistic = append(stationWarningStatistic, stationWarningStatisticYesterday...)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": errors.New("当前还没有开始记录数据")})
	} else {
		c.JSON(http.StatusOK, stationWarningStatistic)
	}
}
