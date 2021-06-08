package moduleStatistic

import (
	"SuperxonWebSite/Models/ModuleQaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleAllPnByTimeHandler(c *gin.Context) {
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	queryCondition.StartTime = c.DefaultQuery("startTime", "None")
	queryCondition.EndTime = c.DefaultQuery("endTime", "None")

	if queryCondition.StartTime == "None" || queryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaPnList, err := ModuleQaStatisticDisplay.GetQaPnList(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaPnList)
	}
}
