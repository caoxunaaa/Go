package qaModuleStatisticBroad

import (
	"SuperxonWebSite/Models/ModuleQaStatisticDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetQaCpkInfoListHandler(c *gin.Context) {
	var err error
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	StartTime, EndTime := Utils.GetCurrentAndZeroTime()
	queryCondition.Pn = c.DefaultQuery("pn", "")
	queryCondition.Process = c.DefaultQuery("process", "")
	queryCondition.StartTime = c.DefaultQuery("startTime", StartTime)
	queryCondition.EndTime = c.DefaultQuery("endTime", EndTime)
	queryCondition.WorkOrderType = c.DefaultQuery("workOrderType", "")
	result, err := ModuleQaStatisticDisplay.RedisGetQaCpkInfoList(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func GetQaCpkRssiListHandler(c *gin.Context) {
	var err error
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	StartTime, EndTime := Utils.GetCurrentAndZeroTime()
	queryCondition.Pn = c.DefaultQuery("pn", "")
	queryCondition.Process = c.DefaultQuery("process", "")
	queryCondition.StartTime = c.DefaultQuery("startTime", StartTime)
	queryCondition.EndTime = c.DefaultQuery("endTime", EndTime)
	queryCondition.WorkOrderType = c.DefaultQuery("workOrderType", "")
	result, err := ModuleQaStatisticDisplay.RedisGetQaCpkRssiList(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
