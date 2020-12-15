package Router

import (
	"SuperxonWebSite/Middlewares"
	"SuperxonWebSite/apps/equipment"
	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

var options []Option

// 注册app的路由配置
func Include(opts ...Option) {
	options = append(options, opts...)
}

// 初始化
func Init() *gin.Engine {
	r := gin.Default()
	r.Use(Middlewares.Cors())
	v1 := r.Group("/product")
	{
		v1.GET("/", equipment.GetProductList)
		v1.GET("/productInfo/:pn", equipment.GetModuleInfoList)
	}
	return r
}
