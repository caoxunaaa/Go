package equipment

import (
	"SuperxonWebSite/Models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStationStatus(c *gin.Context) {
	stationStatusList, err := Models.GetStationStatus()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stationStatusList)
	}
}
