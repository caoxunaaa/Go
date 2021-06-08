package OsaRunDisplay

import (
	"SuperxonWebSite/Models/ProductionLineOracleRelation"
)

func GetOsaWaringInfo(workOrderType, pn, startTime, endTime string) ([]struct {
	OsaInfo
	YellowLine int
	RedLine    int
}, error) {
	moduleOsa := "osa"
	//所有Osa良率
	var qsis = make([]OsaInfo, 0)
	qsis, err := GetOsaInfoList(&OsaQueryCondition{
		Pn:        pn,
		StartTime: startTime,
		EndTime:   endTime,
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
		OsaInfo
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
			OsaInfo
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
				if (swts[si].Pn == res[ri].Pn.String) && (swts[si].Process == res[ri].Process) {
					res[ri].YellowLine = swts[si].YellowLine
					res[ri].RedLine = swts[si].RedLine
					break
				}
			}
		}
	}
	return res, nil
}
