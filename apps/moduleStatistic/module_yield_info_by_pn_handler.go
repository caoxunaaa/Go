package moduleStatistic

import (
	"SuperxonWebSite/Models/ModuleQaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleYieldInfoByPnHandler(c *gin.Context) {
	var err error
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	var qaStatisticInfoList []ModuleQaStatisticDisplay.QaStatisticInfo
	queryCondition.Pn = c.DefaultQuery("pn", "None")
	queryCondition.StartTime = c.DefaultQuery("startTime", "None")
	queryCondition.EndTime = c.DefaultQuery("endTime", "None")
	queryCondition.WorkOrderType = c.DefaultQuery("workOrderType", "")
	if queryCondition.Pn == "None" || queryCondition.StartTime == "None" || queryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaStatisticInfoList, err = ModuleQaStatisticDisplay.GetQaStatisticOrderInfoList(&queryCondition)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaStatisticInfoList)
	}
}
