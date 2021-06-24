package ModuleRunDisplay

import (
	"SuperxonWebSite/Models/ModuleStatisticDisplay"
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
	"strconv"
)

type ModuleWarningInfo struct {
	Pn            string
	Sequence      string
	Process       string
	TotalInput    uint32
	FinalOk       uint32
	FinalBad      uint32
	FinalPassRate float64
}

//查看所有的模块良率告警信息 pn="%%"
func GetModuleWaringInfo(workOrderType, pn, startTime, endTime string) ([]struct {
	ModuleStatisticDisplay.QaStatisticInfo
	YellowLine int
	RedLine    int
}, error) {
	moduleOsa := "module"
	//所有模块良率
	var qsis = make([]ModuleStatisticDisplay.QaStatisticInfo, 0)
	qsis, err := ModuleStatisticDisplay.GetQaStatisticOrderInfoList(&ModuleStatisticDisplay.QueryCondition{
		Pn:            pn,
		WorkOrderType: workOrderType,
		StartTime:     startTime,
		EndTime:       endTime,
	})
	if err != nil {
		return nil, err
	}
	//获取良率的告警配置
	var swts = make([]ProductionLineOracleRelation.SettingWarningThreshold, 0)
	swts, err = ProductionLineOracleRelation.FindSomeSettingWarningThresholdByOrderTypeAndModuleOsa(workOrderType, moduleOsa)
	if err != nil {
		return nil, err
	}

	res := make([]struct {
		ModuleStatisticDisplay.QaStatisticInfo
		YellowLine int
		RedLine    int
	}, 0)

	//进行告警的适配，找出符合要求的
	//全设成默认值
	swtDefault, err := ProductionLineOracleRelation.FindDefaultSettingWarningThresholdByOrderTypeAndModuleOsa(workOrderType, moduleOsa)
	if err != nil {
		return nil, err
	}
	for qi := 0; qi < len(qsis); qi++ {
		res = append(res, struct {
			ModuleStatisticDisplay.QaStatisticInfo
			YellowLine int
			RedLine    int
		}{qsis[qi], swtDefault.YellowLine, swtDefault.RedLine})
	}
	//遍历所有结果，找到Pn和Process对应的设置相应的值
	for si := 0; si < len(swts); si++ {
		if swts[si].Pn == "DEFAULT" {
			continue
		} else {
			for ri := 0; ri < len(res); ri++ {
				if (swts[si].Pn == res[ri].Pn) && (swts[si].Process == res[ri].Process) {
					res[ri].YellowLine = swts[si].YellowLine
					res[ri].RedLine = swts[si].RedLine
					break
				}
			}
		}
	}
	// 找出所有良率低于红线的结果
	result := make([]struct {
		ModuleStatisticDisplay.QaStatisticInfo
		YellowLine int
		RedLine    int
	}, 0)
	for ri := 0; ri < len(res); ri++ {
		fpr, err := strconv.ParseFloat(res[ri].FinalPassRate, 64)
		if err != nil {
			return nil, err
		}
		if fpr < float64(res[ri].YellowLine) {
			result = append(result, res[ri])
		}
	}
	return result, nil
}

//查看所有的模块良率告警信息 pn="%%"
func GetModuleWaringInfoWithStation(workOrderType, pn, startTime, endTime string) ([]struct {
	StationStatus
	YellowLine int
	RedLine    int
}, error) {
	moduleOsa := "module"
	//所有模块良率
	var qsis = make([]StationStatus, 0)
	qsis, err := GetStationStatus(pn, workOrderType, startTime, endTime)
	if err != nil {
		return nil, err
	}
	//获取良率的告警配置
	var swts = make([]ProductionLineOracleRelation.SettingWarningThreshold, 0)
	swts, err = ProductionLineOracleRelation.FindSomeSettingWarningThresholdByOrderTypeAndModuleOsa(workOrderType, moduleOsa)
	if err != nil {
		return nil, err
	}

	res := make([]struct {
		StationStatus
		YellowLine int
		RedLine    int
	}, 0)

	//进行告警的适配，找出符合要求的
	//全设成默认值
	swtDefault, err := ProductionLineOracleRelation.FindDefaultSettingWarningThresholdByOrderTypeAndModuleOsa(workOrderType, moduleOsa)
	if err != nil {
		return nil, err
	}
	for qi := 0; qi < len(qsis); qi++ {
		res = append(res, struct {
			StationStatus
			YellowLine int
			RedLine    int
		}{qsis[qi], swtDefault.YellowLine, swtDefault.RedLine})
	}
	//遍历所有结果，找到Pn和Process对应的设置相应的值
	for si := 0; si < len(swts); si++ {
		if swts[si].Pn == "DEFAULT" {
			continue
		} else {
			for ri := 0; ri < len(res); ri++ {
				if (swts[si].Pn == res[ri].Pn) && (swts[si].Process == res[ri].LogAction) {
					res[ri].YellowLine = swts[si].YellowLine
					res[ri].RedLine = swts[si].RedLine
					break
				}
			}
		}
	}
	// 找出所有良率低于红线的结果
	result := make([]struct {
		StationStatus
		YellowLine int
		RedLine    int
	}, 0)
	for ri := 0; ri < len(res); ri++ {
		fpr, err := strconv.ParseFloat(res[ri].PassRate, 64)
		if err != nil {
			return nil, err
		}
		if fpr < float64(res[ri].YellowLine) {
			result = append(result, res[ri])
		}
	}
	return result, nil
}
