package moduleStatistic

import (
	"SuperxonWebSite/Models/ModuleStatisticDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetInputAndOutputSummaryInfoListHandler(c *gin.Context) {
	var err error
	var queryCondition ModuleStatisticDisplay.QueryCondition
	StartTime, EndTime := Utils.GetCurrentAndZeroTime()
	queryCondition.StartTime = c.DefaultQuery("startTime", StartTime)
	queryCondition.EndTime = c.DefaultQuery("endTime", EndTime)
	result, err := ModuleStatisticDisplay.Get10GLineIOSummaryInfoList(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
