package runModuleDisplayBroad

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStationStatusHandler(c *gin.Context) {
	stationStatusList, err := ModuleRunDisplay.GetStationStatus()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stationStatusList)
	}
}
