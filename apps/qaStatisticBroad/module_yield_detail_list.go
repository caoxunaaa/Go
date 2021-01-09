package qaStatisticBroad

import (
	"SuperxonWebSite/Models/QaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetQaPnListHandler(c *gin.Context) {
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")
	if startTime == "None" || endTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaPnList, err := QaStatisticDisplay.GetQaPnList(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaPnList)
	}
}

func GetQaStatisticInfoListHandler(c *gin.Context) {
	pn := c.DefaultQuery("pn", "None")
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")
	if pn == "None" || startTime == "None" || endTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaStatisticInfoList, err := QaStatisticDisplay.GetQaStatisticInfoList(pn, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaStatisticInfoList)
	}
}

func GetQaDefectsInfoListHandler(c *gin.Context) {
	pn := c.DefaultQuery("pn", "None")
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")
	if pn == "None" || startTime == "None" || endTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaDefectsInfoList, err := QaStatisticDisplay.GetQaDefectsInfoList(pn, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaDefectsInfoList)
	}
}
