package ProductionParameterChange

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProductionParameterChangedRelationByTableNameHandler(c *gin.Context) {
	tn := c.Query("table_name")
	res, err := ProductionLineOracleRelation.FindProductionParameterChangedRelationByTableName(tn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
