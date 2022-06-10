package service

// 一些user、advisor都会用到的service接口 例如密码、手机号码重复校验等。
import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
)

var ROLETABLES = []string{ADVISORTABLE, USERTABLE}

// CommonTranDefer 用在包含事务的service函数中，自动回滚事务
func CommonTranDefer(code *int, Tran *sql.Tx) {
	if *code != errmsg.SUCCESS {
		if err := Tran.Rollback(); err != nil {
			logger.Log.Error("事务回滚失败", zap.Error(err))
		}
	}
}

// CheckPhoneExist 检查手机号是否重复
func CheckPhoneExist(table string, phone interface{}) (code int) {
	// 生产sql语句
	code, _ = GetTableItemByWhere(
		table,
		map[string]interface{}{"phone": phone},
		"phone",
	)
	if code != errmsg.ErrorMysqlNoRows {
		code = errmsg.ErrorUserPhoneUsed
		return
	}
	return errmsg.SUCCESS
}

// ChangePWD 更改用户密码
func ChangePWD(tableName string, id int64, newPwd string) (code int) {
	// 密码加密
	newPwd = GetEncryptPwd(newPwd)

	updates := map[string]interface{}{
		"password": newPwd,
	}
	// 构造sql 执行更新
	if code = UpdateTableItemById(tableName, id, updates); code != errmsg.SUCCESS {
		return
	}
	return errmsg.SUCCESS
}

// GetEncryptPwd 获取加密的密码
func GetEncryptPwd(pwd string) string {
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
func CheckRolePwd(table string, id int64, pwd string) (code int) {
	var encryptPwd string
	var res interface{}
	// 从数据库中查加密后的密码
	if code, res = GetTableItemById(table, id, "password"); code != errmsg.SUCCESS {
		return
	}
	encryptPwd = fmt.Sprintf("%s", res)
	return checkPwd(pwd, encryptPwd)
}

// CheckIdExist 用于Token中检查Id和table是否存在
func CheckIdExist(id int64, table string) int {
	valid := false
	for _, v := range ROLETABLES {
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
