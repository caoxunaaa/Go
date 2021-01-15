package FileManage

import (
	"SuperxonWebSite/Databases"
	"time"
)

type VideoInfo struct {
	ID         uint      `gorm:"primary_key" db:"id"`
	Uploader   string    `gorm:"default:'unknown'" db:"uploader"`
	Name       string    `gorm:"default:'unnamed'" db:"name"`
	StorePath  string    `db:"store_path"`
	UploadTime time.Time `db:"upload_time"`
}

func GetVideoInfoList() (videoInfoList []*VideoInfo, err error) {
	sqlStr := "SELECT * FROM video_infos"
	err = Databases.SuperxonDbDevice.Select(&videoInfoList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func CreateVideoInfo(videoInfo *VideoInfo) (err error) {
	sqlStr := "INSERT INTO video_infos(uploader, name, store_path) values (?, ?, ?)"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr,
		videoInfo.Uploader,
		videoInfo.Name,
		videoInfo.StorePath)
	if err != nil {
		return err
	}
	return
}

func DeleteVideoInfo(id uint) (length int64, err error) {
	sqlStr := "DELETE FROM video_infos where id = ?"
	res, err := Databases.SuperxonDbDevice.Exec(sqlStr, id)
	if err != nil {
		return length, err
	}
	length, err = res.RowsAffected()
	return
}
