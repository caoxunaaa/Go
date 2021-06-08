package moduleRunning

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

//所有模块良率告警 - workOrderType: (TRX正常品 or TRX返工制品) startTime: 2006-01-02 15:04:05 endTime: 2006-01-02 15:04:05
func GetModuleAllWaringInfoHandler(c *gin.Context) {
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")
	workOrderType := c.DefaultQuery("workOrderType", "None")
	pn := "%%"

	if workOrderType == "None" || startTime == "None" || endTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}

	res, err := ModuleRunDisplay.GetModuleWaringInfo(workOrderType, pn, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
