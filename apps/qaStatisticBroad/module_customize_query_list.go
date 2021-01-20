package qaStatisticBroad

import (
	"SuperxonWebSite/Models/QaStatisticDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPnSetParamsListHandler(c *gin.Context) {
	var queryCondition QaStatisticDisplay.QueryCondition
	queryCondition.Pn = c.DefaultQuery("pn", "")
	queryCondition.WorkOrderId = c.DefaultQuery("workOrderId", "")
	queryCondition.BomId = c.DefaultQuery("bomId", "")
	queryCondition.Process = c.DefaultQuery("process", "")

	resultList, err := QaStatisticDisplay.GetPnSetParams(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resultList)
	}
}

func GetPnWorkOrderYieldsListHandler(c *gin.Context) {
	var queryCondition QaStatisticDisplay.QueryCondition
	StartTime, EndTime := Utils.GetCurrentAndZeroTime()
	queryCondition.Pn = c.DefaultQuery("pn", "")
	queryCondition.WorkOrderType = c.DefaultQuery("workOrderType", "")
	queryCondition.StartTime = c.DefaultQuery("startTime", StartTime)
	queryCondition.EndTime = c.DefaultQuery("endTime", EndTime)
	resultList, err := QaStatisticDisplay.GetPnWorkOrderYields(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resultList)
	}
}
