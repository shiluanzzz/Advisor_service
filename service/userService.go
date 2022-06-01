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

// NewUser 新增用户或者顾问
func NewUser(table string, user *model.Login, tx *sql.Tx) (int, int64) {
	// 转化数据并生成sql语句
	var data []map[string]interface{}
	// 去出token字段
	userMap := structs.Map(user)
	delete(userMap, "token")
	data = append(data, userMap)
	cond, values, err := qb.BuildInsert(table, data)
	if err != nil {
		logger.Log.Error("新增用户错误，编译SQL错误", zap.Error(err))
		return errmsg.ErrorSqlBuild, -1
	}
	var row sql.Result
	// 执行sql语句
	if tx != nil {
		row, err = tx.Exec(cond, values...)
	} else {
		row, err = utils.DbConn.Exec(cond, values...)
	}
	if err != nil {
		logger.Log.Error("数据库插入错误", zap.Error(err))
		return errmsg.ErrorMysql, -1
	}
	// 获取用户的主键ID
	user.Id, err = row.LastInsertId()
	if err != nil {
		logger.Log.Error("获取数据库主键错误", zap.Error(err))
	}
	return errmsg.SUCCESS, user.Id
}
func GetId(table, phone string) (id int64, errCode int) {
	where := map[string]interface{}{
		"phone": phone,
	}
	selects := []string{"id"}
	cond, values, err := qb.BuildSelect(table, where, selects)
	if err != nil {
		logger.Log.Error("获取ID,编译SQL错误", zap.Error(err))
		return -1, errmsg.ErrorSqlBuild
	}
	row := utils.DbConn.QueryRow(cond, values...)
	var res int64
	err = row.Scan(&res)
	if err == sql.ErrNoRows {
		return -1, errmsg.ErrorUserNotExist
	} else {
		return res, errmsg.SUCCESS
	}
}

// GetUserInfo 获取用户的全部信息
func GetUserInfo(userId int64) (int, map[string]interface{}) {
	where := map[string]interface{}{
		"id": userId,
	}
	selects := []string{"*"}
	cond, values, err := qb.BuildSelect(USERTABLE, where, selects)
	if err != nil {
		logger.Log.Error("获取用户信息错误，编译SQL错误", zap.Error(err))
		return errmsg.ErrorSqlBuild, nil
	}

	row, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		logger.Log.Error("数据库查询出错", zap.Error(err))
		return errmsg.ErrorMysql, nil
	}
	res, err := scanner.ScanMapDecodeClose(row)
	if err != nil {
		logger.Log.Error("gendry scanner赋值出错", zap.Error(err))
	}
	if len(res) > 1 {
		return errmsg.ErrorRowNotExpect, nil
	} else if len(res) == 0 {
		return errmsg.ErrorNoResult, nil
	} else {
		t := res[0]
		delete(t, "password")
		return errmsg.SUCCESS, t
	}
}
