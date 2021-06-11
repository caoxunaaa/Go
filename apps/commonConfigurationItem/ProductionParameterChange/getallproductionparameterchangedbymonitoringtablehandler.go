package ProductionParameterChange

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllProductionParameterChangedByMonitoringTableHandler(c *gin.Context) {
	mt := c.Query("monitoring_table")
	res, err := ProductionLineOracleRelation.FindAllProductionParameterChangedByMonitoringTable(mt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
