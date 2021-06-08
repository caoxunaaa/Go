package osaRunning

import (
	"SuperxonWebSite/Models/OsaRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllOsaTxCoupleInfoListHandler(c *gin.Context) {
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	osaQueryCondition.StartTime, osaQueryCondition.EndTime = Utils.GetCurrentAndZeroTime()
	osaTxCoupleInfoList, err := OsaRunDisplay.GetAllOsaTxCoupleInfoList(&osaQueryCondition)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, osaTxCoupleInfoList)
	}
}

func GetYesterdayOsaInfoListHandler(c *gin.Context) {
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	var ok bool
	osaQueryCondition.Pn, ok = c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	moduleInfo, err := OsaRunDisplay.GetYesterdayOsaInfoList(&osaQueryCondition)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfo)
	}
}
