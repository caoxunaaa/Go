package PersonInChargeWarningInfo

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UpdatePersonInChargeWarningInfoHandler(c *gin.Context) {
	var p ProductionLineOracleRelation.PersonInChargeWarningInfo
	id, err := strconv.Atoi(c.Param("id"))
	p.Username = c.PostForm("username")
	p.Nickname = c.PostForm("nickname")
	p.ModuleOsa = c.PostForm("module_osa")
	p.Pn = c.PostForm("pn")
	err = ProductionLineOracleRelation.UpdatePersonInChargeWarningInfo(&p, int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, "Ok")
	}
}
