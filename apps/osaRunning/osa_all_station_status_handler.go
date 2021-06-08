package osaRunning

import (
	"SuperxonWebSite/Models/OsaRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOsaAllStationStatusHandler(c *gin.Context) {
	stationStatusList, err := OsaRunDisplay.GetStationStatus(Utils.GetCurrentAndZeroTime())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stationStatusList)
	}
}
