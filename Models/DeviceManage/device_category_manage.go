package DeviceManage

import (
	"SuperxonWebSite/Databases"
	"database/sql"
	"errors"
)

type DeviceCategory struct {
	ID                 uint           `gorm:"primary_key" db:"id"`
	Name               string         `gorm:"unique;not null" db:"category_root"`
	ParentCategoryName sql.NullString `db:"parent_category_id"`
}

func GetAllDeviceCategoryRootList() (deviceCategoryRootNameList []string, err error) {
	sqlStr := "select name from device_categories where parent_category_name is NUll"
	err = Databases.SuperxonDbDevice.Select(&deviceCategoryRootNameList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func GetAllDeviceCategoryChildList(rootCategory string) (deviceCategoryChildNameList []string, err error) {
	sqlStr := "select Child.name from device_categories Child, (select id, name from device_categories where parent_category_name is NUll) As Root WHERE Child.parent_category_name = Root.name and Root.name = ?"
	err = Databases.SuperxonDbDevice.Select(&deviceCategoryChildNameList, sqlStr, rootCategory)
	if err != nil {
		return nil, err
	}
	return
}

func CreateDeviceCategoryChild(rootCategory string, childCategory string) (err error) {
	allDeviceCategoryRootList, _ := GetAllDeviceCategoryRootList()
	for index, v := range allDeviceCategoryRootList {
		if v == rootCategory {
			break
		}
		if index == (len(allDeviceCategoryRootList) - 1) {
			return errors.New("没有相应的一级类型")
		}
	}
	sqlStr := "INSERT INTO device_categories(name, parent_category_name) values (?, ?)"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr, childCategory, rootCategory)
	if err != nil {
		return err
	}
	return
}
