package runOsaDisplayBroad

import (
	"SuperxonWebSite/Models/OsaRunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStationStatusHandler(c *gin.Context) {
	stationStatusList, err := OsaRunDisplay.GetStationStatus()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stationStatusList)
	}
}
