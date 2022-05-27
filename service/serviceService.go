package service

import (
	qb "github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
	"service/model"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
)

var SERVICETABLE = "service"

// NewService 新增一个顾客的服务 根据顾问的ID直接新建一套服务
func NewService(advisorId int64) int {
	var data []map[string]interface{}
	for k, v := range model.ServiceKind {
		data = append(data,
			structs.Map(model.Service{
				AdvisorId:   advisorId,
				ServiceName: v,
				ServiceId:   k,
				Price:       1,
				Status:      0,
			}),
		)
	}
	cond, values, err := qb.BuildInsert(SERVICETABLE, data)
	if err != nil {
		logger.GendryError("NewService", err)
		return errmsg.ERROR_SQL_BUILD
	}
	// 执行sql语句
	_, err = utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.SqlError("NewService", "Insert", err)
		return errmsg.ERROR_SQL_BUILD
	}
	return errmsg.SUCCESS
}

// ModifyServicePrice 修改服务的价格
func ModifyServicePrice(advisorId int64, serviceId int, price float32) int {
	where := map[string]interface{}{
		"service_id": serviceId,
		"advisor_id": advisorId,
	}
	updates := map[string]interface{}{
		"price": price,
	}
	cond, values, err := qb.BuildUpdate(SERVICETABLE, where, updates)
	if err != nil {
		logger.GendryError("ModifyServicePrice", err)
		return errmsg.ERROR_SQL_BUILD
	}
	_, err = utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.SqlError("ModifyServicePrice", "Modify", err)
		return errmsg.ERROR_MYSQL
	}
	return errmsg.SUCCESS
}

// ModifyServiceStatus 修改服务状态
func ModifyServiceStatus(advisorId int64, serviceId int, newStatus int) int {
	where := map[string]interface{}{
		"service_id": serviceId,
		"advisor_id": advisorId,
	}
	updates := map[string]interface{}{
		"status": newStatus,
	}
	cond, values, err := qb.BuildUpdate(SERVICETABLE, where, updates)
	if err != nil {
		logger.GendryError("ModifyServiceStatus", err)
		return errmsg.ERROR_SQL_BUILD
	}
	_, err = utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.SqlError("ModifyServiceStatus", "update", err)
		return errmsg.ERROR_MYSQL
	}
	return errmsg.SUCCESS
}

func GetAdvisorService(id int64) (int, interface{}) {
	where := map[string]interface{}{
		"advisor_id": id,
		"status":     1,
	}
	selects := []string{"service_name", "service_id", "price"}
	cond, values, err := qb.BuildSelect(SERVICETABLE, where, selects)
	if err != nil {
		logger.GendryError("GetAdvisorService", err)
		return errmsg.ERROR_SQL_BUILD, nil
	}
	rows, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		logger.SqlError("GetAdvisorService", "select", err)
		return errmsg.ERROR_MYSQL, nil
	}
	res, err := scanner.ScanMapDecodeClose(rows)
	if err != nil {
		logger.GendryError("GetAdvisorService", err)
		return errmsg.ERROR_SQL_BUILD, nil
	}
	return errmsg.SUCCESS, res
}
