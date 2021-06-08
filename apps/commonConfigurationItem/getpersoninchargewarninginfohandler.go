package commonConfigurationItem

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllPersonInChargeWarningInfoHandler(c *gin.Context) {
	res, err := ProductionLineOracleRelation.FindAllPersonInChargeWarningInfo()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
