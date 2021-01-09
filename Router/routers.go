package Router

import (
	"SuperxonWebSite/Middlewares"
	"SuperxonWebSite/apps/deviceManangeApp"
	"SuperxonWebSite/apps/qaStatisticBroad"
	"SuperxonWebSite/apps/runDisplayBroad"
	"github.com/gin-gonic/gin"
)

//type Option func(*gin.Engine)

//var options []Option

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
	v3 := r.Group("/deviceManage")
	{
		v3.GET("/deviceRootCategory", deviceManangeApp.GetAllDeviceCategoryRootList)
		v3.GET("/deviceChildCategory/:rootCategory", deviceManangeApp.GetAllDeviceCategoryChildList)
		v3.POST("/deviceChildCategory", deviceManangeApp.CreateDeviceCategoryChild)

		v3.GET("/deviceBaseInfo", deviceManangeApp.GetAllDeviceBaseInfoList)
		v3.GET("/deviceBaseInfo/:snAssetsIc", deviceManangeApp.GetDeviceBaseInfo)

		v3.GET("/deviceTransmit", deviceManangeApp.GetAllDeviceTransmitInfoList)

		v3.GET("/deviceRepair", deviceManangeApp.GetAllDeviceRepairInfoList)
		v3.GET("/deviceRepair/:deviceSn", deviceManangeApp.GetDeviceRepairInfo)

		v3.GET("/deviceMaintenanceItem", deviceManangeApp.GetAllDeviceMaintenanceCategoryList)
		v3.GET("/deviceMaintenanceItem/:category", deviceManangeApp.GetDeviceMaintenanceItemOfCategory)

		v3.GET("/deviceCurrentMaintenance", deviceManangeApp.GetAllDeviceMaintenanceCurrentInfoList)
		v3.GET("/deviceCurrentMaintenance/:snAssets", deviceManangeApp.GetDeviceMaintenanceCurrentInfo)

		v3.GET("/deviceMaintenanceRecord", deviceManangeApp.GetAllDeviceMaintenanceAllRecords)
		v3.GET("/deviceMaintenanceRecord/:itemName", deviceManangeApp.GetAllDeviceMaintenanceRecordsOfItemName)

		v3.GET("/deviceMaintenanceRecordOfDevice/:snAssets", deviceManangeApp.GetDeviceMaintenanceRecords)
		v3.GET("/deviceMaintenanceRecordOfDevice/:snAssets/:itemName", deviceManangeApp.GetDeviceMaintenanceRecordOfItemName)
	}
	return r
}
