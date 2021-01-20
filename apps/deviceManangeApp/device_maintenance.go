package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllDeviceMaintenanceCategoryListHandler(c *gin.Context) {
	deviceMaintenanceItemCategoryList, err := DeviceManage.GetAllDeviceMaintenanceCategoryList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceItemCategoryList)
	}
}

func GetDeviceMaintenanceItemOfCategoryHandler(c *gin.Context) {
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

func GetAllDeviceMaintenanceCurrentInfoListHandler(c *gin.Context) {
	deviceMaintenanceCurrentInfoList, err := DeviceManage.GetAllDeviceMaintenanceCurrentInfoList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceCurrentInfoList)
	}
}

func GetDeviceMaintenanceCurrentInfoHandler(c *gin.Context) {
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

func GetAllDeviceMaintenanceAllRecordsHandler(c *gin.Context) {
	deviceMaintenanceRecords, err := DeviceManage.GetAllDeviceMaintenanceRecords("")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceMaintenanceRecords)
	}
}

func GetAllDeviceMaintenanceRecordsOfItemNameHandler(c *gin.Context) {
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

func GetDeviceMaintenanceRecordsHandler(c *gin.Context) {
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

func GetDeviceMaintenanceRecordOfItemNameHandler(c *gin.Context) {
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

func BindDeviceMaintenanceItemHandler(c *gin.Context) {
	var deviceMaintenanceItems []*DeviceManage.DeviceMaintenanceItem
	deviceSn, ok := c.Params.Get("deviceSn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn"})
		return
	}
	if err := c.ShouldBindJSON(&deviceMaintenanceItems); err == nil {
		err = DeviceManage.BindDeviceMaintenanceItem(deviceSn, deviceMaintenanceItems)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"DeviceSn":         deviceSn,
				"MaintenanceItems": deviceMaintenanceItems,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func UnBindDeviceMaintenanceItemHandler(c *gin.Context) {
	deviceSn, ok := c.Params.Get("deviceSn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn"})
		return
	}
	err := DeviceManage.UnBindDeviceMaintenanceItem(deviceSn)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"DeviceSn": deviceSn,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
