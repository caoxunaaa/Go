package moduleRunning

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"SuperxonWebSite/Utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetModuleAllUsingPnHandler(c *gin.Context) {
	startTime, endTime := Utils.GetAgoAndCurrentTime(Utils.Ago{Years: 0, Months: -1, Days: 0})
	moduleList, err := ModuleRunDisplay.GetModuleList(startTime, endTime)
	var product []ModuleRunDisplay.Product
	if moduleList != nil {
		product = append(product, moduleList...)
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, product)
	}
}
