package osaStatistic

import (
	"SuperxonWebSite/Models/OsaQaStatisticDisplay"
	"SuperxonWebSite/Models/OsaRunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOsaYieldInfoByWorkOrderIdHandler(c *gin.Context) {
	var err error
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	osaQueryCondition.WorkOrderId = c.DefaultQuery("workOrderId", "None")
	if osaQueryCondition.WorkOrderId == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaOsaStatisticInfoByWorkOrderIdList, err := OsaQaStatisticDisplay.GetQaOsaStatisticInfoListByWorkOrderId(&osaQueryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaOsaStatisticInfoByWorkOrderIdList)
	}
}
