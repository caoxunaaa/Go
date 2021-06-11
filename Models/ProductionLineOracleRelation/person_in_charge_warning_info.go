package ProductionLineOracleRelation

import (
	"SuperxonWebSite/Databases"
	"errors"
)

type PersonInChargeWarningInfo struct {
	Id        int64  `db:"id"`
	Username  string `db:"username"`
	Nickname  string `db:"nickname"`
	ModuleOsa string `db:"module_osa"`
	Pn        string `db:"pn"`
}

func FindAllPersonInChargeWarningInfo() ([]PersonInChargeWarningInfo, error) {
	var res = make([]PersonInChargeWarningInfo, 0)
	sqlStr := "SELECT id, username, nickname, module_osa, pn FROM person_in_charge_warning_info"
	err := Databases.SuperxonProductionLineOracleRelationDb.Select(&res, sqlStr)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FindAllPersonInChargeWarningInfoByNickname(nickname string) ([]PersonInChargeWarningInfo, error) {
	var res = make([]PersonInChargeWarningInfo, 0)
	sqlStr := "SELECT id, username, nickname, module_osa, pn FROM person_in_charge_warning_info WHERE nickname=?"
	err := Databases.SuperxonProductionLineOracleRelationDb.Select(&res, sqlStr, nickname)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func CreatePersonInChargeWarningInfo(p *PersonInChargeWarningInfo) error {
	sqlStr := "INSERT INTO person_in_charge_warning_info(username, nickname, module_osa, pn) VALUES (?,?,?,?)"
	res, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(sqlStr, p.Username, p.Nickname, p.ModuleOsa, p.Pn)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("告警负责人没有行被创建")
	}
	return nil
}

func UpdatePersonInChargeWarningInfo(p *PersonInChargeWarningInfo, id int64) error {
	sqlStr := "UPDATE person_in_charge_warning_info SET username=?, nickname=?, module_osa=?, pn=? WHERE id=?"
	res, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(
		sqlStr,
		p.Username,
		p.Nickname,
		p.ModuleOsa,
		p.Pn,
		id)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("告警负责人没有行被更新")
	}
	return nil
}

func DeletePersonInChargeWarningInfo(id int64) error {
	sqlStr := "DELETE FROM person_in_charge_warning_info WHERE id=?"
	res, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("告警负责人没有行被删除")
	}
	return nil
}
