package SettingWarningThreshold

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateSettingWarningThresholdHandler(c *gin.Context) {
	var err error
	var p ProductionLineOracleRelation.SettingWarningThreshold
	p.OrderType = c.PostForm("order_type")
	p.ModuleOsa = c.PostForm("module_osa")
	p.Pn = c.PostForm("pn")
	p.Process = c.PostForm("process")
	fmt.Println()
	p.YellowLine, err = strconv.Atoi(c.PostForm("yellow_line"))
	p.RedLine, err = strconv.Atoi(c.PostForm("red_line"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = ProductionLineOracleRelation.CreateSettingWarningThreshold(&p)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, "Ok")
	}
}
