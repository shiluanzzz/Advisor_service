package service

import (
	"database/sql"
	qb "github.com/didi/gendry/builder"
	"github.com/fatih/structs"
	"go.uber.org/zap"
	"service/model"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
)

var SERVICEKINDTABLE = "service_kind"
var SERVICETABLE = "service"

func NewService(model *model.Service) int {
	var data []map[string]interface{}
	newData := structs.Map(model)
	// id表为另外一个表
	delete(newData, "service_name")
	data = append(data, newData)
	cond, vals, err := qb.BuildInsert(SERVICETABLE, data)
	if err != nil {
		logger.GendryError(err)
		return errmsg.ERROR_SQL_BUILD
	}
	// 执行sql语句
	_, err = utils.DbConn.Exec(cond, vals...)
	if err != nil {
		logger.SqlInsertError(err)
		return errmsg.ERROR_SQL_BUILD
	}
	return errmsg.SUCCESS
}

// 获取service的ID 如果不存在就新建一个ID并返回
func GetServiceId(serviceName string) int {
	row := utils.DbConn.QueryRow("select id from service_kind where name=?", serviceName)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return newServiceKind(serviceName)
		}
		logger.Log.Error("获取服务ID错误", zap.Error(err))
		return -1
	}
	return id
}
func newServiceKind(serviceName string) int {
	res, err := utils.DbConn.Exec("insert into service_kind(name) values (?)", serviceName)
	if err != nil {
		logger.SqlInsertError(err)
		return -1
	}
	newId, _ := res.LastInsertId()
	return int(newId)
}

// CheckService 检查顾问是否存在两种相同的服务
func CheckService(model *model.Service) int {
	return errmsg.SUCCESS
}
func ModifyServicePrice(phone string, id int, price float32) int {
	where := map[string]interface{}{
		"service_id":    id,
		"advisor_phone": phone,
	}
	updates := map[string]interface{}{
		"price": price,
	}
	cond, vals, err := qb.BuildUpdate(SERVICETABLE, where, updates)
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
func ModifyServiceStatus(phone string, id int, newStatus int) int {
	where := map[string]interface{}{
		"service_id":    id,
		"advisor_phone": phone,
	}
	updates := map[string]interface{}{
		"status": newStatus,
	}
	cond, vals, err := qb.BuildUpdate(SERVICETABLE, where, updates)
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
