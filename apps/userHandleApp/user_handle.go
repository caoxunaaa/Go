package userHandleApp

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Middlewares"
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

func getAllProfileList() (profileList []*User.Profile, err error) {
	sqlStr := "SELECT id, username,password, is_superuser FROM profiles"
	err = Databases.SuperxonDbDevice.Select(&profileList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func AuthLoginHandler(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user User.Profile
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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

	c.JSON(http.StatusNotFound, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}

func AuthRegisterHandler(c *gin.Context) {
	var registerInfo User.RegisterInfo
	if err := c.ShouldBindJSON(&registerInfo); err == nil {
		err = User.Register(&registerInfo)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"DeviceName": registerInfo.Username,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
