package Router

import (
	"SuperxonWebSite/Middlewares"
	"SuperxonWebSite/apps/deviceManangeApp"
	"SuperxonWebSite/apps/fileManage"
	"SuperxonWebSite/apps/qaModuleStatisticBroad"
	"SuperxonWebSite/apps/qaOsaStatisticBroad"
	"SuperxonWebSite/apps/runModuleDisplayBroad"
	"SuperxonWebSite/apps/runOsaDisplayBroad"
	"SuperxonWebSite/apps/userHandleApp"
	"SuperxonWebSite/apps/waringDisplayBroad"
	"github.com/gin-gonic/gin"
)

// 初始化
func Init() *gin.Engine {
	r := gin.Default()
	//pprof.Register(r)
	r.Static("/assets", "./assets")
	//r.StaticFS("/assets", http.Dir("assets"))
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.Use(Middlewares.Cors())

	//趋势图表
	vCharts := r.Group("/TrendChart")
	{
		vCharts.GET("/warningCharts", waringDisplayBroad.GetWarningCountChartDataListHandler)
	}

	//实时运行看板页面
	v1 := r.Group("/runDisplayBroad")
	{
		v1.GET("/moduleList", runModuleDisplayBroad.GetModuleListHandler)
		v1.GET("/moduleInfo/:pn", runModuleDisplayBroad.GetModuleInfoListHandler)

		v1.GET("/allModuleInfo", runModuleDisplayBroad.GetAllModuleInfoListHandler)
		v1.GET("/allOsaInfo", runOsaDisplayBroad.GetAllOsaInfoListHandler)
		v1.GET("/allOsaTxCoupleInfo", runOsaDisplayBroad.GetAllOsaTxCoupleInfoListHandler)
		v1.GET("/warningCorrespondingTable", waringDisplayBroad.GetWarningToPersonListHandler)

		v1.GET("/moduleYesterdayInfo/:pn", runModuleDisplayBroad.GetYesterdayModuleInfoListHandler)
		v1.GET("/moduleStationStatus", runModuleDisplayBroad.GetStationStatusHandler)
		v1.GET("/moduleProjectPlanList", runModuleDisplayBroad.GetProjectPlanListHandler)
		v1.GET("/UndoneProjectPlanList", runModuleDisplayBroad.GetUndoneProjectPlanInfoListHandler)
		v1.GET("/moduleWipInfoList/:pn", runModuleDisplayBroad.GetWipInfoListHandler)

		v1.GET("/osaList", runOsaDisplayBroad.GetOsaListHandler)
		v1.GET("/osaInfoList/:pn", runOsaDisplayBroad.GetOsaInfoListHandler)
		v1.GET("/osaYesterdayInfo/:pn", runOsaDisplayBroad.GetYesterdayOsaInfoListHandler)
		v1.GET("/osaStationDetail", runOsaDisplayBroad.GetStationStatusHandler)
		v1.GET("/osaWipInfoList/:pn", runOsaDisplayBroad.GetOsaWipInfoListHandler)
	}
	v1Permission := r.Group("/runDisplayBroad").Use(Middlewares.JWTSuperuserMiddleware())
	{
		v1Permission.POST("/UndoneProjectPlanList", runModuleDisplayBroad.CreateUndoneProjectPlanInfoHandler)
		v1Permission.PUT("/UndoneProjectPlanList/:id", runModuleDisplayBroad.UpdateUndoneProjectPlanInfoHandler)
		v1Permission.DELETE("/UndoneProjectPlanList/:id", runModuleDisplayBroad.DeleteUndoneProjectPlanInfoHandler)
	}
	//模块质量统计查询页面
	v2 := r.Group("/qaStatisticBroad").Use(Middlewares.JWTAuthMiddleware())
	{
		v2.GET("/ioSummary", qaModuleStatisticBroad.Get10GLineIOSummaryInfoListHandler)

		v2.GET("/qaWorkOrderIdList", qaModuleStatisticBroad.GetWorkOrderIdsHandler)
		v2.GET("/qaWorkOrderYieldsByWorkOrderId", qaModuleStatisticBroad.GetWorkOrderYieldsByWorkOrderIdListHandler)
		v2.GET("/qaDefectsInfoByWorkOrderId", qaModuleStatisticBroad.GetQaDefectsInfoByWorkOrderIdListHandler)
		v2.GET("/qaDefectsDetailByWorkOrderId", qaModuleStatisticBroad.GetQaDefectsDetailByWorkOrderIdHandler)

		v2.GET("/qaPnList", qaModuleStatisticBroad.GetQaPnListHandler)
		v2.GET("/qaStatisticsInfo", qaModuleStatisticBroad.GetQaStatisticInfoListHandler)
		v2.GET("/pnWorkOrderYields", qaModuleStatisticBroad.GetPnWorkOrderYieldsListHandler)
		v2.GET("/qaDefectsInfo", qaModuleStatisticBroad.GetQaDefectsInfoListByPnHandler)
		v2.GET("/qaDefectsDetailByPn", qaModuleStatisticBroad.GetQaDefectsDetailByPnHandler)

		v2.GET("/qaCpkInfo", qaModuleStatisticBroad.GetQaCpkInfoListHandler)
		v2.GET("/qaCpkRssi", qaModuleStatisticBroad.GetQaCpkRssiListHandler)

		v2.GET("/pnSetParams", qaModuleStatisticBroad.GetPnSetParamsListHandler)
	}
	//OSA质量统计查询页面
	v2Osa := r.Group("/qaOsaStatisticBroad").Use(Middlewares.JWTAuthMiddleware())
	{
		v2Osa.GET("/qaOsaPnList", qaOsaStatisticBroad.GetQaOsaPnListHandler)
		v2Osa.GET("/qaOsaStatisticsInfo", qaOsaStatisticBroad.GetQaOsaStatisticInfoListHandler)
		v2Osa.GET("/qaOsaDefectsInfo", qaOsaStatisticBroad.GetQaOsaDefectsInfoListByPnHandler)
		v2Osa.GET("/qaOsaStatisticsInfoByWorkOrderId", qaOsaStatisticBroad.GetQaOsaStatisticInfoListByWorkOrderIdHandler)
	}
	//v3设备管理页面
	v3 := r.Group("/deviceManage").Use(Middlewares.JWTAuthMiddleware())
	{
		v3.GET("/deviceRootCategory", deviceManangeApp.GetAllDeviceCategoryRootListHandler)
		v3.GET("/deviceChildCategory/:rootCategory", deviceManangeApp.GetAllDeviceCategoryChildListHandler)

		v3.GET("/deviceBaseInfo", deviceManangeApp.GetAllDeviceBaseInfoListHandler)
		v3.GET("/deviceBaseInfo/:snAssetsIc", deviceManangeApp.GetDeviceBaseInfoHandler)

		v3.GET("/deviceTransmit", deviceManangeApp.GetAllDeviceTransmitInfoListHandler)
		v3.GET("/deviceTransmit/:deviceSn", deviceManangeApp.GetDeviceTransmitInfoHandler)

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
	//v3设备管理页面用户权限
	v3Permission := r.Group("/deviceManage").Use(Middlewares.JWTSuperuserMiddleware())
	{
		v3Permission.POST("/deviceChildCategory", deviceManangeApp.CreateDeviceCategoryChildHandler)

		v3Permission.POST("/deviceBaseInfo", deviceManangeApp.CreateDeviceBaseInfoHandler)
		v3Permission.DELETE("/deviceBaseInfo/:deviceSn", deviceManangeApp.DeleteDeviceBaseInfoHandler)
		v3Permission.PUT("/deviceBaseInfo/:deviceSn", deviceManangeApp.UpdateDeviceBaseInfoHandler)

		v3Permission.POST("/deviceTransmit", deviceManangeApp.CreateDeviceTransmitInfoHandler)
		v3Permission.DELETE("/deviceTransmit/:id", deviceManangeApp.DeleteDeviceTransmitInfoHandler)

		v3Permission.POST("/deviceRepair", deviceManangeApp.CreateCreateDeviceRepairInfoHandler)
		v3Permission.PUT("/deviceRepair/:id", deviceManangeApp.UpdateDeviceRepairInfoHandler)

		v3Permission.POST("/deviceMaintenanceItem", deviceManangeApp.CreateDeviceMaintenanceItemHandler)
		v3Permission.PUT("/deviceMaintenanceItem/:id", deviceManangeApp.UpdateDeviceMaintenanceItemHandler)
		v3Permission.DELETE("/deviceMaintenanceItem/:id", deviceManangeApp.DeleteDeviceMaintenanceItemHandler)

		v3Permission.POST("/bindDeviceMaintenanceItem/:deviceSn", deviceManangeApp.BindDeviceMaintenanceItemHandler)
		v3Permission.POST("/unbindDeviceMaintenanceItem/:deviceSn", deviceManangeApp.UnBindDeviceMaintenanceItemHandler)

		v3Permission.POST("/deviceMaintenanceRecord", deviceManangeApp.CreateDeviceMaintenanceRecordHandler)
	}
	//用户页面
	v4 := r.Group("/userHandle")
	{
		v4.GET("/profile", userHandleApp.GetAllProfileListHandler)
		v4.POST("/login", userHandleApp.AuthLoginHandler)
		v4.POST("/register", userHandleApp.AuthRegisterHandler)
	}
	v5 := r.Group("/home").Use(Middlewares.JWTAuthMiddleware())
	{
		v5.GET("/home", userHandleApp.HomeHandler)
	}
	//视频管理页面
	v6 := r.Group("/fileManage").Use(Middlewares.JWTAuthMiddleware())
	{
		v6.GET("/videoInfo", fileManage.GetVideoInfoListHandler)
	}
	//视频管理页面
	v6Permission := r.Group("/fileManage").Use(Middlewares.JWTSuperuserMiddleware())
	{
		v6Permission.POST("/videoInfo", fileManage.UploadVideoFileHandler)
		v6Permission.DELETE("/videoInfo/:id", fileManage.DeleteVideoInfoHandler)
	}
	return r
}
