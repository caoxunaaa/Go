package moduleRunning

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleAllInfoTodayHandler(c *gin.Context) {
	startTimeStr, endTimeStr := Utils.GetCurrentAndZeroTime()
	moduleInfoList, err := ModuleRunDisplay.GetAllModuleInfoList(startTimeStr, endTimeStr)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfoList)
	}
}
