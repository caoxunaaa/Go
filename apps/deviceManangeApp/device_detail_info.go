// @Title  device_detail_info.go
// @Description  设备基本信息的增删改查app
// @Author  曹迅 (时间 2021/01/01  12:00)
// @Update  曹迅 (时间 2021/02/03  12:00)
package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//获取所有设备基本信息
func GetAllDeviceBaseInfoListHandler(c *gin.Context) {
	deviceBaseInfoList, err := DeviceManage.GetAllDeviceBaseInfoList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceBaseInfoList)
	}
}

//获取某个sn对应的设备基本信息
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

//创建一个设备基本信息
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

//删除一个设备基本信息
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

//更新一个设备基本信息
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
