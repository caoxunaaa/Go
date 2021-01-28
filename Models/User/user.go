package User

import (
	"SuperxonWebSite/Databases"
	"errors"
)

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

func GetProfile(username string) (profile *Profile, err error) {
	sqlStr := "SELECT id, username, is_superuser FROM profiles where username=?"
	err = Databases.SuperxonDbDevice.Get(profile, sqlStr, username)
	if err != nil {
		return nil, err
	}
	return
}

type RegisterInfo struct {
	Username      string
	Password      string
	AgainPassword string
}

func Register(registerInfo *RegisterInfo) (err error) {
	if registerInfo.Password != registerInfo.AgainPassword {
		return errors.New("两次密码不一致，请重新确认")
	}

	sqlStr := "INSERT INTO profiles(username, password) values (?, ?)"
	_, err = Databases.SuperxonDbDevice.Exec(sqlStr, registerInfo.Username, registerInfo.Password)
	if err != nil {
		return errors.New("用户名已存在")
	}
	return
}
