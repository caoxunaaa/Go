package waringDisplayBroad

import (
	"SuperxonWebSite/Models/WaringDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWarningCountChartDataListHandler(c *gin.Context) {
	var chartQueryCondition WaringDisplay.ChartQueryCondition
	chartQueryCondition.StartTime = c.DefaultQuery("startTime", "None")
	chartQueryCondition.EndTime = c.DefaultQuery("endTime", "None")
	chartQueryCondition.Classification = c.DefaultQuery("classification", "None")
	if chartQueryCondition.Classification == "None" || chartQueryCondition.StartTime == "None" || chartQueryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	warningCountChartData, err := WaringDisplay.GetWarningCountChartDataList(&chartQueryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, warningCountChartData)
	}

}

func GetWarningToPersonListHandler(c *gin.Context) {
	personToPnManage, err := Utils.GetWarningToPersonList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, personToPnManage)
	}
}
