package runModuleDisplayBroad

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleListHandler(c *gin.Context) {
	startTime, endTime := Utils.GetAgoAndCurrentTime(Utils.Ago{Years: 0, Months: -1, Days: 0})
	moduleList, err := ModuleRunDisplay.GetModuleList(startTime, endTime)
	//fmt.Println(moduleList)
	var product []ModuleRunDisplay.Product
	if moduleList != nil {
		product = append(product, moduleList...)
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, product)
	}
}

func GetModuleInfoListHandler(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	startTimeStr, endTimeStr := Utils.GetCurrentAndZeroTime()
	moduleInfoList, err := ModuleRunDisplay.GetModuleInfoList(pn, startTimeStr, endTimeStr)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfoList)
	}
}

func GetAllModuleInfoListHandler(c *gin.Context) {
	startTimeStr, endTimeStr := Utils.GetCurrentAndZeroTime()
	moduleInfoList, err := ModuleRunDisplay.GetAllModuleInfoList(startTimeStr, endTimeStr)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfoList)
	}
}

func GetYesterdayModuleInfoListHandler(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	moduleInfo, err := ModuleRunDisplay.GetYesterdayModuleInfoList(pn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfo)
	}
}
