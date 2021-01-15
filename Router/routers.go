package Router

import (
	"SuperxonWebSite/Middlewares"
	"SuperxonWebSite/apps/deviceManangeApp"
	"SuperxonWebSite/apps/fileManage"
	"SuperxonWebSite/apps/qaStatisticBroad"
	"SuperxonWebSite/apps/runDisplayBroad"
	"SuperxonWebSite/apps/userHandleApp"
	"github.com/gin-gonic/gin"
)

// 初始化
func Init() *gin.Engine {
	r := gin.Default()
	r.Static("/assets", "./assets")
	//r.StaticFS("/assets", http.Dir("assets"))
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.Use(Middlewares.Cors())
	v1 := r.Group("/runDisplayBroad") //实时运行看板页面
	{
		v1.GET("/moduleList", runDisplayBroad.GetModuleListHandler)
		v1.GET("/moduleInfo/:pn", runDisplayBroad.GetModuleInfoListHandler)
		v1.GET("/osaList", runDisplayBroad.GetOsaListHandler)
		v1.GET("/osaInfo/:pn", runDisplayBroad.GetOsaInfoListHandler)
		v1.GET("/moduleYesterdayInfo/:pn", runDisplayBroad.GetYesterdayModuleInfoListHandler)
		v1.GET("/osaYesterdayInfo/:pn", runDisplayBroad.GetYesterdayOsaInfoListHandler)
		v1.GET("/stationStatus", runDisplayBroad.GetStationStatusHandler)
		v1.GET("/projectPlanList", runDisplayBroad.GetProjectPlanListHandler)
		v1.GET("/wipInfoList/:pn", runDisplayBroad.GetWipInfoListHandler)
	}
	v2 := r.Group("/qaStatisticBroad") //质量统计查询页面
	{
		v2.GET("/qaPnList", qaStatisticBroad.GetQaPnListHandler)
		v2.GET("/qaStatisticsInfo", qaStatisticBroad.GetQaStatisticInfoListHandler)
		v2.GET("/qaDefectsInfo", qaStatisticBroad.GetQaDefectsInfoListHandler)
		v2.GET("/pnSetParams", qaStatisticBroad.GetPnSetParamsListHandler)
		v2.GET("/pnWorkOrderYields", qaStatisticBroad.GetPnWorkOrderYieldsListHandler)
		v2.GET("/qaCpkInfo", qaStatisticBroad.GetQaCpkInfoListHandler)
		v2.GET("/qaCpkRssi", qaStatisticBroad.GetQaCpkRssiListHandler)
	}
	v3 := r.Group("/deviceManage") //设备管理页面
	{
		v3.GET("/deviceRootCategory", deviceManangeApp.GetAllDeviceCategoryRootListHandler)
		v3.GET("/deviceChildCategory/:rootCategory", deviceManangeApp.GetAllDeviceCategoryChildListHandler)
		v3.POST("/deviceChildCategory", deviceManangeApp.CreateDeviceCategoryChildHandler)

		v3.GET("/deviceBaseInfo", deviceManangeApp.GetAllDeviceBaseInfoListHandler)
		v3.GET("/deviceBaseInfo/:snAssetsIc", deviceManangeApp.GetDeviceBaseInfoHandler)
		v3.POST("/deviceBaseInfo", deviceManangeApp.CreateDeviceBaseInfoHandler)
		v3.DELETE("/deviceBaseInfo/:deviceSn", deviceManangeApp.DeleteDeviceBaseInfoHandler)
		v3.PUT("/deviceBaseInfo/:deviceSn", deviceManangeApp.UpdateDeviceBaseInfoHandler)

		v3.GET("/deviceTransmit", deviceManangeApp.GetAllDeviceTransmitInfoListHandler)
		v3.GET("/deviceTransmit/:deviceSn", deviceManangeApp.GetDeviceTransmitInfoHandler)
		v3.POST("/deviceTransmit", deviceManangeApp.CreateDeviceTransmitInfoHandler)
		v3.DELETE("/deviceTransmit/:id", deviceManangeApp.DeleteDeviceTransmitInfoHandler)

		v3.GET("/deviceRepair", deviceManangeApp.GetAllDeviceRepairInfoListHandler)
		v3.GET("/deviceRepair/:deviceSn", deviceManangeApp.GetDeviceRepairInfoHandler)

		v3.GET("/deviceMaintenanceItem", deviceManangeApp.GetAllDeviceMaintenanceCategoryListHandler)
		v3.GET("/deviceMaintenanceItem/:category", deviceManangeApp.GetDeviceMaintenanceItemOfCategoryHandler)

		v3.GET("/deviceCurrentMaintenance", deviceManangeApp.GetAllDeviceMaintenanceCurrentInfoListHandler)
		v3.GET("/deviceCurrentMaintenance/:snAssets", deviceManangeApp.GetDeviceMaintenanceCurrentInfoHandler)

		v3.GET("/deviceMaintenanceRecord", deviceManangeApp.GetAllDeviceMaintenanceAllRecordsHandler)
		v3.GET("/deviceMaintenanceRecord/:itemName", deviceManangeApp.GetAllDeviceMaintenanceRecordsOfItemNameHandler)

		v3.GET("/deviceMaintenanceRecordOfDevice/:snAssets", deviceManangeApp.GetDeviceMaintenanceRecordsHandler)
		v3.GET("/deviceMaintenanceRecordOfDevice/:snAssets/:itemName", deviceManangeApp.GetDeviceMaintenanceRecordOfItemNameHandler)
	}
	v4 := r.Group("/userHandle")
	{
		v4.GET("/profile", userHandleApp.GetAllProfileListHandler)
		v4.POST("/auth", userHandleApp.AuthJwtHandler)
	}
	v5 := r.Group("/home").Use(Middlewares.JWTAuthMiddleware())
	{
		v5.GET("/home", userHandleApp.HomeHandler)
	}
	v6 := r.Group("/fileManage") //视频管理页面
	{
		v6.GET("/videoInfo", fileManage.GetVideoInfoListHandler)
		v6.POST("/videoInfo", fileManage.UploadVideoFileHandler)
		v6.DELETE("/videoInfo/:id", fileManage.DeleteVideoInfoHandler)
	}
	return r
}
