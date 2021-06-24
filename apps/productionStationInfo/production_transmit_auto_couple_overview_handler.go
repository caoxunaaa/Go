package productionStationInfo

import (
	"SuperxonWebSite/Models/ProductionLineStation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProductionTransmitAutoCoupleOverviewHandler(c *gin.Context) {
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")

	if startTime == "None" || endTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	res, err := ProductionLineStation.GetAllTransmitAutoCoupleStationOverView(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
