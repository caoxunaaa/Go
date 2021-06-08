package moduleRunning

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleAllStationStatusHandler(c *gin.Context) {
	stationStatusList, err := ModuleRunDisplay.GetStationStatus(Utils.GetCurrentAndZeroTime())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stationStatusList)
	}
}
