package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllDeviceTransmitInfoListHandler(c *gin.Context) {
	deviceTransmitInfoList, err := DeviceManage.GetAllDeviceTransmitInfoList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceTransmitInfoList)
	}
}

func GetDeviceTransmitInfoHandler(c *gin.Context) {
	deviceSn, ok := c.Params.Get("deviceSn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn"})
	}
	deviceTransmitInfo, err := DeviceManage.GetDeviceTransmitInfo(deviceSn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceTransmitInfo)
	}
}

func CreateDeviceTransmitInfoHandler(c *gin.Context) {
	var deviceTransmitInfo DeviceManage.DeviceTransmitInfo
	if err := c.ShouldBindJSON(&deviceTransmitInfo); err == nil {
		err = DeviceManage.CreateDeviceTransmitInfo(&deviceTransmitInfo)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"DeviceName": deviceTransmitInfo.DeviceName,
				"DeviceSn":   deviceTransmitInfo.DeviceSn,
				"OldOwner":   deviceTransmitInfo.OldOwner,
				"NewOwner":   deviceTransmitInfo.NewOwner,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func DeleteDeviceTransmitInfoHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
	}
	idInt, _ := strconv.Atoi(id)
	length, err := DeviceManage.DeleteDeviceTransmitInfo(idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": strconv.FormatInt(length, 10) + "行已经被删除"})
	}
}
