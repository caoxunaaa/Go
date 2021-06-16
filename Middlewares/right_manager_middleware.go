package Middlewares

import (
	"SuperxonWebSite/Models/UserRelation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ProductionParameterRightMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		u, exists := c.Get("username")
		username, ok := u.(string)
		if ok != true {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2003,
				"msg":  "用户没有登录",
			})
			c.Abort()
			return
		}
		if exists == false || username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2003,
				"msg":  "用户没有登录",
			})
			c.Abort()
			return
		}
		//获取username是否有修改产品工艺权限
		res, err := UserRelation.FindAllRightManagerByUsername(username)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 2005,
				"msg":  "没有修改产品工艺权限1",
			})
			c.Abort()
			return
		}
		if len(res) <= 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 2005,
				"msg":  "没有修改产品工艺权限2",
			})
			c.Abort()
			return
		}

		for _, re := range res {
			if re.RightItem == "生产工艺参数" {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{
			"code": 2005,
			"msg":  "没有修改产品工艺权限3",
		})
		c.Abort()
		return
	}
}
