package Router

import (
	"SuperxonWebSite/Middlewares"
	"SuperxonWebSite/apps/qaStatisticBroad"
	"SuperxonWebSite/apps/runDisplayBroad"
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
	v1 := r.Group("/runDisplayBroad")
	{
		v1.GET("/moduleList", runDisplayBroad.GetModuleList)
		v1.GET("/moduleInfo/:pn", runDisplayBroad.GetModuleInfoList)
		v1.GET("/osaList", runDisplayBroad.GetOsaList)
		v1.GET("/osaInfo/:pn", runDisplayBroad.GetOsaInfoList)
		v1.GET("/moduleYesterdayInfo/:pn", runDisplayBroad.GetYesterdayModuleInfoList)
		v1.GET("/osaYesterdayInfo/:pn", runDisplayBroad.GetYesterdayOsaInfoList)
		v1.GET("/stationStatus", runDisplayBroad.GetStationStatus)
		v1.GET("/projectPlanList", runDisplayBroad.GetProjectPlanList)
		v1.GET("/wipInfoList/:pn", runDisplayBroad.GetWipInfoList)
	}
	v2 := r.Group("/qaStatisticBroad")
	{
		v2.GET("/qaPnList", qaStatisticBroad.GetQaPnList)
		v2.GET("/qaStatisticsInfo", qaStatisticBroad.GetQaStatisticInfoList)
		v2.GET("/qaDefectsInfo", qaStatisticBroad.GetQaDefectsInfoList)
	}
	return r
}
