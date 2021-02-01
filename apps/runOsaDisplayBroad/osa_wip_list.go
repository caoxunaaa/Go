package runOsaDisplayBroad

import (
	"SuperxonWebSite/Models/OsaRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOsaWipInfoListHandler(c *gin.Context) {
	var osaQueryCondition OsaRunDisplay.OsaQueryCondition
	var ok bool
	osaQueryCondition.Pn, ok = c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的PN"})
	}

	osaQueryCondition.StartTime, osaQueryCondition.EndTime = Utils.GetCurrentAndZeroTime()
	wipInfoList, err := OsaRunDisplay.GetOsaWipInfoList(&osaQueryCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, wipInfoList)
	}
}
