package Utils

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strconv"
	"strings"
)

func InitIni(path string) (cfg *ini.File, err error) {
	cfg, err = ini.Load(path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	return
}

func GetCommonPnList() (pnList []string, err error) {
	dir, _ := os.Getwd()
	cfg, err := InitIni(dir + "\\Services\\Common.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return
	}
	pnList = strings.Split(cfg.Section("PnForCpkRedis").Key("PnList").String(), ",")
	return
}

func GetChartsPnList() (pnList []string, err error) {
	dir, _ := os.Getwd()
	cfg, err := InitIni(dir + "\\Services\\Common.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return
	}
	pnList = strings.Split(cfg.Section("PnForCharts").Key("PnList").String(), ",")
	return
}

func GetErrorCodeDescribe(errorCode string) (errorCodeDescribe string, err error) {
	dir, _ := os.Getwd()
	cfg, err := InitIni(dir + "\\Services\\ErrorCodeDescribe.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return
	}
	errorCodeDescribe = cfg.Section("ErrorCode").Key(errorCode).String()
	return
}

func GetWarningClassificationList() (warningClassificationList []string, err error) {
	dir, _ := os.Getwd()
	cfg, err := InitIni(dir + "\\Services\\Common.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return
	}
	warningClassificationList = strings.Split(cfg.Section("WarningClassification").Key("Classification").String(), ",")
	return
}

func GetWarningToPersonList() (map[string][]string, error) {
	personToPnManage := make(map[string][]string)
	dir, _ := os.Getwd()
	cfg, err := InitIni(dir + "\\Services\\Common.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return nil, nil
	}
	People := cfg.Section("WarningToPerson").KeyStrings()
	for _, person := range People {
		personToPnManage[person] = strings.Split(cfg.Section("WarningToPerson").Key(person).String(), ",")
	}
	return personToPnManage, err
}

func GetIOSummaryList(pn string) (pnCode string, err error) {
	dir, _ := os.Getwd()
	cfg, err := InitIni(dir + "\\Services\\Common.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return
	}
	pnCode = cfg.Section("10GLine").Key(pn).String()
	return
}

func GetWarningThresholdList() ([]interface{}, error) {
	res := make([]interface{}, 0)

	dir, _ := os.Getwd()
	cfg, err := InitIni(dir + "\\Services\\Common.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return nil, nil
	}
	keys := cfg.Section("WarningThreshold").KeyStrings()
	for i := 0; i < len(keys); i++ {
		var WarningThreshold map[string]interface{}
		WarningThreshold = make(map[string]interface{})
		tempKey := keys[i]
		tempKeyList := strings.Split(tempKey, "/")
		WarningThreshold["Pn"] = tempKeyList[0]
		WarningThreshold["Process"] = tempKeyList[1]
		tempValueList := strings.Split(cfg.Section("WarningThreshold").Key(tempKey).String(), ",")
		if len(tempValueList) != 2 {
			return nil, errors.New("配置文件中[WarningThreshold]不符合规范")
		}
		for _, value := range tempValueList {
			if ok := strings.Contains(value, "R"); ok {
				temp, err := strconv.Atoi(value[1:])
				if err != nil {
					return nil, errors.New("配置文件中[WarningThreshold]不符合规范")
				}
				WarningThreshold["R"] = temp
			}
			if ok := strings.Contains(value, "Y"); ok {
				temp, err := strconv.Atoi(value[1:])
				if err != nil {
					return nil, errors.New("配置文件中[WarningThreshold]不符合规范")
				}
				WarningThreshold["Y"] = temp
			}
		}
		res = append(res, WarningThreshold)
		fmt.Println(res)
	}
	return res, nil
}
