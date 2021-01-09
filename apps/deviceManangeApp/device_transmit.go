package deviceManangeApp

import (
	"SuperxonWebSite/Models/DeviceManage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllDeviceTransmitInfoList(c *gin.Context) {
	deviceTransmitInfoList, err := DeviceManage.GetAllDeviceTransmitInfoList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, deviceTransmitInfoList)
	}
}
