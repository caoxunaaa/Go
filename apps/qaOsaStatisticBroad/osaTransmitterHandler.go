package qaOsaStatisticBroad

import (
	"SuperxonWebSite/Models/OsaQaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOsaTransmitterHandler(c *gin.Context) {
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")
	if startTime == "None" || endTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	reply, err := OsaQaStatisticDisplay.GetOsaTransmitter(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, reply)
	}
}
