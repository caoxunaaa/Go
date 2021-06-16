package UserRelation

import (
	"SuperxonWebSite/Databases"
	"errors"
)

type RightManager struct {
	Id        int64  `db:"id"`
	Username  string `db:"username"`
	Nickname  string `db:"nickname"`
	RightItem string `db:"right_item"`
}

func FindAllRightManager() ([]RightManager, error) {
	var res = make([]RightManager, 0)
	sqlStr := "SELECT * FROM right_manager"
	err := Databases.SuperxonUserDb.Select(&res, sqlStr)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FindAllRightManagerByUsername(username string) ([]RightManager, error) {
	var res = make([]RightManager, 0)
	sqlStr := "SELECT id, username, nickname, right_item FROM right_manager WHERE username=?"
	err := Databases.SuperxonUserDb.Select(&res, sqlStr, username)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func CreateRightManager(p *RightManager) error {
	sqlStr := "INSERT INTO right_manager(username, nickname, right_item) VALUES (?,?,?)"
	res, err := Databases.SuperxonUserDb.Exec(sqlStr, p.Username, p.Nickname, p.RightItem)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("权限管理中没有行被创建")
	}
	return nil
}

func UpdateRightManager(p *RightManager, id int64) error {
	sqlStr := "UPDATE right_manager SET username=?, nickname=?, right_item=? WHERE id=?"
	res, err := Databases.SuperxonUserDb.Exec(sqlStr, p.Username, p.Nickname, p.RightItem, id)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("权限管理中没有行被更新")
	}
	return nil
}

func DeleteRightManager(id int64) error {
	sqlStr := "DELETE FROM RightManager WHERE id=?"
	res, err := Databases.SuperxonUserDb.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("权限管理中没有行被删除")
	}
	return nil
}
