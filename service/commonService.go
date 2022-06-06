package service

// 一些user、advisor都会用到的service接口 例如密码、手机号码重复校验等。
import (
	"database/sql"
	"fmt"
	qb "github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
)

var TABLES = []string{ADVISORTABLE, USERTABLE}

// CheckPhoneExist 检查手机号是否重复 true=已经存在 false=不存在
func CheckPhoneExist(tableName string, phone interface{}) int {
	// 生产sql语句
	where := map[string]interface{}{
		"phone": phone,
	}
	selectFields := []string{"phone"}
	cond, values, err := qb.BuildSelect(tableName, where, selectFields)
	if err != nil {
		logger.GendryBuildError("CheckPhoneExist", err)
		return errmsg.ErrorSqlBuild
	}
	// 查询
	rows, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		logger.SqlError("CheckPhoneExist", "select", err)
		return errmsg.ErrorMysql
	}
	// 判断是否存在重复key
	res, err := scanner.ScanMapDecodeClose(rows)
	if err != nil {
		if err == scanner.ErrNilRows {
			return errmsg.SUCCESS
		}
		logger.GendryScannerError("CheckPhoneExist", err)
		return errmsg.ErrorSqlScanner
	}
	if len(res) != 0 {
		return errmsg.ErrorUserPhoneUsed
	} else {
		return errmsg.SUCCESS
	}
}

// ChangePWD 更改用户密码
func ChangePWD(tableName string, id int64, newPwd string) int {
	// 密码加密
	newPwd = GetPwd(newPwd)
	// 构造sql
	where := map[string]interface{}{
		"id": id,
	}
	updates := map[string]interface{}{
		"password": newPwd,
	}
	// 构造sql 执行更新
	cond, values, err := qb.BuildUpdate(tableName, where, updates)
	if err != nil {
		logger.GendryBuildError("ChangePWD", err)
		return errmsg.ErrorSqlBuild
	}
	_, err = utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.SqlError("ChangePWD", "update", err)
		return errmsg.ErrorMysql
	}
	return errmsg.SUCCESS
}

// GetPwd 获取加密的密码
func GetPwd(pwd string) string {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if err != nil {
		logger.Log.Error("生成密码错误", zap.Error(err))
		return pwd
	}
	return string(hashPwd)
}

// checkPwd 检查用户输入的密码和数据库中加密的密码是否一致
// pwd:用户输入的密码 encryptPwd 数据库中加密的密码
func checkPwd(pwd string, encryptPwd string) int {
	err := bcrypt.CompareHashAndPassword([]byte(encryptPwd), []byte(pwd))
	if err != nil {
		return errmsg.ErrorPasswordWrong
	}
	return errmsg.SUCCESS
}

// CheckRolePwd 检查不同的角色对应的用户密码是否对应
// table:不同角色对应的表名 phone:手机号 pwd:密码
func CheckRolePwd(table string, id int64, pwd string) int {
	var encryptPwd string
	// 从数据库中查加密后的密码
	code, res := GetTableItem(table, id, "password")
	if code == errmsg.SUCCESS {
		encryptPwd = fmt.Sprintf("%s", res)
		// 查到了加密密码在比对
		return checkPwd(pwd, encryptPwd)
	}
	return code
}

// Update 更新用户或者顾问的信息
func Update(table string, Info map[string]interface{}) int {

	where := map[string]interface{}{
		"id": Info["id"],
	}
	// 构造sql 执行更新
	cond, values, err := qb.BuildUpdate(table, where, Info)
	if err != nil {
		logger.Log.Error("更新信息错误，编译SQL错误", zap.Error(err))
		return errmsg.ErrorSqlBuild
	}
	_, err = utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.SqlError("UpdateUser", "update", err)
		return errmsg.ErrorMysql
	}
	return errmsg.SUCCESS
}

// CheckIdExist 用于Token中检查Id和table是否存在
func CheckIdExist(id int64, table string) int {
	valid := false
	for _, v := range TABLES {
		if v == table {
			valid = true
			break
		}
	}
	if !valid {
		return errmsg.ErrorTokenRoleNotExist
	}
	code := CheckRolePwd(table, id, "")
	if code == errmsg.ErrorUserNotExist {
		return errmsg.ErrorTokenIdNotExist
	} else if code == errmsg.ErrorPasswordWrong {
		return errmsg.SUCCESS
	} else {
		return code
	}
}

// GetTableItem 通用的单项查询结构 字符串类型返回的是uint8
func GetTableItem(tableName string, tableId int64, fieldName string, tx ...*sql.Tx) (code int, res interface{}) {
	var err error
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Error("通用查询出错", zap.Error(err), zap.String("errorMsg", errmsg.GetErrMsg(code)))
		}
	}()
	where := map[string]interface{}{
		"id": tableId,
	}
	selects := []string{fieldName}
	cond, values, err := qb.BuildSelect(tableName, where, selects)
	if err != nil {
		logger.GendryBuildError("GetTableItem", err)
		return errmsg.ErrorSqlBuild, nil
	}
	var row *sql.Row
	// 可能是事务调用的
	if len(tx) != 0 {
		row = tx[0].QueryRow(cond, values...)
	} else {
		row = utils.DbConn.QueryRow(cond, values...)
	}
	if err = row.Scan(&res); err != nil {
		logger.Log.Error(fmt.Sprintf("无法从表%s根据%d匹配到%s字段,请检查。", tableName, tableId, fieldName), zap.Error(err))
		return errmsg.ErrorMysql, nil
	}
	return errmsg.SUCCESS, res
}

// GetManyTableItemsById 通过Id字段查从数据表中查多个字段
func GetManyTableItemsById(tableName string, tableId int64, selects []string, tx ...*sql.Tx) (code int, res map[string]interface{}) {
	var err error
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Error("通用查询出错", zap.Error(err), zap.String("errorMsg", errmsg.GetErrMsg(code)))
		}
	}()
	where := map[string]interface{}{
		"id": tableId,
	}
	cond, values, err := qb.BuildSelect(tableName, where, selects)
	if err != nil {
		logger.GendryBuildError("GetManyTableItemsById", err)
		return errmsg.ErrorSqlBuild, nil
	}
	var rows *sql.Rows
	if len(tx) != 0 {
		rows, err = tx[0].Query(cond, values...)
	} else {
		rows, err = utils.DbConn.Query(cond, values...)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return errmsg.ErrorNoResult, nil
		} else {
			return errmsg.ErrorMysql, nil
		}
	}
	results, err := scanner.ScanMapDecodeClose(rows)
	if err != nil {
		logger.GendryScannerError("GetManyTableItemsById", err, "cond", cond, "values", values)
		return errmsg.ErrorSqlScanner, nil
	}
	if len(results) > 1 {
		logger.Log.Error("主键查出多条数据来了", zap.Int64("id", tableId), zap.String("table", tableName))
		return errmsg.ErrorRowNotExpect, nil
	} else if len(results) == 0 {
		return errmsg.ErrorNoResult, nil
	}
	return errmsg.SUCCESS, results[0]

}

// GetManyTableItemsByWhere 通过条件判断从数据表中查多个字段
func GetManyTableItemsByWhere(tableName string, where map[string]interface{}, selects []string, tx ...*sql.Tx) (code int, res []map[string]interface{}) {
	var err error
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Error("通用查询出错", zap.Error(err), zap.String("errorMsg", errmsg.GetErrMsg(code)))
		}
	}()
	cond, values, err := qb.BuildSelect(tableName, where, selects)
	if err != nil {
		logger.GendryBuildError("GetManyTableItemsByWhere", err)
		return errmsg.ErrorSqlBuild, nil
	}
	var rows *sql.Rows
	if len(tx) != 0 {
		rows, err = tx[0].Query(cond, values...)
	} else {
		rows, err = utils.DbConn.Query(cond, values...)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return errmsg.ErrorNoResult, nil
		} else {
			return errmsg.ErrorMysql, nil
		}
	}
	results, err := scanner.ScanMapDecodeClose(rows)
	if err != nil {
		logger.GendryScannerError("GetManyTableItemsByWhere", err, "cond", cond, "values", values)
		return errmsg.ErrorSqlScanner, nil
	}
	return errmsg.SUCCESS, results
}

// UpdateTableItem 单项更新表值，传入表名，表id，map，tx为事务可选
func UpdateTableItem(tableName string, tableId int64, updates map[string]interface{}, tx ...*sql.Tx) (code int) {
	defer func() {
		fields := []zapcore.Field{
			zap.String("table", tableName),
			zap.Int64("table", tableId),
			zap.String("updates", fmt.Sprintf("%v", updates)),
		}
		if code == errmsg.SUCCESS {
			logger.Log.Info("更新表", fields...)
		} else {
			logger.Log.Error("更新表", fields...)
		}
	}()
	where := map[string]interface{}{
		"id": tableId,
	}
	cond, values, err := qb.BuildUpdate(tableName, where, updates)
	if err != nil {
		logger.GendryBuildError("UpdateTableItem", err)
		code = errmsg.ErrorSqlBuild
		return code
	}
	var res sql.Result
	if len(tx) != 0 {
		res, err = tx[0].Exec(cond, values...)
	} else {
		res, err = utils.DbConn.Exec(cond, values...)
	}
	if err != nil {
		logger.Log.Error("通用service接口更新参数值失败", zap.Error(err), zap.String("cond", cond), zap.String("values", fmt.Sprintf("%v", values)))
		code = errmsg.ErrorMysql
		return code
	}
	affectRow, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error("获取res.RowsAffected()失败", zap.Error(err))
		code = errmsg.ErrorMysql
		return code
	}
	if affectRow != 1 {
		logger.Log.Error(errmsg.GetErrMsg(errmsg.ErrorAffectsNotOne), zap.String("id", fmt.Sprintf("%v", tableId)), zap.String("table", tableName))
		code = errmsg.ErrorAffectsNotOne
		return code
	}
	return errmsg.SUCCESS
}
