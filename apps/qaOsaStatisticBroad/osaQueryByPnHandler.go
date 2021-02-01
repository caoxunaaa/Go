package qaOsaStatisticBroad

import (
	"SuperxonWebSite/Models/OsaQaStatisticDisplay"
	"SuperxonWebSite/Models/OsaRunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetQaOsaPnListHandler(c *gin.Context) {
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	osaQueryCondition.StartTime = c.DefaultQuery("startTime", "None")
	osaQueryCondition.EndTime = c.DefaultQuery("endTime", "None")

	if osaQueryCondition.StartTime == "None" || osaQueryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaOsaPnList, err := OsaQaStatisticDisplay.GetQaOsaPnList(&osaQueryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaOsaPnList)
	}
}

func GetQaOsaStatisticInfoListHandler(c *gin.Context) {
	var err error
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	osaQueryCondition.Pn = c.DefaultQuery("pn", "None")
	osaQueryCondition.StartTime = c.DefaultQuery("startTime", "None")
	osaQueryCondition.EndTime = c.DefaultQuery("endTime", "None")
	if osaQueryCondition.Pn == "None" || osaQueryCondition.StartTime == "None" || osaQueryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaOsaStatisticInfoList, err := OsaQaStatisticDisplay.RedisGetQaOsaStatisticInfoList(&osaQueryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaOsaStatisticInfoList)
	}
}

func GetQaOsaDefectsInfoListByPnHandler(c *gin.Context) {
	var err error
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	osaQueryCondition.Pn = c.DefaultQuery("pn", "None")
	osaQueryCondition.StartTime = c.DefaultQuery("startTime", "None")
	osaQueryCondition.EndTime = c.DefaultQuery("endTime", "None")
	if osaQueryCondition.Pn == "None" || osaQueryCondition.StartTime == "None" || osaQueryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaOsaDefectsInfoList, err := OsaQaStatisticDisplay.RedisGetQaOsaDefectsInfoListByPn(&osaQueryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaOsaDefectsInfoList)
	}
}
