package DeviceManage

import (
	"SuperxonWebSite/Databases"
	"database/sql"
	"errors"
)

type DeviceCategory struct {
	ID                 uint           `gorm:"primary_key" db:"id"`
	Name               string         `gorm:"unique;not null" db:"name"`
	ParentCategoryName sql.NullString `db:"parent_category_name"`
}

func GetAllDeviceCategoryRootList() (deviceCategoryRootNameList []*DeviceCategory, err error) {
	sqlStr := "select * from device_categories where parent_category_name is NUll"
	err = Databases.SuperxonDbDevice.Select(&deviceCategoryRootNameList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func GetAllDeviceCategoryChildList(rootCategory string) (deviceCategoryChildNameList []*DeviceCategory, err error) {
	sqlStr := "select * from device_categories where parent_category_name = ?"
	//sqlStr := "select Child.name from device_categories Child, (select id, name from device_categories where parent_category_name is NUll) As Root WHERE Child.parent_category_name = Root.name and Root.name = ?"
	err = Databases.SuperxonDbDevice.Select(&deviceCategoryChildNameList, sqlStr, rootCategory)
	if err != nil {
		return nil, err
	}
	return
}

func CreateDeviceCategoryChild(deviceCategory *DeviceCategory) (err error) {
	allDeviceCategoryRootList, _ := GetAllDeviceCategoryRootList()
	for index, v := range allDeviceCategoryRootList {
		if v.Name == deviceCategory.ParentCategoryName.String {
			break
		}
		if index == (len(allDeviceCategoryRootList) - 1) {
			return errors.New("没有相应的一级类型")
		}
	}
	sqlStr := "INSERT INTO device_categories(name, parent_category_name) values (?, ?)"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr, deviceCategory.Name, deviceCategory.ParentCategoryName.String)
	if err != nil {
		return err
	}
	return
}
