package SettingWarningThreshold

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UpdateSettingWarningThresholdHandler(c *gin.Context) {
	var err error
	var p ProductionLineOracleRelation.SettingWarningThreshold
	id, err := strconv.Atoi(c.Param("id"))
	p.OrderType = c.PostForm("order_type")
	p.ModuleOsa = c.PostForm("module_osa")
	p.Pn = c.PostForm("pn")
	p.Process = c.PostForm("process")
	p.YellowLine, err = strconv.Atoi(c.PostForm("yellow_line"))
	p.RedLine, err = strconv.Atoi(c.PostForm("red_line"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = ProductionLineOracleRelation.UpdateSettingWarningThreshold(&p, int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, "Ok")
	}
}
