package qaModuleStatisticBroad

import (
	"SuperxonWebSite/Models/ModuleQaStatisticDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPnSetParamsListHandler(c *gin.Context) {
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	queryCondition.Pn = c.DefaultQuery("pn", "")
	queryCondition.WorkOrderId = c.DefaultQuery("workOrderId", "")
	queryCondition.BomId = c.DefaultQuery("bomId", "")
	queryCondition.Process = c.DefaultQuery("process", "")

	resultList, err := ModuleQaStatisticDisplay.GetPnSetParams(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resultList)
	}
}

func GetPnWorkOrderYieldsListHandler(c *gin.Context) {
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	StartTime, EndTime := Utils.GetCurrentAndZeroTime()
	queryCondition.Pn = c.DefaultQuery("pn", "")
	queryCondition.WorkOrderType = c.DefaultQuery("workOrderType", "")
	queryCondition.StartTime = c.DefaultQuery("startTime", StartTime)
	queryCondition.EndTime = c.DefaultQuery("endTime", EndTime)
	resultList, err := ModuleQaStatisticDisplay.GetWorkOrderYieldsByPn(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, resultList)
	}
}

func GetWorkOrderIdsHandler(c *gin.Context) {
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

func GetWorkOrderYieldsByWorkOrderIdListHandler(c *gin.Context) {
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	queryCondition.WorkOrderId = c.DefaultQuery("workOrderId", "None")
	if queryCondition.WorkOrderId == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}

	workOrderYieldList, err := ModuleQaStatisticDisplay.GetWorkOrderYieldsByWorkOrderId(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, workOrderYieldList)
	}
}

func GetQaDefectsInfoByWorkOrderIdListHandler(c *gin.Context) {
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

func GetQaDefectsDetailByWorkOrderIdHandler(c *gin.Context) {
	var err error
	var queryCondition ModuleQaStatisticDisplay.QueryCondition
	var qaDefectsDetailInfoList []ModuleQaStatisticDisplay.QaDefectsDetailInfo
	queryCondition.WorkOrderId = c.DefaultQuery("workOrderId", "None")
	queryCondition.Process = c.DefaultQuery("process", "None")
	queryCondition.ErrorCode = c.DefaultQuery("errorCode", "None")
	if queryCondition.WorkOrderId == "None" || queryCondition.Process == "None" || queryCondition.ErrorCode == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少参数"})
		return
	}
	qaDefectsDetailInfoList, err = ModuleQaStatisticDisplay.GetQaDefectsDetailByWorkOrderId(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaDefectsDetailInfoList)
	}
}
