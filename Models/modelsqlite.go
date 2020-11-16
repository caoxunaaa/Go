package Models

import (
	"SuperxonWebSite/Databases"
	"gorm.io/gorm"
)

type Equipment struct {
	gorm.Model
	Name string `json:"name"`
	Sort string `json:"sort"`
	Sn string `json:"sn"`
}

func CreateAEquipment(equipment *Equipment) (err error) {
	err = Databases.DB.Create(&equipment).Error
	return
}

func GetAEquipment(id string) (equipment *Equipment, err error) {
	equipment = new(Equipment)
	if err = Databases.DB.Where("id = ?", id).First(&equipment).Error; err!=nil{
		return nil, err
	}
	return
}

func GetAllEquipment() (equipments []*Equipment, err error) {
	if err = Databases.DB.Find(&equipments).Error; err != nil{
		return nil, err
	}
	return
}

func UpdateAEquipment(equipment *Equipment)(err error){
	err = Databases.DB.Where("id = ?", equipment.ID).First(&equipment).Error
	if err = Databases.DB.Model(&equipment).Updates(&equipment).Error; err != nil{
		return err
	}
	return
}

func DeleteAEquipment(id string)(err error){
	err = Databases.DB.Where("id = ?", id).Delete(&Equipment{}).Error
	return
}

