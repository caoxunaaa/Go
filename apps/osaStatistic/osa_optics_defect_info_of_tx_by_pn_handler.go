package osaStatistic

import (
	"SuperxonWebSite/Models/OsaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOsaOpticsDefectInfoOfTxByOsaPnHandler(c *gin.Context) {
	osaPn := c.DefaultQuery("osaPn", "None")
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")

	if startTime == "None" || endTime == "None" || osaPn == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	res, err := OsaStatisticDisplay.GetOsaOpticsDefectInfoOfTxByPn(osaPn, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
