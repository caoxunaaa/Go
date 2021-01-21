package qaStatisticBroad

import (
	"SuperxonWebSite/Models/QaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetQaPnListHandler(c *gin.Context) {
	var queryCondition QaStatisticDisplay.QueryCondition
	queryCondition.StartTime = c.DefaultQuery("startTime", "None")
	queryCondition.EndTime = c.DefaultQuery("endTime", "None")
	if queryCondition.StartTime == "None" || queryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}

	qaPnList, err := QaStatisticDisplay.GetQaPnList(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaPnList)
	}
}

func GetQaStatisticInfoListHandler(c *gin.Context) {
	var err error
	var queryCondition QaStatisticDisplay.QueryCondition
	var qaStatisticInfoList []QaStatisticDisplay.QaStatisticInfo
	queryCondition.Pn = c.DefaultQuery("pn", "None")
	queryCondition.StartTime = c.DefaultQuery("startTime", "None")
	queryCondition.EndTime = c.DefaultQuery("endTime", "None")
	queryCondition.WorkOrderType = c.DefaultQuery("workOrderType", "None")
	if queryCondition.Pn == "None" || queryCondition.StartTime == "None" || queryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	if queryCondition.WorkOrderType == "TRX正常品" || queryCondition.WorkOrderType == "TRX改制返工品" {
		qaStatisticInfoList, err = QaStatisticDisplay.GetQaStatisticOrderInfoList(&queryCondition)
	} else {
		qaStatisticInfoList, err = QaStatisticDisplay.GetQaStatisticInfoList(&queryCondition)
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaStatisticInfoList)
	}
}

func GetQaDefectsInfoListByPnHandler(c *gin.Context) {
	var err error
	var queryCondition QaStatisticDisplay.QueryCondition
	var qaDefectsInfoList []QaStatisticDisplay.QaDefectsInfoByPn
	queryCondition.Pn = c.DefaultQuery("pn", "None")
	queryCondition.StartTime = c.DefaultQuery("startTime", "None")
	queryCondition.EndTime = c.DefaultQuery("endTime", "None")
	queryCondition.WorkOrderType = c.DefaultQuery("workOrderType", "None")
	if queryCondition.Pn == "None" || queryCondition.StartTime == "None" || queryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	if queryCondition.WorkOrderType == "TRX正常品" || queryCondition.WorkOrderType == "TRX改制返工品" {
		qaDefectsInfoList, err = QaStatisticDisplay.GetQaDefectsOrderInfoListByPn(&queryCondition)
	} else {
		qaDefectsInfoList, err = QaStatisticDisplay.GetQaDefectsInfoListByPn(&queryCondition)
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaDefectsInfoList)
	}
}

func GetQaDefectsDetailByPnHandler(c *gin.Context) {
	var err error
	var queryCondition QaStatisticDisplay.QueryCondition
	var qaDefectsDetailInfoList []QaStatisticDisplay.QaDefectsDetailInfo
	queryCondition.Pn = c.DefaultQuery("pn", "None")
	queryCondition.StartTime = c.DefaultQuery("startTime", "None")
	queryCondition.EndTime = c.DefaultQuery("endTime", "None")
	queryCondition.Process = c.DefaultQuery("process", "None")
	queryCondition.ErrorCode = c.DefaultQuery("errorCode", "None")
	queryCondition.WorkOrderType = c.DefaultQuery("workOrderType", "")
	if queryCondition.Pn == "None" || queryCondition.StartTime == "None" || queryCondition.Process == "None" || queryCondition.EndTime == "None" || queryCondition.WorkOrderType == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaDefectsDetailInfoList, err = QaStatisticDisplay.GetQaDefectsDetailByPn(&queryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaDefectsDetailInfoList)
	}
}
