package moduleStatistic

import (
	"SuperxonWebSite/Models/ModuleQaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleWorkOrderIdHandler(c *gin.Context) {
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	queryCondition.Pn = c.DefaultQuery("pn", "None")
	queryCondition.WorkOrderType = c.DefaultQuery("workOrderType", "")
	queryCondition.StartTime = c.DefaultQuery("startTime", "None")
	queryCondition.EndTime = c.DefaultQuery("endTime", "None")
	isFinish := c.DefaultQuery("isFinish", "yes")

	if queryCondition.Pn == "None" || queryCondition.StartTime == "None" || queryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少参数"})
		return
	}
	qaWorkOrderIdList, err := ModuleQaStatisticDisplay.GetWorkOrderIds(&queryCondition, isFinish)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaWorkOrderIdList)
	}
}
