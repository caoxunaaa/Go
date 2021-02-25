package Utils

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
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
