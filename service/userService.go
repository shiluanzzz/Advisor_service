package service

import (
	"database/sql"
	qb "github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/fatih/structs"
	"go.uber.org/zap"
	"service/model"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
)

var USERTABLE = "user"

// NewUser 新增用户
func NewUser(user *model.UserLogin) int {
	// 转化数据并生成sql语句
	var data []map[string]interface{}
	// 去出token字段
	userMap := structs.Map(user)
	delete(userMap, "token")
	data = append(data, userMap)
	cond, values, err := qb.BuildInsert(USERTABLE, data)
	if err != nil {
		logger.Log.Error("新增用户错误，编译SQL错误", zap.Error(err))
		return errmsg.ERROR_SQL_BUILD
	}

	// 执行sql语句
	row, err := utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.Log.Error("数据库插入错误", zap.Error(err))
		return errmsg.ERROR_MYSQL
	}
	// 获取用户的主键ID
	user.Id, err = row.LastInsertId()
	if err != nil {
		logger.Log.Error("获取数据库主键错误", zap.Error(err))
	}
	return errmsg.SUCCESS
}
func GetUserId(phone string) (id int64, errCode int) {
	where := map[string]interface{}{
		"phone": phone,
	}
	selects := []string{"id"}
	cond, values, err := qb.BuildSelect(USERTABLE, where, selects)
	if err != nil {
		logger.Log.Error("获取用户ID,编译SQL错误", zap.Error(err))
		return -1, errmsg.ERROR_SQL_BUILD
	}
	row := utils.DbConn.QueryRow(cond, values...)
	var res int64
	err = row.Scan(&res)
	if err == sql.ErrNoRows {
		return -1, errmsg.ERROR_USER_NOT_EXIST
	} else {
		return res, errmsg.SUCCESS
	}
}

// UpdateUser 更新用户的信息
func UpdateUser(userInfo map[string]interface{}) int {

	where := map[string]interface{}{
		"id": userInfo["id"],
	}
	// 密码、coin不能直接更新
	InValidField := []string{"coin", "password"}
	for _, filed := range InValidField {
		if userInfo[filed] != nil {
			return errmsg.ERROR_UPDATE_VALID
		}
	}
	// 构造sql 执行更新
	cond, values, err := qb.BuildUpdate(USERTABLE, where, userInfo)
	if err != nil {
		logger.Log.Error("更新用户信息错误，编译SQL错误", zap.Error(err))
		return errmsg.ERROR_SQL_BUILD
	}
	_, err = utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.SqlUpdateError(err)
		return errmsg.ERROR_MYSQL
	}
	return errmsg.SUCCESS
}

// GetUser 获取用户的全部信息
func GetUser(userId int64) (int, interface{}) {
	where := map[string]interface{}{
		"id": userId,
	}
	// 一次获取全部key? TODO
	selects := []string{"name", "phone", "birth", "gender", "bio", "about", "coin"}
	cond, values, err := qb.BuildSelect(USERTABLE, where, selects)
	if err != nil {
		logger.Log.Error("获取用户信息错误，编译SQL错误", zap.Error(err))
		return errmsg.ERROR_SQL_BUILD, nil
	}

	row, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		logger.Log.Error("数据库查询出错", zap.Error(err))
		return errmsg.ERROR_MYSQL, nil
	}
	res, err := scanner.ScanMapDecodeClose(row)
	if err != nil {
		logger.Log.Error("gendry scanner赋值出错", zap.Error(err))
	}
	return errmsg.SUCCESS, res
}
