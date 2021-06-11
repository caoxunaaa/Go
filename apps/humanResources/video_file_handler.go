package humanResources

import (
	"SuperxonWebSite/Models/FileManager"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetVideoInfoListHandler(c *gin.Context) {
	videoInfoList, err := FileManager.GetVideoInfoList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, videoInfoList)
	}
}

func UploadVideoFileHandler(c *gin.Context) {
	var videoInfo *FileManager.VideoInfo
	videoInfo = new(FileManager.VideoInfo)

	videoInfo.Uploader = c.PostForm("uploader")
	videoFile, err := c.FormFile("videoFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	videoInfo.Name = videoFile.Filename
	fmt.Println("fileName", videoInfo.Name)

	nameSplit := strings.Split(videoInfo.Name, ".")
	dir := nameSplit[len(nameSplit)-1]

	videoInfo.StorePath = "/assets/" + dir + "/" + videoInfo.Name

	_, err = os.Stat("./assets")
	if os.IsNotExist(err) {
		fmt.Println("目录不存在,创建目录")
		err = os.Mkdir("./assets", 0777)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	_, err = os.Stat("./assets/" + dir)
	if os.IsNotExist(err) {
		fmt.Println("文件不存在,创建目录")
		err = os.Mkdir("./assets/"+dir, 0777)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	err = c.SaveUploadedFile(videoFile, "./assets/"+dir+"/"+videoInfo.Name)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = FileManager.CreateVideoInfo(videoInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": videoInfo.Name + "已经成功上传"})
	}
}

func DeleteVideoInfoHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
		return
	}
	idInt, _ := strconv.Atoi(id)
	length, err := FileManager.DeleteVideoInfo(uint(idInt))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": strconv.FormatInt(length, 10) + "行已经被删除"})
	}
}
