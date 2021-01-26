package runDisplayBroad

import (
	"SuperxonWebSite/Models/RunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProjectPlanListHandler(c *gin.Context) {
	projectPlanInfoList, err := RunDisplay.RedisGetProjectPlanList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, projectPlanInfoList)
	}
}
