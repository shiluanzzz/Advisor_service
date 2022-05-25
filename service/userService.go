package service

import (
	qb "github.com/didi/gendry/builder"
	"github.com/fatih/structs"
	"go.uber.org/zap"
	"service/model"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
)

var USERTABLE = "user"

// NewUser 新增用户
func NewUser(user *model.User) int {
	// 转化数据并生成sql语句
	var data []map[string]interface{}
	data = append(data, structs.Map(user))
	cond, vals, err := qb.BuildInsert(USERTABLE, data)
	if err != nil {
		logger.Log.Error("SQL编译错误", zap.Error(err))
		return errmsg.ERROR_SQL_BUILD
	}

	// 执行sql语句
	_, err = utils.DbConn.Exec(cond, vals...)
	if err != nil {
		logger.Log.Error("数据库插入错误", zap.Error(err))
		return errmsg.ERROR_SQL_BUILD
	}
	return errmsg.SUCCESS
}

// UpdateUser 更新用户的信息
func UpdateUser(user *model.User) int {

	where := map[string]interface{}{
		"phone": user.Phone,
	}
	// 把新的用户角色直接转化为map,去掉不能直接更新的字段
	// phone,password,coin不可直接更新
	updates := structs.Map(user)
	delete(updates, "phone")
	delete(updates, "password")
	delete(updates, "coin")

	// 构造sql 执行更新
	cond, vals, err := qb.BuildUpdate(USERTABLE, where, updates)
	if err != nil {
		logger.GendryError(err)
		return errmsg.ERROR_SQL_BUILD
	}
	_, err = utils.DbConn.Exec(cond, vals...)
	if err != nil {
		logger.SqlUpdateError(err)
		return errmsg.ERROR_MYSQL
	}
	return errmsg.SUCCESS
}

// GetUser 获取用户的全部信息
func GetUser(phone string) (int, model.User) {
	where := map[string]interface{}{
		"phone": phone,
	}
	selects := []string{"name", "phone", "birth", "gender", "bio", "about", "coin"}
	cond, values, err := qb.BuildSelect(USERTABLE, where, selects)
	if err != nil {
		logger.GendryError(err)
		return errmsg.ERROR_SQL_BUILD, model.User{}
	}
	row := utils.DbConn.QueryRow(cond, values...)
	res := model.User{}
	// 能不能直接自动赋值到结构体对应的字段?
	err = row.Scan(&res.Name, &res.Phone, &res.Birth, &res.Gender, &res.Bio, &res.About, &res.Coin)
	if err != nil {
		logger.SqlSelectError(err)
		return errmsg.ERROR_MYSQL, model.User{}
	}
	return errmsg.SUCCESS, res
}
