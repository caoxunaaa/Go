// @Title  device_category.go
// @Description  设备类型的查询和添加app
// @Author  曹迅 (时间 2021/01/01  12:00)
// @Update  曹迅 (时间 2021/02/03  12:00)
package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
)

//获取设备类型一级目录
func GetAllDeviceCategoryRootListHandler(c *gin.Context) {
	deviceCategoryRootNameList, err := DeviceManage.GetAllDeviceCategoryRootList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceCategoryRootNameList)
	}
}

//获取一级设备类型下的所有二级设备类型目录
func GetAllDeviceCategoryChildListHandler(c *gin.Context) {
	rootCategory, ok := c.Params.Get("rootCategory")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	deviceCategoryChildNameList, err := DeviceManage.GetAllDeviceCategoryChildList(rootCategory)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceCategoryChildNameList)
	}
}

//创建一个设备类型
func CreateDeviceCategoryChildHandler(c *gin.Context) {
	var deviceCategory DeviceManage.DeviceCategory
	if err := c.BindJSON(&deviceCategory); err == nil {
		err = DeviceManage.CreateDeviceCategoryChild(&deviceCategory)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"ChildName": deviceCategory.Name,
				"RootName":  deviceCategory.ParentCategoryName.String,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
