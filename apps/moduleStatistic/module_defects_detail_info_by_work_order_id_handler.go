package moduleStatistic

import (
	"SuperxonWebSite/Models/ModuleStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetQaDefectsDetailByWorkOrderIdHandler(c *gin.Context) {
	var err error
	var queryCondition ModuleStatisticDisplay.QueryCondition
	var qaDefectsDetailInfoList []ModuleStatisticDisplay.QaDefectsDetailInfo
	queryCondition.WorkOrderId = c.DefaultQuery("workOrderId", "None")
	queryCondition.Process = c.DefaultQuery("process", "None")
	queryCondition.ErrorCode = c.DefaultQuery("errorCode", "None")
	if queryCondition.WorkOrderId == "None" || queryCondition.Process == "None" || queryCondition.ErrorCode == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少参数"})
		return
	}
	qaDefectsDetailInfoList, err = ModuleStatisticDisplay.GetQaDefectsDetailByWorkOrderId(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaDefectsDetailInfoList)
	}
}
