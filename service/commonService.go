package service

// 一些user、advisor都会用到的service接口 例如密码、手机号码重复校验等。
import (
	"database/sql"
	qb "github.com/didi/gendry/builder"
	"golang.org/x/crypto/bcrypt"
	"log"
	"service/utils"
	"service/utils/errmsg"
)

// CheckPhoneExist 检查手机号是否重复 true=已经存在 false=不存在
func CheckPhoneExist(tableName, phone string) int {
	// 生产sql语句
	where := map[string]interface{}{
		"phone": phone,
	}
	selectFields := []string{"phone"}
	cond, values, err := qb.BuildSelect(tableName, where, selectFields)
	if err != nil {
		log.Println("gendry SQL生成错误", err)
		return errmsg.ERROR_SQL_BUILD
	}
	// 查询
	rows, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		log.Println("数据库查询错误", err)
		return errmsg.ERROR_MYSQL
	}
	// 判断是否存在重复key
	var flag = false
	for rows.Next() {
		flag = true
		break
	}
	if flag {
		return errmsg.ERROR_USERPHONE_USED
	} else {
		return errmsg.SUCCESS
	}
}

// GetPwd 获取加密的密码
func GetPwd(pwd string) string {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if err != nil {
		log.Println("generate password error,", err)
		return pwd
	}
	return string(hashPwd)
}

// checkPwd 检查用户输入的密码和数据库中加密的密码是否一致
// checkPwd pwd:用户输入的密码 encryptPwd 数据库中加密的密码
func checkPwd(pwd string, encryptPwd string) int {
	err := bcrypt.CompareHashAndPassword([]byte(encryptPwd), []byte(pwd))
	if err != nil {
		return errmsg.ERROR_PASSWORD_WORON
	}
	return errmsg.SUCCESS
}

// CheckRolePwd 检查不同的角色对应的用户密码是否对应
// table:不同角色对应的表名 username:用户名 pwd:密码
func CheckRolePwd(table, username string, pwd string) int {
	var encryptPwd string
	// 从数据库中查加密后的密码
	where := map[string]interface{}{
		"name": username,
	}
	selectFiled := []string{"password"}
	cond, value, err := qb.BuildSelect(table, where, selectFiled)
	if err != nil {
		log.Println("gendry SQL生成错误", err)
		return errmsg.ERROR_SQL_BUILD
	}
	rows := utils.DbConn.QueryRow(cond, value...)
	err = rows.Scan(&encryptPwd)
	if err != nil {
		if err == sql.ErrNoRows {
			return errmsg.ERROR_USER_NOT_EXIST
		} else {
			log.Println("数据查询用户密码错误", err)
			return errmsg.ERROR_PASSWORD_WORON
		}
	}
	// 查到了加密密码在比对
	return checkPwd(pwd, encryptPwd)
}
