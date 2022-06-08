package service

import (
	"database/sql"
	qb "github.com/didi/gendry/builder"
	"github.com/fatih/structs"
	"service-backend/model"
	"service-backend/utils"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
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
		logger.GendryBuildError(err, "cond", cond, "values", values)
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
		logger.SqlError(err, "cond", cond, "values", values)
		return errmsg.ErrorMysql, -1
	}
	// 获取用户的主键ID
	user.Id, err = row.LastInsertId()
	if err != nil {
		logger.SqlError(err, "cond", cond, "values", values)
	}
	return errmsg.SUCCESS, user.Id
}

// GetUser 对查询用户信息的方法再次封装，补充消息
func GetUser(id int64) (code int, res model.User) {
	code = GetTableRows2StructByWhere(USERTABLE, map[string]interface{}{"id": id}, []string{"*"}, &res)
	res.UpdateShow("02-01-2006")
	return code, res
}
