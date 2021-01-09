package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllDeviceBaseInfoList(c *gin.Context) {
	deviceBaseInfoList, err := DeviceManage.GetAllDeviceBaseInfoList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceBaseInfoList)
	}
}

func GetDeviceBaseInfo(c *gin.Context) {
	snAssetsIc, ok := c.Params.Get("snAssetsIc")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn Assets 或者InternalCoding"})
		return
	}
	deviceBaseInfo, err := DeviceManage.GetDeviceBaseInfo(snAssetsIc)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceBaseInfo)
	}
}
