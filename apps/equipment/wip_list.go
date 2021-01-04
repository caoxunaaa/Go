package equipment

import (
	"SuperxonWebSite/Models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWipInfoList(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的PN"})
	}
	wipInfoList, err := Models.GetWipInfoList(pn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, wipInfoList)
	}
}
