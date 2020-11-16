package Controllers

import (
	"SuperxonWebSite/Models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}


func CreateAEquipment(c *gin.Context) {
	// 前端页面填写待办事项 点击提交 会发请求到这里
	// 1. 从请求中把数据拿出来
	var todo Models.Equipment
	err := c.BindJSON(&todo)
	// 2. 存入数据库
	err = Models.CreateAEquipment(&todo)
	if err != nil{
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}else{
		c.JSON(http.StatusOK, todo)
	}
}

func GetEquipmentList(c *gin.Context) {
	// 查询todo这个表里的所有数据
	todoList, err := Models.GetAllEquipment()
	if err!= nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}else {
		c.JSON(http.StatusOK, todoList)
	}
}

func GetAEquipment(c *gin.Context) {
	// 查询todo这个表里的所有数据
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
		return
	}
	todo, err := Models.GetAEquipment(id)
	if err!= nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}else {
		c.JSON(http.StatusOK, todo)
	}
}

func UpdateAEquipment(c *gin.Context) {
	var equipment Models.Equipment
	id, ok := c.Params.Get("id")
	err :=c.ShouldBindJSON(&equipment)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
		return
	}
	err = Models.UpdateAEquipment(&equipment)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}else{
		c.JSON(http.StatusOK, gin.H{id:"updated"})
	}
}

func DeleteAEquipment(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
		return
	}
	if err := Models.DeleteAEquipment(id);err!=nil{
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}else{
		c.JSON(http.StatusOK, gin.H{id:"deleted"})
	}
}