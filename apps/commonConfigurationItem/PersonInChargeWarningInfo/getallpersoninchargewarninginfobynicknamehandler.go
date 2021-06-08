package PersonInChargeWarningInfo

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllPersonInChargeWarningInfoByNicknameHandler(c *gin.Context) {
	nickname := c.Param("nickname")
	res, err := ProductionLineOracleRelation.FindAllPersonInChargeWarningInfoByNickname(nickname)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
