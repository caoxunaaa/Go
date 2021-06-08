package ProductionLineOracleRelation

import "SuperxonWebSite/Databases"

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
	_, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(sqlStr, p.Username, p.Nickname, p.ModuleOsa, p.Pn)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePersonInChargeWarningInfo(p *PersonInChargeWarningInfo, id int64) error {
	sqlStr := "UPDATE person_in_charge_warning_info SET username=?, nickname=?, module_osa=?, pn=? WHERE id=?"
	_, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(
		sqlStr,
		p.Username,
		p.Nickname,
		p.ModuleOsa,
		p.Pn,
		id)
	if err != nil {
		return err
	}
	return nil
}

func DeletePersonInChargeWarningInfo(id int64) error {
	sqlStr := "DELETE FROM person_in_charge_warning_info WHERE id=?"
	_, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	return nil
}
