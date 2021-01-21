package runDisplayBroad

import (
	"SuperxonWebSite/Models/RunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleListHandler(c *gin.Context) {
	startTime, endTime := Utils.GetCurrentAndZeroTime()
	startTime = c.DefaultQuery("startTime", startTime)
	endTime = c.DefaultQuery("endTime", endTime)
	moduleList, err := RunDisplay.GetModuleList(startTime, endTime)
	//fmt.Println(moduleList)
	var product []RunDisplay.Product
	if moduleList != nil {
		product = append(product, moduleList...)
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, product)
	}
}

func GetOsaListHandler(c *gin.Context) {
	startTime, endTime := Utils.GetCurrentAndZeroTime()
	startTime = c.DefaultQuery("startTime", startTime)
	endTime = c.DefaultQuery("endTime", endTime)
	osaList, err := RunDisplay.GetOsaList(startTime, endTime)
	var product []RunDisplay.Product
	if osaList != nil {
		product = append(product, osaList...)
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
	moduleInfo, err := RunDisplay.GetModuleInfoList(pn, startTimeStr, endTimeStr)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfo)
	}
}

func GetOsaInfoListHandler(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	startTimeStr, endTimeStr := Utils.GetCurrentAndZeroTime()
	moduleInfo, err := RunDisplay.GetOsaInfoList(pn, startTimeStr, endTimeStr)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfo)
	}
}

func GetYesterdayModuleInfoListHandler(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	moduleInfo, err := RunDisplay.GetYesterdayModuleInfoList(pn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfo)
	}
}

func GetYesterdayOsaInfoListHandler(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	moduleInfo, err := RunDisplay.GetYesterdayOsaInfoList(pn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfo)
	}
}
