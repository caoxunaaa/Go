package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllDeviceMaintenanceCategoryList(c *gin.Context) {
	deviceMaintenanceItemCategoryList, err := DeviceManage.GetAllDeviceMaintenanceCategoryList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceItemCategoryList)
	}
}

func GetDeviceMaintenanceItemOfCategory(c *gin.Context) {
	category, ok := c.Params.Get("category")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的保养类型"})
		return
	}
	deviceMaintenanceItems, err := DeviceManage.GetDeviceMaintenanceItemOfCategory(category)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceItems)
	}
}

func GetAllDeviceMaintenanceCurrentInfoList(c *gin.Context) {
	deviceMaintenanceCurrentInfoList, err := DeviceManage.GetAllDeviceMaintenanceCurrentInfoList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceCurrentInfoList)
	}
}

func GetDeviceMaintenanceCurrentInfo(c *gin.Context) {
	snAssets, ok := c.Params.Get("snAssets")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的保养类型"})
		return
	}
	deviceMaintenanceCurrentInfo, err := DeviceManage.GetDeviceMaintenanceCurrentInfo(snAssets)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceCurrentInfo)
	}
}

func GetAllDeviceMaintenanceAllRecords(c *gin.Context) {
	deviceMaintenanceRecords, err := DeviceManage.GetAllDeviceMaintenanceRecords("")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceRecords)
	}
}

func GetAllDeviceMaintenanceRecordsOfItemName(c *gin.Context) {
	itemName, ok := c.Params.Get("itemName")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的保养类型"})
		return
	}
	deviceMaintenanceRecords, err := DeviceManage.GetAllDeviceMaintenanceRecords(itemName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceRecords)
	}
}

func GetDeviceMaintenanceRecords(c *gin.Context) {
	snAssets, ok := c.Params.Get("snAssets")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn 或者 固资号"})
		return
	}

	deviceMaintenanceCurrentInfo, err := DeviceManage.GetDeviceMaintenanceRecords(snAssets, "")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceCurrentInfo)
	}
}

func GetDeviceMaintenanceRecordOfItemName(c *gin.Context) {
	snAssets, ok := c.Params.Get("snAssets")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn 或者 固资号"})
		return
	}
	itemName, ok := c.Params.Get("itemName")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的保养类型"})
		return
	}

	deviceMaintenanceCurrentInfo, err := DeviceManage.GetDeviceMaintenanceRecords(snAssets, itemName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceCurrentInfo)
	}
}
