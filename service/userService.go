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
func UpdateUser(user *model.User) int {

	where := map[string]interface{}{
		"name": user.Name,
	}
	// 把新的用户角色直接转化为map,去掉其中的value为空的key 和 username,password.
	// username,password,coin不可直接更新
	updates := structs.Map(user)
	delete(updates, "name")
	delete(updates, "password")
	delete(updates, "coin")
	for k, v := range updates {
		if v == "" {
			delete(updates, k)
		}
	}
	// 构造sql 执行更新
	cond, vals, err := qb.BuildUpdate(USERTABLE, where, updates)
	if err != nil {
		log.Println("gendry SQL生成错误", err)
		return errmsg.ERROR_SQL_BUILD
	}
	_, err = utils.DbConn.Exec(cond, vals...)
	if err != nil {
		log.Println("数据库更新数据出错", err)
		return errmsg.ERROR_MYSQL
	}
	return errmsg.SUCCESS
}

// ChangeUserPWD 更改用户密码
func ChangeUserPWD() {

}

// GetUser 获取用户的全部信息，测试用.
func GetUser() {

}
