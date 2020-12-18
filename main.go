package main

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Router"
	"fmt"
)

func main() {
	Databases.InitOracle()
	defer Databases.CloseOracle()

	Databases.RedisInit()
	defer Databases.RedisClose()

	r := Router.Init()
	if err := r.Run("0.0.0.0:8002"); err != nil {
		fmt.Printf("startup service failed, err:%v\n\n", err)
	}
}
