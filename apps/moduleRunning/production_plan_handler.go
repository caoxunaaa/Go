package moduleRunning

import (
	"SuperxonWebSite/Models/ModuleRunDisplay"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetProjectPlanListHandler(c *gin.Context) {
	projectPlanInfoList, err := ModuleRunDisplay.GetProjectPlanList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, projectPlanInfoList)
	}
}

func GetUndoneProjectPlanInfoListHandler(c *gin.Context) {
	undoneProjectPlanInfoList, err := ModuleRunDisplay.GetUndoneProjectPlanInfoList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, undoneProjectPlanInfoList)
	}
}

func CreateUndoneProjectPlanInfoHandler(c *gin.Context) {
	var undoneProjectPlanInfo ModuleRunDisplay.UndoneProjectPlanInfo
	if err := c.ShouldBindJSON(&undoneProjectPlanInfo); err == nil {
		err = ModuleRunDisplay.CreateUndoneProjectPlanInfo(&undoneProjectPlanInfo)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"Code": undoneProjectPlanInfo.Code,
				"Pn":   undoneProjectPlanInfo.Pn,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func DeleteUndoneProjectPlanInfoHandler(c *gin.Context) {
	idStr, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	id, _ := strconv.Atoi(idStr)
	err := ModuleRunDisplay.DeleteUndoneProjectPlanInfo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	}
}

func UpdateUndoneProjectPlanInfoHandler(c *gin.Context) {
	var undoneProjectPlanInfo ModuleRunDisplay.UndoneProjectPlanInfo
	idStr, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Sn"})
		return
	}
	err := c.ShouldBindJSON(&undoneProjectPlanInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	id, _ := strconv.Atoi(idStr)
	err = ModuleRunDisplay.UpdateUndoneProjectPlanInfo(&undoneProjectPlanInfo, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "数据已更新"})
	}
}
