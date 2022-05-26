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

// NewService 新增一个顾客的服务
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

// GetServiceId 获取service的ID 如果不存在就新建一个ID并返回
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

// 新增一个服务类型
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
func CheckService(ServiceId int, AdvisorPhone string) int {
	where := map[string]interface{}{
		"service_id":    ServiceId,
		"advisor_phone": AdvisorPhone,
	}
	selects := []string{"price"}
	cond, vals, err := qb.BuildSelect(SERVICETABLE, where, selects)
	if err != nil {
		logger.GendryError(err)
		return errmsg.ERROR_SQL_BUILD
	}
	var temp float32
	row := utils.DbConn.QueryRow(cond, vals...)
	err = row.Scan(&temp)
	if err == sql.ErrNoRows {
		return errmsg.ERROR_SERVICE_NOT_EXIST
	} else if err != nil {
		logger.SqlSelectError(err)
		return errmsg.ERROR_MYSQL
	}
	return errmsg.ERROR_SERVICE_EXIST
}

// ModifyServicePrice 修改服务的价格
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

// ModifyServiceStatus 修改服务状态
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
