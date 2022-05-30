package service

// 一些user、advisor都会用到的service接口 例如密码、手机号码重复校验等。
import (
	"database/sql"
	qb "github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"go.uber.org/zap"
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
		return errmsg.ErrorUserphoneUsed
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
		return errmsg.ErrorPasswordWoron
	}
	return errmsg.SUCCESS
}

// CheckRolePwd 检查不同的角色对应的用户密码是否对应
// table:不同角色对应的表名 phone:手机号 pwd:密码
func CheckRolePwd(table string, id int64, pwd string) int {
	var encryptPwd string
	// 从数据库中查加密后的密码
	where := map[string]interface{}{
		"id": id,
	}
	selectFiled := []string{"password"}
	cond, value, err := qb.BuildSelect(table, where, selectFiled)
	if err != nil {
		logger.GendryBuildError("CheckRolePwd", err)
		return errmsg.ErrorSqlBuild
	}
	rows := utils.DbConn.QueryRow(cond, value...)
	// TODO
	err = rows.Scan(&encryptPwd)
	if err != nil {
		if err == sql.ErrNoRows {
			return errmsg.ErrorUserNotExist
		} else {
			logger.SqlError("CheckRolePwd", "select", err)
			return errmsg.ErrorPasswordWoron
		}
	}
	// 查到了加密密码在比对
	return checkPwd(pwd, encryptPwd)
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
	} else if code == errmsg.ErrorPasswordWoron {
		return errmsg.SUCCESS
	} else {
		return code
	}
}
