package moduleStatistic

import (
	"SuperxonWebSite/Models/ModuleQaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleDefectsInfoByWorkOrderIdHandler(c *gin.Context) {
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	queryCondition.WorkOrderId = c.DefaultQuery("workOrderId", "None")
	if queryCondition.WorkOrderId == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}

	workOrderYieldList, err := ModuleQaStatisticDisplay.GetQaDefectsInfoByWorkOrderIdList(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, workOrderYieldList)
	}
}
