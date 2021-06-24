package ProductionParameterChange

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"SuperxonWebSite/Services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func UpdateProductionParameterChangedHandler(c *gin.Context) {
	var p ProductionLineOracleRelation.ProductionParameterChanged
	var err error
	p.Username = c.PostForm("username")
	p.Nickname = c.PostForm("nickname")
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
	p.OldValue = c.PostForm("old_value")
	p.NewValue = c.PostForm("new_value")
	time1 := c.PostForm("update_time")
	fmt.Println("time1", time1)
	p.UpdateTime, err = time.ParseInLocation("2006-01-02 15:04:05", time1, time.Local)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	p.UpdateReason = c.PostForm("update_reason")

	err = ProductionLineOracleRelation.UpdateParameterByMonitoringTableAndField(&p)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		mes := "<html><b>" + p.Nickname + "</b> 于 <b>" + p.UpdateTime.Format("2006-01-02 15:04:05") + "</b>修改<br>" +
			"----------------------------------------------<br>" +
			"表: <b>" + p.MonitoringTable + "</b><br>" +
			"字段1: <b>" + p.OnlyField1Name + "</b> 值为: <b>" + p.OnlyField1Value + "</b><br>" +
			"字段2: <b>" + p.OnlyField2Name + "</b> 值为: <b>" + p.OnlyField2Value + "</b><br>" +
			"字段3: <b>" + p.OnlyField3Name + "</b> 值为: <b>" + p.OnlyField3Value + "</b><br>" +
			"字段4: <b>" + p.OnlyField4Name + "</b> 值为: <b>" + p.OnlyField4Value + "</b><br>" +
			"----------------------------------------------<br>" +
			"修改项为: <b>" + p.ChangedItem + "</b><br>" +
			"由<b>" + p.OldValue + "</b>更改为<b>" + p.NewValue + "</b><br>" +
			"更改原因为: <b>" + p.UpdateReason + "</b></html>"
		err := Services.EmailOfPtrChangeInfoKqPusher.Push(mes)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "Ok")
	}
}
