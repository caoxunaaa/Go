package main

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Router"
	"fmt"
)

func main() {
	var err error

	Databases.InitOracle()
	defer Databases.CloseOracle()

	Databases.InitMysql()
	defer Databases.CloseMysql()

	Databases.RedisInit()
	defer Databases.RedisClose()

	//Services.KafkaInit()
	//err = Services.InitCron()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer Services.CloseCron()

	r := Router.Init()
	if err = r.Run("0.0.0.0:8002"); err != nil {
		fmt.Printf("startup service failed, err:%v\n\n", err)
	}
}
