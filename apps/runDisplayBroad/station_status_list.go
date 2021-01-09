package runDisplayBroad

import (
	"SuperxonWebSite/Models/RunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStationStatusHandler(c *gin.Context) {
	stationStatusList, err := RunDisplay.GetStationStatus()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stationStatusList)
	}
}
