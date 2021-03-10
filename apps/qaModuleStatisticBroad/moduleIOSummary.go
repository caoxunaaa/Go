package qaModuleStatisticBroad

import (
	"SuperxonWebSite/Models/ModuleQaStatisticDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Get10GLineIOSummaryInfoListHandler(c *gin.Context) {
	var err error
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	StartTime, EndTime := Utils.GetCurrentAndZeroTime()
	queryCondition.StartTime = c.DefaultQuery("startTime", StartTime)
	queryCondition.EndTime = c.DefaultQuery("endTime", EndTime)
	result, err := ModuleQaStatisticDisplay.Get10GLineIOSummaryInfoList(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
