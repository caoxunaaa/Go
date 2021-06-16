package ProductionLineOracleRelation

import (
	"SuperxonWebSite/Databases"
	"errors"
	"fmt"
	"time"
)

//生产参数变更关联数据
type ProductionParameterChangedRelation struct {
	Id              int64  `db:"id"`
	MonitoringTable string `db:"monitoring_table"`
	OnlyField1Name  string `db:"only_field_1_name"`
	OnlyField2Name  string `db:"only_field_2_name"`
	OnlyField3Name  string `db:"only_field_3_name"`
	OnlyField4Name  string `db:"only_field_4_name"`
}

// 获取所有监控的表tableName的唯一索引
func FindAllProductionParameterChangedRelation() ([]ProductionParameterChangedRelation, error) {
	var res = make([]ProductionParameterChangedRelation, 0)
	sqlStr := "SELECT * FROM production_parameter_changed_relation"
	err := Databases.SuperxonProductionLineOracleRelationDb.Select(&res, sqlStr)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取需要监控的表tableName的唯一索引
func FindProductionParameterChangedRelationByTableName(tableName string) (*ProductionParameterChangedRelation, error) {
	var res = new(ProductionParameterChangedRelation)
	sqlStr := "SELECT * FROM production_parameter_changed_relation WHERE monitoring_table=?"
	err := Databases.SuperxonProductionLineOracleRelationDb.Get(res, sqlStr, tableName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//生产参数变更数据
type ProductionParameterChanged struct {
	Id              int64     `db:"id"`
	Username        string    `db:"username"`
	Nickname        string    `db:"nickname"`
	MonitoringTable string    `db:"monitoring_table"`
	OnlyField1Name  string    `db:"only_field_1_name"`
	OnlyField1Value string    `db:"only_field_1_value"`
	OnlyField2Name  string    `db:"only_field_2_name"`
	OnlyField2Value string    `db:"only_field_2_value"`
	OnlyField3Name  string    `db:"only_field_3_name"`
	OnlyField3Value string    `db:"only_field_3_value"`
	OnlyField4Name  string    `db:"only_field_4_name"`
	OnlyField4Value string    `db:"only_field_4_value"`
	ChangedItem     string    `db:"changed_item"`
	OldValue        string    `db:"old_value"`
	NewValue        string    `db:"new_value"`
	UpdateTime      time.Time `db:"update_time"`
	UpdateReason    string    `db:"update_reason"`
}

func FindAllProductionParameterChanged() ([]ProductionParameterChanged, error) {
	var res = make([]ProductionParameterChanged, 0)
	sqlStr := "SELECT * FROM production_parameter_changed"
	err := Databases.SuperxonProductionLineOracleRelationDb.Select(&res, sqlStr)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FindAllProductionParameterChangedByMonitoringTable(monitoringTable string) ([]ProductionParameterChanged, error) {
	var res = make([]ProductionParameterChanged, 0)
	sqlStr := "SELECT * FROM production_parameter_changed WHERE monitoring_table=?"
	err := Databases.SuperxonProductionLineOracleRelationDb.Select(&res, sqlStr, monitoringTable)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FindAllProductionParameterChangedByMonitoringTableAndOnlyField(monitoringTable, only_field_1_value, only_field_2_value, only_field_3_value, only_field_4_value string) ([]ProductionParameterChanged, error) {
	var res = make([]ProductionParameterChanged, 0)
	sqlStr := "SELECT * FROM production_parameter_changed WHERE monitoring_table=? AND only_field_1_value=? AND only_field_2_value=? AND only_field_3_value=? AND only_field_4_value=?"
	err := Databases.SuperxonProductionLineOracleRelationDb.Select(&res, sqlStr, monitoringTable, only_field_1_value, only_field_2_value, only_field_3_value, only_field_4_value)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FindAllProductionParameterChangedByMonitoringTableAndOnlyFieldAndChangedItem(monitoringTable, only_field_1_value, only_field_2_value, only_field_3_value, only_field_4_value, changed_item string) ([]ProductionParameterChanged, error) {
	var res = make([]ProductionParameterChanged, 0)
	sqlStr := "SELECT * FROM production_parameter_changed WHERE monitoring_table=? AND only_field_1_value=? AND only_field_2_value=? AND only_field_3_value=? AND only_field_4_value=? AND changed_item=?"
	err := Databases.SuperxonProductionLineOracleRelationDb.Select(&res, sqlStr, monitoringTable, only_field_1_value, only_field_2_value, only_field_3_value, only_field_4_value, changed_item)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func CreateProductionParameterChanged(p *ProductionParameterChanged) error {
	sqlStr := "INSERT INTO production_parameter_changed(username, nickname, monitoring_table, only_field_1_name, only_field_1_value, " +
		"only_field_2_name, only_field_2_value, only_field_3_name, only_field_3_value, only_field_4_name, only_field_4_value, changed_item, " +
		"old_value, new_value, update_time, update_reason) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	res, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(sqlStr, p.Username, p.Nickname, p.MonitoringTable,
		p.OnlyField1Name, p.OnlyField1Value, p.OnlyField2Name, p.OnlyField2Value, p.OnlyField3Name, p.OnlyField3Value, p.OnlyField4Name, p.OnlyField4Value,
		p.ChangedItem, p.OldValue, p.NewValue, p.UpdateTime, p.UpdateReason)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("生产工艺记录没有行被创建")
	}
	return nil
}

func UpdateProductionParameterChanged(p *ProductionParameterChanged, id int64) error {
	sqlStr := "UPDATE production_parameter_changed SET username=?, nickname=?, monitoring_table=?, only_field_1_name=?, only_field_1_value=?, " +
		"only_field_2_name=?, only_field_2_value=?, only_field_3_name=?, only_field_3_value=?, only_field_4_name=?, only_field_4_value=?, changed_item=?, " +
		"old_value=?, new_value=?, update_time=?, update_reason=? WHERE id=?"
	res, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(sqlStr, p.Username, p.Nickname, p.MonitoringTable,
		p.OnlyField1Name, p.OnlyField1Value, p.OnlyField2Name, p.OnlyField2Value, p.OnlyField3Name, p.OnlyField3Value, p.OnlyField4Name, p.OnlyField4Value,
		p.ChangedItem, p.OldValue, p.NewValue, p.UpdateTime, p.UpdateReason, id)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("生产工艺记录没有行被更新")
	}
	return nil
}

func DeleteProductionParameterChanged(id int64) error {
	sqlStr := "DELETE FROM production_parameter_changed WHERE id=?"
	res, err := Databases.SuperxonProductionLineOracleRelationDb.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("生产工艺记录没有行被删除")
	}
	return nil
}

//Oracle产线具体生产参数查询与更新
//获取Oracle监控表中的所有字段
func GetAllFieldByMonitoringTable(monitoringTable string) ([]string, error) {
	sqlStr := "SELECT * FROM " + monitoringTable + " WHERE ROWNUM <= 1"
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	res, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetParameterByMonitoringTableAndOnlyFieldAndChangedItem(p *ProductionParameterChanged) (string, error) {
	var sqlStr string
	var err error
	// 读ORACLE 该唯一行更改行的值
	if p.OnlyField4Name == "" {
		sqlStr = fmt.Sprint("SELECT "+p.ChangedItem+" FROM "+p.MonitoringTable+" WHERE ", p.OnlyField1Name, "='", p.OnlyField1Value, "' AND ", p.OnlyField2Name, "='", p.OnlyField2Value, "' AND ", p.OnlyField3Name, "='", p.OnlyField3Value, "'")
	} else {
		sqlStr = fmt.Sprint("SELECT "+p.ChangedItem+" FROM "+p.MonitoringTable+" WHERE ", p.OnlyField1Name, "='", p.OnlyField1Value, "' AND ", p.OnlyField2Name, "='", p.OnlyField2Value, "' AND ", p.OnlyField3Name, "='", p.OnlyField3Value, "' AND ", p.OnlyField4Name, "='", p.OnlyField4Value, "'")
	}
	fmt.Println(sqlStr)
	rows, err := Databases.OracleDB.Query(sqlStr)
	if err != nil {
		return "", err
	}
	var oldChangedItemValue = new(string)
	for rows.Next() {
		err = rows.Scan(oldChangedItemValue)
		if err != nil {
			return "", err
		}
	}
	return *oldChangedItemValue, nil
}

//修改ORACLE监控表中的某个字段并在生产工艺变更记录中修改动作
func UpdateParameterByMonitoringTableAndField(p *ProductionParameterChanged) error {
	var sqlStr string
	var err error

	// /修改ORACLE监控表中的某个字段的值
	if p.OnlyField4Name == "" {
		sqlStr = fmt.Sprint("UPDATE ", p.MonitoringTable, " SET ", p.ChangedItem, "='", p.NewValue, "' WHERE ", p.OnlyField1Name, "='", p.OnlyField1Value, "' AND ", p.OnlyField2Name, "='", p.OnlyField2Value, "' AND ", p.OnlyField3Name, "='", p.OnlyField3Value, "'")
	} else {
		sqlStr = fmt.Sprint("UPDATE ", p.MonitoringTable, " SET ", p.ChangedItem, "='", p.NewValue, "' WHERE ", p.OnlyField1Name, "='", p.OnlyField1Value, "' AND ", p.OnlyField2Name, "='", p.OnlyField2Value, "' AND ", p.OnlyField3Name, "='", p.OnlyField3Value, "' AND ", p.OnlyField4Name, "='", p.OnlyField4Value, "'")
	}
	fmt.Println(sqlStr)
	res, err := Databases.OracleDB.Exec(sqlStr)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("ORACLE 中的表" + p.MonitoringTable + "没有被更新")
	}

	// 回读ORACLE 该唯一行是否已经更改
	newValue, err := GetParameterByMonitoringTableAndOnlyFieldAndChangedItem(p)
	if err != nil {
		return err
	}
	if newValue != p.NewValue {
		return errors.New("回读oracle中该字段没有发生改变")
	}

	//在生产工艺变更记录中修改动作
	err = CreateProductionParameterChanged(p)
	if err != nil {
		return err
	}
	return nil
}
