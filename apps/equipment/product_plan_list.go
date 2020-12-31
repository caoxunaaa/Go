package equipment

import (
	"SuperxonWebSite/Models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProjectPlanList(c *gin.Context) {
	projectPlanInfoList, err := Models.GetProjectPlanList()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, projectPlanInfoList)
	}
}
