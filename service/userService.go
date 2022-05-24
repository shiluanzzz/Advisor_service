package service

import (
	qb "github.com/didi/gendry/builder"
	"github.com/fatih/structs"
	"log"
	"service/model"
	"service/utils"
	"service/utils/errmsg"
)

var USERTABLE = "user"

// NewUser 新增用户
func NewUser(user *model.User) int {
	// 转化数据并生成sql语句
	var data []map[string]interface{}
	data = append(data, structs.Map(user))
	cond, vals, err := qb.BuildInsert(USERTABLE, data)
	if err != nil {
		log.Println("gendry SQL生成错误", err)
		return errmsg.ERROR_SQL_BUILD
	}

	// 执行sql语句
	_, err = utils.DbConn.Exec(cond, vals...)
	if err != nil {
		log.Println("数据新增错误", err)
		return errmsg.ERROR_SQL_BUILD
	}
	return errmsg.SUCCESS
}

// 数据校验

// CheckUserPWD 检查用户名密码是否匹配
func CheckUserPWD() {

}

// CheckUserName 检查用户名是否重复 true=已经存在 false=不存在
func CheckUserName(name string) int {
	// 生产sql语句
	where := map[string]interface{}{
		"name": name,
	}
	selectFields := []string{"name"}
	cond, values, err := qb.BuildSelect(USERTABLE, where, selectFields)
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
		return errmsg.ERROR_USERNAME_USED
	} else {
		return errmsg.SUCCESS
	}
}

// 更新数据

// UpdateUser 更新用户的信息
func UpdateUser() {

}

// ChangeUserPWD 更改用户密码
func ChangeUserPWD() {

}

// GetUser 获取用户的全部信息，测试用.
func GetUser() {

}
