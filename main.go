package main

import (
	"SuperxonWebSite/Databases"
	"SuperxonWebSite/Models"
	routers "SuperxonWebSite/Router"
	"fmt"
)

func main(){
	err := Databases.InitSqlite()
	if err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer Databases.CloseSqlite() // 程序退出关闭数据库连接
	// 模型绑定
	
	Databases.DB.AutoMigrate(&Models.Equipment{})
	// 注册路由
	r := routers.SetupRouter()
	if err := r.Run(); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
	}
}
