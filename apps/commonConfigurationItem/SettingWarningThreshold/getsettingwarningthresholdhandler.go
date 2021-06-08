package SettingWarningThreshold

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetSettingWarningThresholdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	res, err := ProductionLineOracleRelation.FindOneSettingWarningThresholdById(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
