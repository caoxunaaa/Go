package ProductionParameterChange

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllProductionParameterChangedByMonitoringTableAndOnlyFieldAndChangedItemHandler(c *gin.Context) {
	mt := c.Query("monitoring_table")
	f1 := c.Query("only_field_1_value")
	f2 := c.Query("only_field_2_value")
	f3 := c.Query("only_field_3_value")
	f4 := c.Query("only_field_4_value")
	ci := c.Query("changed_item")
	res, err := ProductionLineOracleRelation.FindAllProductionParameterChangedByMonitoringTableAndOnlyFieldAndChangedItem(mt, f1, f2, f3, f4, ci)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
