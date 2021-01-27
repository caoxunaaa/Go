package userHandleApp

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Middlewares"
	"SuperxonWebSite/Models/User"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAllProfileList() (profileList []*User.Profile, err error) {
	sqlStr := "SELECT id, username,password, is_superuser FROM profiles"
	err = Databases.SuperxonDbDevice.Select(&profileList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func AuthJwtHandler(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user User.Profile
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}

	// 校验用户名和密码是否正确
	profileList, err := getAllProfileList()
	if err != nil {
		return
	}

	for _, profile := range profileList {
		if user.Username == profile.Username && user.Password == profile.Password {
			tokenString, _ := Middlewares.GenToken(user.Username)
			c.JSON(http.StatusOK, gin.H{
				"code":     2000,
				"msg":      "success",
				"Token":    tokenString,
				"username": profile.Username,
				"nickname": profile.Username,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}
