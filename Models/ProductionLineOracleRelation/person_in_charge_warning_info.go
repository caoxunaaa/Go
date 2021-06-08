package ProductionLineOracleRelation

import "SuperxonWebSite/Databases"

type PersonInChargeWarningInfo struct {
	Id        int64  `db:"id"`
	Username  string `db:"username"`
	Nickname  string `db:"nickname"`
	ModuleOsa string `db:"module_osa"`
	PnList    string `db:"pn_list"`
}

func FindAllPersonInChargeWarningInfo() ([]PersonInChargeWarningInfo, error) {
	var res = make([]PersonInChargeWarningInfo, 0)
	sqlStr := "SELECT id, username, nickname, module_osa, pn_list FROM person_in_charge_warning_info"
	err := Databases.SuperxonProductionLineOracleRelationDb.Select(&res, sqlStr)
	if err != nil {
		return nil, err
	}
	return res, nil
}
