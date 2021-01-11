package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllDeviceBaseInfoListHandler(c *gin.Context) {
	deviceBaseInfoList, err := DeviceManage.GetAllDeviceBaseInfoList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceBaseInfoList)
	}
}

func GetDeviceBaseInfoHandler(c *gin.Context) {
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

func CreateDeviceBaseInfoHandler(c *gin.Context) {
	var deviceBaseInfo DeviceManage.DeviceBaseInfo
	if err := c.ShouldBindJSON(&deviceBaseInfo); err == nil {
		err = DeviceManage.CreateDeviceBaseInfo(&deviceBaseInfo)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"DeviceName": deviceBaseInfo.Name,
				"DeviceSn":   deviceBaseInfo.Sn,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func DeleteDeviceBaseInfoHandler(c *gin.Context) {
	deviceSn, ok := c.Params.Get("deviceSn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn"})
		return
	}
	length, err := DeviceManage.DeleteDeviceBaseInfo(deviceSn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": strconv.FormatInt(length, 10) + "行已经被删除"})
	}
}

func UpdateDeviceBaseInfoHandler(c *gin.Context) {
	var deviceBaseInfo DeviceManage.DeviceBaseInfo
	deviceSn, ok := c.Params.Get("deviceSn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn"})
		return
	}
	err := c.ShouldBindJSON(&deviceBaseInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	length, err := DeviceManage.UpdateDeviceBaseInfo(&deviceBaseInfo, deviceSn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "有" + strconv.FormatInt(length, 10) + "行数据已经被更新"})
	}
}
