package equipment

import (
	"SuperxonWebSite/Models"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProductList(c *gin.Context) {
	startTimeStr, endTimeStr := Utils.GetTwoMonthAgoAndCurrentTime()
	moduleList, err := Models.GetModuleList(startTimeStr, endTimeStr)
	osaList, err := Models.GetOsaList(startTimeStr, endTimeStr)
	var product []Models.Product
	if osaList != nil {
		product = append(product, osaList...)
	}
	if moduleList != nil {
		product = append(product, moduleList...)
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
