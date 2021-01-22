package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//保养计划表
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

func CreateDeviceMaintenanceItemHandler(c *gin.Context) {
	var deviceMaintenanceItem DeviceManage.DeviceMaintenanceItem
	if err := c.ShouldBindJSON(&deviceMaintenanceItem); err == nil {
		err = DeviceManage.CreateDeviceMaintenanceItem(&deviceMaintenanceItem)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"MaintenanceName": deviceMaintenanceItem.Name,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func UpdateDeviceMaintenanceItemHandler(c *gin.Context) {
	var deviceMaintenanceItem DeviceManage.DeviceMaintenanceItem
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	if err := c.ShouldBindJSON(&deviceMaintenanceItem); err == nil {
		oldId, _ := strconv.Atoi(id)
		err = DeviceManage.UpdateDeviceMaintenanceItem(&deviceMaintenanceItem, uint(oldId))
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"MaintenanceName": deviceMaintenanceItem.Name,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func DeleteDeviceMaintenanceItemHandler(c *gin.Context) {
	var deviceMaintenanceItem DeviceManage.DeviceMaintenanceItem
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	if err := c.ShouldBindJSON(&deviceMaintenanceItem); err == nil {
		oldId, _ := strconv.Atoi(id)
		err = DeviceManage.DeleteDeviceMaintenanceItem(&deviceMaintenanceItem, uint(oldId))
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"MaintenanceName": deviceMaintenanceItem.Name,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

//绑定保养项目
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

//解绑保养项目
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

//保养当前信息
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

//func UpdateDeviceMaintenanceCurrentInfoHandler(c *gin.Context) {
//	var deviceMaintenanceCurrentInfoList DeviceManage.DeviceMaintenanceCurrentInfo
//	id, ok := c.Params.Get("id")
//	if !ok {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
//		return
//	}
//	if err := c.ShouldBindJSON(&deviceMaintenanceCurrentInfoList); err == nil {
//		oldId, _ := strconv.Atoi(id)
//		_, err = DeviceManage.UpdateDeviceMaintenanceCurrentInfo(&deviceMaintenanceCurrentInfoList, uint(oldId), false)
//		if err == nil {
//			c.JSON(http.StatusOK, gin.H{
//				"MaintenanceName": deviceMaintenanceCurrentInfoList.DeviceName,
//			})
//		} else {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		}
//	} else {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//	}
//}

//保养记录
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

func CreateDeviceMaintenanceRecordHandler(c *gin.Context) {
	var deviceMaintenanceRecord *DeviceManage.DeviceMaintenanceRecord
	deviceMaintenanceRecord = new(DeviceManage.DeviceMaintenanceRecord)

	deviceMaintenanceRecord.DeviceName = c.PostForm("DeviceName")
	deviceMaintenanceRecord.DeviceSn = c.PostForm("DeviceSn")
	deviceMaintenanceRecord.DeviceAssets = c.PostForm("DeviceAssets")
	deviceMaintenanceRecord.DeviceSort = c.PostForm("DeviceSort")
	deviceMaintenanceRecord.ItemCategory = c.PostForm("ItemCategory")
	deviceMaintenanceRecord.ItemName = c.PostForm("ItemName")
	deviceMaintenanceRecord.MaintenanceTime, _ = time.ParseInLocation("2006-01-02 15:04:05", c.PostForm("MaintenanceTime"), time.Local)
	deviceMaintenanceRecord.MaintenanceUser.String = c.PostForm("MaintenanceUser")
	deviceMaintenanceRecord.Remark.String = c.PostForm("Remark")
	maintenanceRecordFile, err := c.FormFile("maintenanceRecordFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fileName := maintenanceRecordFile.Filename
	fmt.Println("fileName", fileName)

	nameSplit := strings.Split(fileName, ".")
	dir := nameSplit[len(nameSplit)-1]

	deviceMaintenanceRecord.FilePath.String = "./assets/" + dir + "/" + fileName

	_, err = os.Stat("./assets")
	if os.IsNotExist(err) {
		fmt.Println("目录不存在,创建目录")
		err = os.Mkdir("./assets", 0777)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	_, err = os.Stat("./assets/" + dir)
	if os.IsNotExist(err) {
		fmt.Println("文件不存在,创建目录")
		err = os.Mkdir("./assets/"+dir, 0777)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	err = c.SaveUploadedFile(maintenanceRecordFile, "./assets/"+dir+"/"+fileName)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = DeviceManage.CreateDeviceMaintenanceRecord(deviceMaintenanceRecord)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": fileName + "已经成功上传"})
	}
}
