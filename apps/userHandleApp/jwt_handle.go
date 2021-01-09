package userHandleApp

import (
	"SuperxonWebSite/Middlewares"
	"SuperxonWebSite/Models/User"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthJwtHandler(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user User.Profile
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}
	// 校验用户名和密码是否正确
	if user.Username == "superxon" && user.Password == "superxon" {
		// 生成Token
		tokenString, _ := Middlewares.GenToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "success",
			"data": gin.H{"token": tokenString},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}
