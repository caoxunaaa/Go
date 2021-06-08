package osaStatistic

import (
	"SuperxonWebSite/Models/OsaRunDisplay"
	"SuperxonWebSite/Models/OsaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOsaDefectsInfoByPnHandler(c *gin.Context) {
	var err error
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	osaQueryCondition.Pn = c.DefaultQuery("pn", "None")
	osaQueryCondition.StartTime = c.DefaultQuery("startTime", "None")
	osaQueryCondition.EndTime = c.DefaultQuery("endTime", "None")
	if osaQueryCondition.Pn == "None" || osaQueryCondition.StartTime == "None" || osaQueryCondition.EndTime == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	qaOsaDefectsInfoList, err := OsaStatisticDisplay.GetQaOsaDefectsInfoListByPn(&osaQueryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, qaOsaDefectsInfoList)
	}
}
