package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func CreateCreateDeviceRepairInfoHandler(c *gin.Context) {
	var deviceRepairInfo DeviceManage.DeviceRepairInfo
	repairStatus := c.DefaultQuery("repairStatus", "")
	if err := c.ShouldBindJSON(&deviceRepairInfo); err == nil {
		err = DeviceManage.CreateDeviceRepairInfo(&deviceRepairInfo, repairStatus)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"DeviceName": deviceRepairInfo.Name,
				"DeviceSn":   deviceRepairInfo.Sn,
				"Status":     repairStatus,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func UpdateDeviceRepairInfoHandler(c *gin.Context) {
	var deviceRepairInfo DeviceManage.DeviceRepairInfo
	repairStatus := c.DefaultQuery("repairStatus", "")
	oldId, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn"})
		return
	}
	if err := c.ShouldBindJSON(&deviceRepairInfo); err == nil {
		id, _ := strconv.Atoi(oldId)
		_, err = DeviceManage.UpdateDeviceRepairInfo(&deviceRepairInfo, uint(id), repairStatus)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"DeviceName": deviceRepairInfo.Name,
				"DeviceSn":   deviceRepairInfo.Sn,
				"Status":     repairStatus,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
