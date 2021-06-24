package moduleRunning

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleAllInfoWithStationTodayHandler(c *gin.Context) {
	startTime, endTime := Utils.GetCurrentAndZeroTime()
	stationStatusList, err := ModuleRunDisplay.GetStationStatus("%%", "", startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stationStatusList)
	}
}
