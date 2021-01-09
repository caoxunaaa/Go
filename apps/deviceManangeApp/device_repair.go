package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllDeviceRepairInfoListHandler(c *gin.Context) {
	deviceRepairInfoList, err := DeviceManage.GetAllDeviceRepairInfoList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceRepairInfoList)
	}
}

func GetDeviceRepairInfoHandler(c *gin.Context) {
	deviceSn, ok := c.Params.Get("deviceSn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn"})
	}
	deviceRepairInfo, err := DeviceManage.GetDeviceRepairInfo(deviceSn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceRepairInfo)
	}
}
