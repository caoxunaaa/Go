package osaRunning

import (
	"SuperxonWebSite/Models/OsaRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOsaAllInfoTodayHandler(c *gin.Context) {
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	osaQueryCondition.StartTime, osaQueryCondition.EndTime = Utils.GetCurrentAndZeroTime()
	osaInfoList, err := OsaRunDisplay.GetAllOsaInfoList(&osaQueryCondition)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, osaInfoList)
	}
}
