package qaStatisticBroad

import (
	"SuperxonWebSite/Models/QaStatisticDisplay"
	"SuperxonWebSite/Utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPnSetParamsListHandler(c *gin.Context) {
	var queryCondition QaStatisticDisplay.QueryCondition

	queryCondition.Pn = c.DefaultQuery("Pn", "")
	queryCondition.WorkOrderId = c.DefaultQuery("WorkOrderId", "")
	queryCondition.BomId = c.DefaultQuery("BomId", "")
	queryCondition.Process = c.DefaultQuery("Process", "")
	fmt.Println(queryCondition)

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
	queryCondition.Pn = c.DefaultQuery("Pn", "")
	queryCondition.WorkOrderType = c.DefaultQuery("WorkOrderType", "")
	queryCondition.StartTime = c.DefaultQuery("StartTime", StartTime)
	queryCondition.EndTime = c.DefaultQuery("EndTime", EndTime)
	fmt.Println(queryCondition)

	resultList, err := QaStatisticDisplay.GetPnWorkOrderYields(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resultList)
	}
}
