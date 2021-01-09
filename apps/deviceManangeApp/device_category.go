package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllDeviceCategoryRootListHandler(c *gin.Context) {
	deviceCategoryRootNameList, err := DeviceManage.GetAllDeviceCategoryRootList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceCategoryRootNameList)
	}
}

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
