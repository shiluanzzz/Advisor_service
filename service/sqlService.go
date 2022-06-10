package service

import (
	"database/sql"
	"fmt"
	qb "github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"go.uber.org/zap"
	"service-backend/utils"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
)

// GetTableItemById 通用的单项查询结构 字符串类型返回的是uint8
func GetTableItemById(tableName string, tableId int64, fieldName string, tx ...*sql.Tx) (code int, res interface{}) {
	where := map[string]interface{}{
		"id": tableId,
	}
	return GetTableItemByWhere(tableName, where, fieldName, tx...)
}

// getTableRows 获取数据库中的某一行，可以绑定到model上去
func getTableRows(tableName string, where map[string]interface{}, fieldName ...string) (code int, res *sql.Rows) {
	var err error
	cond, values, err := qb.BuildSelect(tableName, where, fieldName)
	if err != nil {
		logger.GendryBuildError(err)
		return errmsg.ErrorSqlBuild, nil
	}
	res, err = utils.DbConn.Query(cond, values...)
	if err != nil {
		code = errmsg.ErrorMysql
		return
	}
	return errmsg.SUCCESS, res
}

// GetTableItemByWhere 通用的单项查询结构
func GetTableItemByWhere(tableName string, where map[string]interface{}, fieldName string, tx ...*sql.Tx) (code int, res interface{}) {
	var err error
	msg := fmt.Sprintf("从表 [%s] 根据 [%v] 匹配 [%s] 字段", tableName, where, fieldName)
	defer logger.CommonServiceLog(&code, msg)
	selects := []string{fieldName}
	cond, values, err := qb.BuildSelect(tableName, where, selects)
	if err != nil {
		logger.GendryBuildError(err)
		return errmsg.ErrorSqlBuild, nil
	}

	var results []map[string]interface{}
	if code, results = SQLQuery(cond, values, tx...); code != errmsg.SUCCESS {
		return
	}
	if len(results) == 0 {
		return errmsg.ErrorNoResult, nil
	}
	// TODO 测试
	return errmsg.SUCCESS, results[0][fieldName]
}

// GetTableItemsByWhere 通过条件判断从数据表中查多个字段
func GetTableItemsByWhere(tableName string, where map[string]interface{}, selects []string, tx ...*sql.Tx) (code int, res []map[string]interface{}) {
	var err error
	defer func() {
		msg := fmt.Sprintf("从表 [%s] 通过条件 [%v] 查询 [%v]数据", tableName, where, selects)
		logger.CommonServiceLog(&code, msg)
	}()
	cond, values, err := qb.BuildSelect(tableName, where, selects)
	if err != nil {
		logger.GendryBuildError(err)
		return errmsg.ErrorSqlBuild, nil
	}
	return SQLQuery(cond, values, tx...)
}

// UpDateTableItemByWhere 通过条件更新表的字段
func UpDateTableItemByWhere(tableName string, where map[string]interface{}, updates map[string]interface{}, tx ...*sql.Tx) (code int) {

	defer func() {
		msg := fmt.Sprintf("从表 [%s] 根据 [%v] 更新数据 [%v] ", tableName, where, updates)
		logger.CommonServiceLog(&code, msg)
	}()
	cond, values, err := qb.BuildUpdate(tableName, where, updates)
	if err != nil {
		logger.GendryBuildError(err, "table", tableName, "where", where)
		code = errmsg.ErrorSqlBuild
		return code
	}
	code, _, _ = SQLExec(cond, values, tx...)
	return errmsg.SUCCESS
}

// UpdateTableItemById 单项更新表值，传入表名，表id，map，tx为事务可选
func UpdateTableItemById(tableName string, tableId int64, updates map[string]interface{}, tx ...*sql.Tx) (code int) {
	where := map[string]interface{}{
		"id": tableId,
	}
	return UpDateTableItemByWhere(tableName, where, updates, tx...)
}

func InsertTableItem(tableName string, data []map[string]interface{}, tx ...*sql.Tx) (code int, Id int64) {

	defer func() {
		msg := fmt.Sprintf("向表 [%s] 插入数据 [%v]", tableName, tableName)
		logger.CommonServiceLog(&code, nil, "msg", msg)
	}()

	cond, values, err := qb.BuildInsert(tableName, data)
	if err != nil {
		logger.GendryBuildError(err, "table", tableName, "data", data)
		return errmsg.ErrorMysql, -1
	}
	code, Id, _ = SQLExec(cond, values, tx...)
	return
}

// GetTableItemsById 通过Id字段查从数据表中查多个字段
func GetTableItemsById(tableName string, tableId int64, selects []string, tx ...*sql.Tx) (code int, res map[string]interface{}) {
	var err error
	defer func() {
		msg := fmt.Sprintf("从表 [%s] 中根据主键 [%d] 查询 [%v]数据", tableName, tableId, selects)
		logger.CommonServiceLog(&code, msg)
	}()
	where := map[string]interface{}{
		"id": tableId,
	}
	cond, values, err := qb.BuildSelect(tableName, where, selects)
	if err != nil {
		logger.GendryBuildError(err)
		return errmsg.ErrorSqlBuild, nil
	}

	var results []map[string]interface{}
	if code, results = SQLQuery(cond, values, tx...); code != errmsg.SUCCESS {
		return
	}

	// 主键就一条数据
	if len(results) > 1 {
		logger.Log.Error("主键查出多条数据来了", zap.Int64("id", tableId), zap.String("table", tableName))
		return errmsg.ErrorRowNotExpect, nil
	} else if len(results) == 0 {
		return errmsg.ErrorNoResult, nil
	}
	return errmsg.SUCCESS, results[0]

}

func GetTableRows2StructByWhere(tableName string, where map[string]interface{}, selects []string, object interface{}) (code int) {
	defer func() {
		msg := fmt.Sprintf("从表 [%s] 根据 [%v] 查询数据绑定到结构体对象", tableName, where)
		logger.CommonServiceLog(&code, msg)
	}()

	var rows *sql.Rows
	if code, rows = getTableRows(tableName, where, selects...); code != errmsg.SUCCESS {
		return
	}
	err := scanner.Scan(rows, object)
	if err != nil {
		if err == scanner.ErrNilRows || err == scanner.ErrEmptyResult {
			return errmsg.ErrorNoResult
		}
		return errmsg.ErrorSqlScanner
	}
	return errmsg.SUCCESS
}
func DeleteTableRowByWhere(tableName string, where map[string]interface{}, tx ...*sql.Tx) (code int) {
	defer func() {
		msg := fmt.Sprintf("从表 [%s] 根据 [%v] 删除数据", tableName, where)
		logger.CommonServiceLog(&code, msg)
	}()
	cond, values, err := qb.BuildDelete(tableName, where)
	if err != nil {
		logger.GendryBuildError(err)
	}
	code, _, affects := SQLExec(cond, values, tx...)
	if code != errmsg.SUCCESS {
		return code
	}
	if affects == 0 {
		return errmsg.ErrorNoResult
	}
	return errmsg.SUCCESS
}

// SQLExec 执行SQL语句
func SQLExec(cond string, values []interface{}, tx ...*sql.Tx) (code int, Id int64, Affected int64) {
	var err error
	var row sql.Result
	defer logger.CommonServiceLog(&code, code, "values", fmt.Sprintf("%v", values))

	if len(tx) != 0 {
		row, err = tx[0].Exec(cond, values...)
	} else {
		row, err = utils.DbConn.Exec(cond, values...)
	}
	if err != nil {
		logger.SqlError(err, "cond", cond, "values", values)
		return errmsg.ErrorMysql, -1, -1
	}
	// 获取用户的主键ID
	if Id, err = row.LastInsertId(); err != nil {
		logger.Log.Error("获取数据库主键错误", zap.Error(err))
		return errmsg.ErrorMysql, -1, -1
	}
	if Affected, err = row.RowsAffected(); err != nil {
		logger.Log.Error("获取数据库affects错误", zap.Error(err))
		return errmsg.ErrorMysql, -1, -1
	}

	return errmsg.SUCCESS, Id, Affected
}

// SQLQuery 执行查询语句
func SQLQuery(cond string, values []interface{}, tx ...*sql.Tx) (code int, res []map[string]interface{}) {
	var rows *sql.Rows
	var err error
	defer logger.CommonServiceLog(&code, code, "values", fmt.Sprintf("%v", values))

	if len(tx) != 0 {
		rows, err = tx[0].Query(cond, values...)
	} else {
		rows, err = utils.DbConn.Query(cond, values...)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return errmsg.ErrorNoResult, nil
		} else {
			logger.SqlError(err, "cond", cond, "values", fmt.Sprintf("%v", values))
			return errmsg.ErrorMysql, nil
		}
	}
	if res, err = scanner.ScanMapDecodeClose(rows); err != nil {
		logger.GendryScannerError(err, "cond", cond, "values", fmt.Sprintf("%v", values))
		return errmsg.ErrorSqlScanner, nil
	}
	return errmsg.SUCCESS, res
}
