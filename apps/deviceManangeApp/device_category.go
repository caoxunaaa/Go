package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllDeviceCategoryRootList(c *gin.Context) {
	deviceCategoryRootNameList, err := DeviceManage.GetAllDeviceCategoryRootList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceCategoryRootNameList)
	}
}

func GetAllDeviceCategoryChildList(c *gin.Context) {
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

func CreateDeviceCategoryChild(c *gin.Context) {
	var deviceCategory DeviceManage.DeviceCategory
	if err := c.ShouldBind(&deviceCategory); err == nil {
		err = DeviceManage.CreateDeviceCategoryChild(deviceCategory.ParentCategoryName.String, deviceCategory.Name)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"ChildName": deviceCategory.Name,
				"RootName":  deviceCategory.ParentCategoryName,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
