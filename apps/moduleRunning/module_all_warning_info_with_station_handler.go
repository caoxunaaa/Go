package moduleRunning

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleAllWaringInfoWithStationHandler(c *gin.Context) {
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")
	workOrderType := c.DefaultQuery("workOrderType", "None")
	pn := "%%"
	if workOrderType == "None" || startTime == "None" || endTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	stationStatusList, err := ModuleRunDisplay.GetModuleWaringInfoWithStation(workOrderType, pn, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stationStatusList)
	}
}
