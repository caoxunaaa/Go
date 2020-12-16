package main

import (
	"SuperxonWebSite/Databases"
)

func main() {
	//Databases.InitOracle()
	//defer Databases.CloseOracle()

	Databases.RedisInit()
	defer Databases.RedisClose()

	//r := routers.Init()
	//if err := r.Run("0.0.0.0:8003"); err != nil {
	//	fmt.Printf("startup service failed, err:%v\n\n", err)
	//}
}
