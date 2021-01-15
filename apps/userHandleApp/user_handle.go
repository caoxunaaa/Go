package userHandleApp

import (
	"SuperxonWebSite/Models/User"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllProfileListHandler(c *gin.Context) {
	profileList, err := User.GetAllProfileList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, profileList)
	}
}
