package commonConfigurationItem

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllSettingWarningThresholdHandler(c *gin.Context) {
	res, err := ProductionLineOracleRelation.FindAllSettingWarningThreshold()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
