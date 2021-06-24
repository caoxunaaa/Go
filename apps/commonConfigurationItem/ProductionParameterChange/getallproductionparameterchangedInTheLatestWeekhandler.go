package ProductionParameterChange

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllProductionParameterChangedInTheLatestWeekHandler(c *gin.Context) {
	res, err := ProductionLineOracleRelation.FindAllProductionParameterChangedInTheLatestWeek()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
