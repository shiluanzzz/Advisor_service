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

// UpdateUser 更新用户的信息
func UpdateUser(user *model.User) int {

	where := map[string]interface{}{
		"name": user.Name,
	}
	// 把新的用户角色直接转化为map,去掉其中的value为空的key 和 username,password.
	// phone,password,coin不可直接更新
	updates := structs.Map(user)
	delete(updates, "phone")
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
func ChangeUserPWD(username, newPwd string) int {
	// 密码加密
	newPwd = GetPwd(newPwd)
	// 构造sql
	where := map[string]interface{}{
		"name": username,
	}
	updates := map[string]interface{}{
		"password": newPwd,
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

// GetUser 获取用户的全部信息
func GetUser(username string) (int, model.User) {
	where := map[string]interface{}{
		"name": username,
	}
	selects := []string{"name", "phone", "birth", "gender", "bio", "about", "coin"}
	cond, values, err := qb.BuildSelect(USERTABLE, where, selects)
	if err != nil {
		log.Println("gendry SQL生成错误", err)
		return errmsg.ERROR_SQL_BUILD, model.User{}
	}
	row := utils.DbConn.QueryRow(cond, values...)
	res := model.User{}
	// 能不能直接自动赋值到结构体对应的字段?
	err = row.Scan(&res.Name, &res.Phone, &res.Birth, &res.Gender, &res.Bio, &res.About, &res.Coin)
	if err != nil {
		log.Println("数据库查询用户信息出错", err)
		return errmsg.ERROR_MYSQL, model.User{}
	}
	return errmsg.SUCCESS, res
}
