// @Title  setting_warning_threshold.go
// @Description  对良率告警红黄线的CRUD
// @Author  曹迅 (时间 2021/06/02  10:00)
// @Update  曹迅 (时间 2021/06/02  10:00)
package ProductionLineOracleRelation

import (
	"SuperxonWebSite/Databases"
	"errors"
)

type SettingWarningThreshold struct {
	Id         int64  `db:"id"`
	OrderType  string `db:"order_type"`
	ModuleOsa  string `db:"module_osa"`
	Pn         string `db:"pn"`
	Process    string `db:"process"`
	YellowLine int    `db:"yellow_line"`
	RedLine    int    `db:"red_line"`
}

func FindOneSettingWarningThresholdById(id int64) (SettingWarningThreshold, error) {
	res := SettingWarningThreshold{}
	sqlStr := "SELECT id, order_type, module_osa, pn, process, yellow_line, red_line FROM setting_warning_threshold WHERE id=?"
	err := Databases.SuperxonProductionLineOracleRelationDb.Get(&res, sqlStr, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func FindDefaultSettingWarningThresholdByOrderTypeAndModuleOsa(orderType, moduleOsa string) (SettingWarningThreshold, error) {
	res := SettingWarningThreshold{}
	sqlStr := "SELECT id, order_type, module_osa, pn, process, yellow_line, red_line FROM setting_warning_threshold WHERE pn='默认' AND process='默认' AND order_type=? AND module_osa=?"
	err := Databases.SuperxonProductionLineOracleRelationDb.Get(&res, sqlStr, orderType, moduleOsa)
	if err != nil {
		return res, err
	}
	return res, nil
}

func FindSomeSettingWarningThresholdByOrderTypeAndModuleOsa(orderType, moduleOsa string) ([]SettingWarningThreshold, error) {
	var res = make([]SettingWarningThreshold, 0)
	sqlStr := "SELECT id, order_type, module_osa, pn, process, yellow_line, red_line FROM setting_warning_threshold WHERE order_type=? AND module_osa=?"
	err := Databases.SuperxonProductionLineOracleRelationDb.Select(&res, sqlStr, orderType, moduleOsa)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FindAllSettingWarningThreshold() ([]SettingWarningThreshold, error) {
	var res = make([]SettingWarningThreshold, 0)
	sqlStr := "SELECT id, order_type, module_osa, pn, process, yellow_line, red_line FROM setting_warning_threshold"
	err := Databases.SuperxonProductionLineOracleRelationDb.Select(&res, sqlStr)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func CreateSettingWarningThreshold(swt *SettingWarningThreshold) error {
	sqlStr := "INSERT INTO setting_warning_threshold(order_type, module_osa, pn, process, yellow_line, red_line) VALUES (?,?,?,?,?,?)"
	res, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(sqlStr, swt.OrderType, swt.ModuleOsa, swt.Pn, swt.Process, swt.YellowLine, swt.RedLine)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("告警门限设置没有行被创建")
	}
	return nil
}

func UpdateSettingWarningThreshold(swt *SettingWarningThreshold, id int64) error {
	sqlStr := "UPDATE setting_warning_threshold SET order_type=?, module_osa=?, pn=?, process=?, yellow_line=?, red_line=? WHERE id=?"
	res, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(
		sqlStr,
		swt.OrderType,
		swt.ModuleOsa,
		swt.Pn,
		swt.Process,
		swt.YellowLine,
		swt.RedLine,
		id)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("告警门限设置没有行被更新")
	}
	return nil
}

func DeleteSettingWarningThreshold(id int64) error {
	sqlStr := "DELETE FROM setting_warning_threshold WHERE id=?"
	res, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("告警门限设置没有行被删除")
	}
	return nil
}
