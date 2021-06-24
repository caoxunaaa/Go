package osaStatistic

import (
	"SuperxonWebSite/Models/OsaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOsaOpticsDefectInfoOfTxInChartByOsaPnAndErrorCodeHandler(c *gin.Context) {
	osaPn := c.DefaultQuery("osaPn", "None")
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")
	errorCode := c.DefaultQuery("errorCode", "None")

	if startTime == "None" || endTime == "None" || osaPn == "None" || errorCode == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	res, err := OsaStatisticDisplay.RedisGetOsaOpticsDefectInfoOfTxByPn(osaPn, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	resDetail := make(map[string]interface{})
	//统计某个错误码的散点分布图,分1G和10G
	resDetail["1G-结果值散点分布图"], resDetail["10G-结果值散点分布图"], err = OsaStatisticDisplay.GetErrorCodeDetailByErrorCode(res, errorCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	//统计某个错误码在各测试工位上的分布
	resDetail["测试工位错误码分布图"], err = OsaStatisticDisplay.GetStationIdDistributionByErrorCode(res, errorCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	//统计某个错误码在各耦合工位上的分布
	resDetail["耦合工位错误码分布图"], err = OsaStatisticDisplay.GetInsNameDistributionByErrorCode(res, errorCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resDetail)
}
