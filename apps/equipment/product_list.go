package equipment

import (
	"SuperxonWebSite/Models"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleList(c *gin.Context) {
	startTimeStr, endTimeStr := Utils.GetAgoAndCurrentTime(Utils.Ago{Days: -10})
	moduleList, err := Models.GetModuleList(startTimeStr, endTimeStr)
	//fmt.Println(moduleList)
	var product []Models.Product
	if moduleList != nil {
		product = append(product, moduleList...)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, product)
	}
}

func GetOsaList(c *gin.Context) {
	startTimeStr, endTimeStr := Utils.GetAgoAndCurrentTime(Utils.Ago{Days: -10})
	osaList, err := Models.GetOsaList(startTimeStr, endTimeStr)

	var product []Models.Product
	if osaList != nil {
		product = append(product, osaList...)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, product)
	}
}

func GetModuleInfoList(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	startTimeStr, endTimeStr := Utils.GetCurrentAndZeroTime()
	moduleInfo, err := Models.GetModuleInfoList(pn, startTimeStr, endTimeStr)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfo)
	}
}

func GetOsaInfoList(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	startTimeStr, endTimeStr := Utils.GetCurrentAndZeroTime()
	moduleInfo, err := Models.GetOsaInfoList(pn, startTimeStr, endTimeStr)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfo)
	}
}

func GetYesterdayModuleInfoList(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	moduleInfo, err := Models.GetYesterdayModuleInfoList(pn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfo)
	}
}

func GetYesterdayOsaInfoList(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	moduleInfo, err := Models.GetYesterdayOsaInfoList(pn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, moduleInfo)
	}
}
