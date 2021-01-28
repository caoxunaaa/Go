package runModuleDisplayBroad

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProjectPlanListHandler(c *gin.Context) {
	projectPlanInfoList, err := ModuleRunDisplay.RedisGetProjectPlanList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, projectPlanInfoList)
	}
}
