package osaStatistic

import (
	"SuperxonWebSite/Models/OsaStatisticDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOsaOpticsDefectInfoOfTxByOsaPnHandler(c *gin.Context) {
	osaPn := c.DefaultQuery("osaPn", "None")
	startTime := c.DefaultQuery("startTime", "None")
	endTime := c.DefaultQuery("endTime", "None")

	if startTime == "None" || endTime == "None" || osaPn == "None" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数"})
		return
	}
	var err error
	res := make(map[string]interface{})
	r, err := OsaStatisticDisplay.RedisGetOsaOpticsDefectInfoOfTxByPn(osaPn, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	res["完整数据"] = r
	res["所有错误代码分布"], err = OsaStatisticDisplay.GetErrorCodeDistribution(r)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
