package runOsaDisplayBroad

import (
	"SuperxonWebSite/Models/OsaRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOsaListHandler(c *gin.Context) {
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	osaQueryCondition.StartTime, osaQueryCondition.EndTime = Utils.GetAgoAndCurrentTime(Utils.Ago{Years: 0, Months: -1, Days: 0})
	osaList, err := OsaRunDisplay.GetOsaList(&osaQueryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, osaList)
	}
}

func GetOsaInfoListHandler(c *gin.Context) {
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	osaQueryCondition.StartTime, osaQueryCondition.EndTime = Utils.GetCurrentAndZeroTime()
	var ok bool
	osaQueryCondition.Pn, ok = c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	osaInfoList, err := OsaRunDisplay.GetOsaInfoList(&osaQueryCondition)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, osaInfoList)
	}
}

func GetAllOsaInfoListHandler(c *gin.Context) {
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	osaQueryCondition.StartTime, osaQueryCondition.EndTime = Utils.GetCurrentAndZeroTime()
	osaInfoList, err := OsaRunDisplay.GetAllOsaInfoList(&osaQueryCondition)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, osaInfoList)
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
