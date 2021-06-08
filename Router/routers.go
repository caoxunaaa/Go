package Router

import (
	"SuperxonWebSite/Middlewares"
	"SuperxonWebSite/apps/commonConfigurationItem"
	"SuperxonWebSite/apps/commonConfigurationItem/PersonInChargeWarningInfo"
	"SuperxonWebSite/apps/commonConfigurationItem/SettingWarningThreshold"
	"SuperxonWebSite/apps/humanResources"
	"SuperxonWebSite/apps/moduleRunning"
	"SuperxonWebSite/apps/moduleStatistic"
	"SuperxonWebSite/apps/osaRunning"
	"SuperxonWebSite/apps/osaStatistic"
	"SuperxonWebSite/apps/trendCharts"
	"github.com/gin-gonic/gin"
)

// 初始化
func Init() *gin.Engine {
	r := gin.Default()

	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.Use(Middlewares.Cors())

	// 通用资源
	common := r.Group("/common-source")
	{
		//所有在用模块Pn
		common.GET("/all-module-pn-list", moduleRunning.GetModuleAllUsingPnHandler)
		//所有在用OSA Pn
		common.GET("/all-osa-pn-list", osaRunning.GetOsaAllUsingPnHandler)
		//某个时间段的所有模块Pn
		common.GET("/all-module-pn-list-in-time-period", moduleStatistic.GetModuleAllPnByTimeHandler)
		//某个时间段的所有Osa Pn
		common.GET("/all-osa-pn-list-in-time-period", osaStatistic.GetOsaAllPnByTimeHandler)
		//告警负责人
		common.GET("/person-in-charge-warning-info", PersonInChargeWarningInfo.GetAllPersonInChargeWarningInfoHandler)
		common.GET("/person-in-charge-warning-info/:nickname", PersonInChargeWarningInfo.GetAllPersonInChargeWarningInfoByNicknameHandler)
		common.POST("/person-in-charge-warning-info", PersonInChargeWarningInfo.CreatePersonInChargeWarningInfoHandler)
		common.PUT("/person-in-charge-warning-info/:id", PersonInChargeWarningInfo.UpdatePersonInChargeWarningInfoHandler)
		common.DELETE("/person-in-charge-warning-info/:id", PersonInChargeWarningInfo.DeletePersonInChargeWarningInfoHandler)
		//获取所有的告警门限设置
		common.GET("/settings-warning-threshold", SettingWarningThreshold.GetAllSettingWarningThresholdHandler)
		common.GET("/settings-warning-threshold/:id", SettingWarningThreshold.GetSettingWarningThresholdHandler)
		common.POST("/settings-warning-threshold", SettingWarningThreshold.CreateSettingWarningThresholdHandler)
		common.PUT("/settings-warning-threshold/:id", SettingWarningThreshold.UpdateSettingWarningThresholdHandler)
		common.DELETE("/settings-warning-threshold/:id", SettingWarningThreshold.DeleteSettingWarningThresholdHandler)
	}
	// 产线生产情况
	productInfo := r.Group("/product-operation-info")
	{
		// 1.模块端信息
		moduleOfProductInfo := productInfo.Group("/module")
		{
			//模块总览
			overview := moduleOfProductInfo.Group("/overview")
			{
				//某个模块pn的信息-当天0点到当前时间
				overview.GET("/module-info-today/:pn", moduleRunning.GetModuleInfoTodayByPnHandler)
				//所有模块的信息-当天0点到当前时间
				overview.GET("/all-module-info-today", moduleRunning.GetModuleAllInfoTodayHandler)
				//昨日某个模块信息
				overview.GET("/module-info-yesterday/:pn", moduleRunning.GetModuleInfoYesterdayHandler)
				//某个模块wip信息-当天0点到当前时间
				overview.GET("/wip-of-module-info-today/:pn", moduleRunning.GetModuleWipHandler)
				//所有工位的模块生产信息--当天0点到当前时间
				overview.GET("/station-product-info-of-module-today", moduleRunning.GetModuleAllStationStatusHandler)
			}
			//模块告警
			warningView := moduleOfProductInfo.Group("warning-view")
			{
				//通过工单类型获取某个时间段的告警信息
				warningView.GET("/warning-info-in-time-period-by-work-order-type", moduleRunning.GetModuleAllWaringInfoHandler)
				//工位动态告警
				warningView.GET("/station-dynamic-warning-monitoring", moduleRunning.GetModuleAllStationDynamicWarningMonitoringHandler)
			}
		}
		// 2.OSA端信息
		osaOfProductInfo := productInfo.Group("/osa")
		{
			//OSA总览
			overview := osaOfProductInfo.Group("/overview")
			{
				//某个OSA的信息-当天0点到当前时间
				overview.GET("/osa-info-today/:pn", osaRunning.GetOsaInfoTodayByPnHandler)
				//所有OSA的信息-当天0点到当前时间
				overview.GET("/all-osa-info-today", osaRunning.GetOsaAllInfoTodayHandler)
				//所有OSA发端耦合信息-当天0点到当前时间
				overview.GET("/all-osa-tx-couple-info-today", osaRunning.GetAllOsaTxCoupleInfoListHandler)
				//昨日某个OSA信息
				overview.GET("/osa-info-yesterday/:pn", osaRunning.GetYesterdayOsaInfoListHandler)
				//所有工位的OSA生产信息--当天0点到当前时间
				overview.GET("/station-product-info-of-osa-today", osaRunning.GetOsaAllStationStatusHandler)
				//某个OSA wip信息-当天0点到当前时间
				overview.GET("/wip-of-osa-info-today/:pn", osaRunning.GetOsaWipHandler)
			}
		}
		// 3.本月生产计划及其变更
		planOfProductInfo := productInfo.Group("/product-plan")
		{
			//计划完成情况
			planOfProductInfo.GET("/completed-situation", commonConfigurationItem.GetProjectPlanListHandler)
			//计划任务
			planOfProductInfo.GET("/plan-info", commonConfigurationItem.GetUndoneProjectPlanInfoListHandler)
			planOfProductInfo.POST("/plan-info", commonConfigurationItem.CreateUndoneProjectPlanInfoHandler)
			planOfProductInfo.PUT("/plan-info/:id", commonConfigurationItem.UpdateUndoneProjectPlanInfoHandler)
			planOfProductInfo.DELETE("/plan-info/:id", commonConfigurationItem.DeleteUndoneProjectPlanInfoHandler)
		}
	}
	// 统计查询
	statisticQuery := r.Group("/product-statistic-query")
	{
		//模块端统计查询
		moduleOfStatisticQuery := statisticQuery.Group("/module")
		{
			//获取产品配置参数
			moduleOfStatisticQuery.GET("/settings-params", moduleStatistic.GetModuleSettingParamHandler)
			//通过pn 工序获取产品良率
			moduleOfStatisticQuery.GET("/production-yield-info-by-pn", moduleStatistic.GetModuleYieldInfoByPnHandler)
			//通过pn 工序获取产品良率带工单号和版本号
			moduleOfStatisticQuery.GET("/production-yield-info-with-order-and-version-by-pn", moduleStatistic.GetModuleYieldInfoWithWorkOrderIdByPnHandler)
			//通过Pn获取不良代码分布
			moduleOfStatisticQuery.GET("/production-bad-code-distribution-by-pn", moduleStatistic.GetModuleDefectsInfoByPnHandler)
			//通过Pn获取不良代码详情
			moduleOfStatisticQuery.GET("/production-bad-code-detail-by-pn", moduleStatistic.GetModuleDefectsDetailInfoByPnHandler)

			//通过pn获取某段时间工单号-分结案和未结案
			moduleOfStatisticQuery.GET("/all-production-work-order-id", moduleStatistic.GetModuleWorkOrderIdHandler)
			//通过工单号获取产品良率
			moduleOfStatisticQuery.GET("/production-yield-info-by-work-order-id", moduleStatistic.GetModuleYieldInfoByWorkOrderIdHandler)
			//通过工单号获取不良代码分布
			moduleOfStatisticQuery.GET("/production-bad-code-distribution-by-work-order-id", moduleStatistic.GetModuleDefectsInfoByWorkOrderIdHandler)
			//通过工单号获取不良代码详情
			moduleOfStatisticQuery.GET("/production-bad-code-detail-by-work-order-id", moduleStatistic.GetQaDefectsDetailByWorkOrderIdHandler)
			//常规项 CPK
			moduleOfStatisticQuery.GET("/normal-items-cpk", moduleStatistic.GetModuleNormalItemCpkHandler)
			//RSSI CPK
			moduleOfStatisticQuery.GET("/rssi-cpk", moduleStatistic.GetModuleRssiCpkHandler)
		}
		//OSA端统计查询
		osaOfStatisticQuery := statisticQuery.Group("/osa")
		{
			//通过Pn工序获取产品良率
			osaOfStatisticQuery.GET("/production-yield-info-by-pn", osaStatistic.GetOsaYieldInfoByPnHandler)
			//通过Pn获取不良代码分布
			osaOfStatisticQuery.GET("/production-bad-code-distribution-by-pn", osaStatistic.GetOsaDefectsInfoByPnHandler)
			//通过工单号获取产品良率
			osaOfStatisticQuery.GET("/production-yield-info-by-work-order-id", osaStatistic.GetOsaYieldInfoByWorkOrderIdHandler)
		}
		//产线投入产出汇总
		statisticQuery.GET("/input-and-output-summary", moduleStatistic.GetInputAndOutputSummaryInfoListHandler)
	}
	//统计趋势图
	charts := r.Group("/trend-chart")
	{
		//每日告警统计
		charts.GET("/warning-statistic-daily", trendCharts.GetWarningStatisticDailyHandler)
	}
	//人力资源
	humanResource := r.Group("/human-resources")
	{
		//文件管理
		fileManage := humanResource.Group("/file-manage")
		{
			fileManage.GET("/videoInfo", humanResources.GetVideoInfoListHandler).Use(Middlewares.JWTAuthMiddleware())
			fileManage.POST("/videoInfo", humanResources.UploadVideoFileHandler).Use(Middlewares.JWTSuperuserMiddleware())
			fileManage.DELETE("/videoInfo/:id", humanResources.DeleteVideoInfoHandler).Use(Middlewares.JWTSuperuserMiddleware())
		}
	}
	return r
}
