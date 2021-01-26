package Utils

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

func InitCommonIni(path string) (cfg *ini.File, err error) {
	cfg, err = ini.Load(path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	return
}

func GetCommonPnList() (pnList []string, err error) {
	dir, _ := os.Getwd()
	fmt.Println(dir)
	cfg, err := InitCommonIni(dir + "\\Services\\Common.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return
	}
	pnList = strings.Split(cfg.Section("PnForCpkRedis").Key("PnList").String(), ",")
	return
}
