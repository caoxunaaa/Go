package runDisplayBroad

import (
	"SuperxonWebSite/Models/RunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWipInfoList(c *gin.Context) {
	pn, ok := c.Params.Get("pn")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的PN"})
	}
	wipInfoList, err := RunDisplay.GetWipInfoList(pn)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, wipInfoList)
	}
}
