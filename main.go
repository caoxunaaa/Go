package main

import (
	"SuperxonWebSite/Databases"
	routers "SuperxonWebSite/Router"
	"fmt"
)

func main() {
	Databases.InitOracle()
	defer Databases.CloseOracle()

	r := routers.Init()
	if err := r.Run("0.0.0.0:8002"); err != nil {
		fmt.Printf("startup service failed, err:%v\n\n", err)
	}
}
