package ProductionParameterChange

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetParameterByMonitoringTableAndOnlyFieldAndChangedItemHandler(c *gin.Context) {
	var p ProductionLineOracleRelation.ProductionParameterChanged
	p.MonitoringTable = c.PostForm("monitoring_table")
	p.OnlyField1Name = c.PostForm("only_field_1_name")
	p.OnlyField1Value = c.PostForm("only_field_1_value")
	p.OnlyField2Name = c.PostForm("only_field_2_name")
	p.OnlyField2Value = c.PostForm("only_field_2_value")
	p.OnlyField3Name = c.PostForm("only_field_3_name")
	p.OnlyField3Value = c.PostForm("only_field_3_value")
	p.OnlyField4Name = c.PostForm("only_field_4_name")
	p.OnlyField4Value = c.PostForm("only_field_4_value")
	p.ChangedItem = c.PostForm("changed_item")
	res, err := ProductionLineOracleRelation.GetParameterByMonitoringTableAndOnlyFieldAndChangedItem(&p)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
