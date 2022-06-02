package service

import (
	"database/sql"
	qb "github.com/didi/gendry/builder"
	"github.com/fatih/structs"
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
		logger.GendryBuildError("NewUser", err, "cond", cond, "values", values)
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
		logger.SqlError("NewUser", "insert", err, "cond", cond, "values", values)
		return errmsg.ErrorMysql, -1
	}
	// 获取用户的主键ID
	user.Id, err = row.LastInsertId()
	if err != nil {
		logger.SqlError("NewUser", "LastInsertId", err, "cond", cond, "values", values)
	}
	return errmsg.SUCCESS, user.Id
}

// GetId 根据手机号获取用户的ID 用在登录的接口上
func GetId(table, phone string) (id int64, errCode int) {
	where := map[string]interface{}{
		"phone": phone,
	}
	selects := []string{"id"}
	cond, values, err := qb.BuildSelect(table, where, selects)
	if err != nil {
		logger.GendryBuildError("GetId", err, "cond", cond, "values", values)
		return -1, errmsg.ErrorSqlBuild
	}
	row := utils.DbConn.QueryRow(cond, values...)
	var res int64
	err = row.Scan(&res)
	if err != nil && err == sql.ErrNoRows {
		if err == sql.ErrNoRows {
			return -1, errmsg.ErrorUserNotExist
		} else {
			logger.SqlError("GetId", "scan", err, "cond", cond, "values", values)
			return -1, errmsg.ErrorMysql
		}
	} else {
		return res, errmsg.SUCCESS
	}
}

// GetUserInfo 对查询用户信息的方法再次封装，补充消息
func GetUserInfo(id int64) (int, map[string]interface{}) {
	code, data := GetManyTableItemsById(USERTABLE, id, []string{"*"})
	delete(data, "password")
	if gender, ok := data["gender"].(int64); ok {
		data["genderShow"] = model.GetGenderNameById(int(gender))
	} else {
		logger.Log.Error("从数据库中查出的gender转int失败")
	}
	return code, data
}
