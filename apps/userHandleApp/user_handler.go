package userHandleApp

import (
	"SuperxonWebSite/Middlewares"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthLoginHandler(c *gin.Context) {
	var user struct {
		Id       int64  `db:"id"`
		Username string `db:"username"`
		Phone    string `db:"phone"`
		Nickname string `db:"nickname"`
		Password string `db:"password"`
		Email    string `db:"email"`
	}
	err := c.ShouldBindJSON(&user)
	fmt.Println(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}

	tokenString, _ := Middlewares.GenToken(user.Username)
	c.JSON(http.StatusOK, gin.H{
		"code":     2000,
		"msg":      "success",
		"Token":    tokenString,
		"username": user.Username,
		"nickname": user.Username,
	})
	return
	c.JSON(http.StatusNotFound, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}
