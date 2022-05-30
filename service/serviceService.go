package service

import (
	"database/sql"
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
func NewService(advisorId int64, tx *sql.Tx) int {
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
		logger.GendryBuildError("NewService", err)
		return errmsg.ErrorSqlBuild
	}
	// 执行sql语句
	_, err = tx.Exec(cond, values...)
	if err != nil {
		logger.SqlError("NewService", "Insert", err)
		return errmsg.ErrorSqlBuild
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
		logger.GendryBuildError("ModifyServicePrice", err)
		return errmsg.ErrorSqlBuild
	}
	_, err = utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.SqlError("ModifyServicePrice", "Modify", err)
		return errmsg.ErrorMysql
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
		logger.GendryBuildError("ModifyServiceStatus", err)
		return errmsg.ErrorSqlBuild
	}
	_, err = utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.SqlError("ModifyServiceStatus", "update", err)
		return errmsg.ErrorMysql
	}
	return errmsg.SUCCESS
}

func GetAdvisorService(id int64) (int, []map[string]interface{}) {
	where := map[string]interface{}{
		"advisor_id": id,
		"status":     1,
	}
	selects := []string{"service_name", "service_id", "price"}
	cond, values, err := qb.BuildSelect(SERVICETABLE, where, selects)
	if err != nil {
		logger.GendryBuildError("GetAdvisorService", err)
		return errmsg.ErrorSqlBuild, nil
	}
	rows, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		logger.SqlError("GetAdvisorService", "select", err)
		return errmsg.ErrorMysql, nil
	}
	res, err := scanner.ScanMapDecodeClose(rows)
	if err != nil {
		logger.GendryBuildError("GetAdvisorService", err)
		return errmsg.ErrorSqlBuild, nil
	}
	return errmsg.SUCCESS, res
}
