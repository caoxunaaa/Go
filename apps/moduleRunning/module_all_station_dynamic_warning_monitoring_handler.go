package moduleRunning

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleAllStationDynamicWarningMonitoringHandler(c *gin.Context) {
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
