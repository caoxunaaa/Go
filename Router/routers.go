package Router

import (
	"SuperxonWebSite/Middlewares"
	"SuperxonWebSite/apps/equipment"
	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

var options []Option

// 注册app的路由配置
//func Include(opts ...Option) {
//	options = append(options, opts...)
//}

// 初始化
func Init() *gin.Engine {
	r := gin.Default()
	r.Static("/assets", "./assets")
	//r.StaticFS("/assets", http.Dir("assets"))
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.Use(Middlewares.Cors())
	v1 := r.Group("/product")
	{
		v1.GET("/moduleList", equipment.GetModuleList)
		v1.GET("/moduleInfo/:pn", equipment.GetModuleInfoList)
		v1.GET("/osaList", equipment.GetOsaList)
		v1.GET("/osaInfo/:pn", equipment.GetOsaInfoList)
		v1.GET("/moduleYesterdayInfo/:pn", equipment.GetYesterdayModuleInfoList)
		v1.GET("/osaYesterdayInfo/:pn", equipment.GetYesterdayOsaInfoList)
		v1.GET("/stationStatus", equipment.GetStationStatus)
		v1.GET("/projectPlanList", equipment.GetProjectPlanList)
	}
	return r
}
