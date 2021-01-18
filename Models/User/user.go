package User

import "SuperxonWebSite/Databases"

type Profile struct {
	ID          uint   `gorm:"primary_key" db:"id"`
	Username    string `db:"username"`
	Password    string `db:"password"`
	IsSuperuser bool   `gorm:"default:'0'" db:"is_superuser"`
}

func GetAllProfileList() (profileList []*Profile, err error) {
	sqlStr := "SELECT id, username, is_superuser FROM profiles"
	err = Databases.SuperxonDbDevice.Select(&profileList, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}
