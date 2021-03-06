package moduleStatistic

import (
	"SuperxonWebSite/Models/ModuleStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleSettingParamHandler(c *gin.Context) {
	var queryCondition ModuleStatisticDisplay.QueryCondition
	queryCondition.Pn = c.DefaultQuery("pn", "")
	queryCondition.WorkOrderId = c.DefaultQuery("workOrderId", "")
	queryCondition.BomId = c.DefaultQuery("bomId", "")
	queryCondition.Process = c.DefaultQuery("process", "")

	resultList, err := ModuleStatisticDisplay.GetPnSetParams(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resultList)
	}
}
