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
	var err error
	var qaStatisticInfoList []QaStatisticDisplay.QaStatisticInfo
	pn := c.DefaultQuery("pn", "None")
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")
	order := c.DefaultQuery("order", "None")
	if pn == "None" || startTime == "None" || endTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	if order == "TRX正常品" || order == "TRX改制返工品" {
		qaStatisticInfoList, err = QaStatisticDisplay.GetQaStatisticOrderInfoList(pn, startTime, endTime, order)
	} else {
		qaStatisticInfoList, err = QaStatisticDisplay.GetQaStatisticInfoList(pn, startTime, endTime)
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaStatisticInfoList)
	}
}

func GetQaDefectsInfoListHandler(c *gin.Context) {
	var err error
	var qaDefectsInfoList []QaStatisticDisplay.QaDefectsInfo
	pn := c.DefaultQuery("pn", "None")
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")
	order := c.DefaultQuery("workOrderType", "None")
	if pn == "None" || startTime == "None" || endTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	if order == "TRX正常品" || order == "TRX改制返工品" {
		qaDefectsInfoList, err = QaStatisticDisplay.GetQaDefectsOrderInfoList(pn, startTime, endTime, order)
	} else {
		qaDefectsInfoList, err = QaStatisticDisplay.GetQaDefectsInfoList(pn, startTime, endTime)
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaDefectsInfoList)
	}
}
