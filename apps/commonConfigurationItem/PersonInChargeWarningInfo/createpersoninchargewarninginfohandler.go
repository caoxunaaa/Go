package PersonInChargeWarningInfo

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreatePersonInChargeWarningInfoHandler(c *gin.Context) {
	var p ProductionLineOracleRelation.PersonInChargeWarningInfo
	p.Username = c.PostForm("username")
	p.Nickname = c.PostForm("nickname")
	p.ModuleOsa = c.PostForm("module_osa")
	p.Pn = c.PostForm("pn")

	err := ProductionLineOracleRelation.CreatePersonInChargeWarningInfo(&p)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, "Ok")
	}
}
