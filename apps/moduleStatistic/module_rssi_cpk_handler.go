package moduleStatistic

import (
	"SuperxonWebSite/Models/ModuleStatisticDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleRssiCpkHandler(c *gin.Context) {
	var err error
	var queryCondition ModuleStatisticDisplay.QueryCondition
	StartTime, EndTime := Utils.GetCurrentAndZeroTime()
	queryCondition.Pn = c.DefaultQuery("pn", "")
	queryCondition.Process = c.DefaultQuery("process", "")
	queryCondition.StartTime = c.DefaultQuery("startTime", StartTime)
	queryCondition.EndTime = c.DefaultQuery("endTime", EndTime)
	queryCondition.WorkOrderType = c.DefaultQuery("workOrderType", "")
	result, err := ModuleStatisticDisplay.GetQaCpkRssiList(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
