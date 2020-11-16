package routers

import (
	"SuperxonWebSite/Controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", Controllers.IndexHandler)

	// v1
	v1Group := r.Group("v1")
	{
		// 待办事项
		// 添加
		v1Group.POST("/todo", Controllers.CreateAEquipment)
		// 查看所有的待办事项
		v1Group.GET("/todo", Controllers.GetEquipmentList)
		// 查看某一个待办事项
		v1Group.GET("/todo/:id", Controllers.GetAEquipment)
		// 修改某一个待办事项
		v1Group.PUT("/todo/:id", Controllers.UpdateAEquipment)
		// 删除某一个待办事项
		v1Group.DELETE("/todo/:id", Controllers.DeleteAEquipment)
	}
	return r
}
