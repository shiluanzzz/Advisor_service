package service

import (
	"database/sql"
	"github.com/didi/gendry/scanner"
	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
	"service-backend/model"
	"service-backend/utils/errmsg"
	"service-backend/utils/tools"
)

var SERVICETABLE = "service"

func GetService(serviceId int64) (code int, res model.Service) {
	where := map[string]interface{}{
		"id": serviceId,
	}
	code, rows := GetTableRows(SERVICETABLE, where, "*")
	if err := scanner.Scan(rows, &res); err != nil {
		return errmsg.ErrorSqlScanner, res
	}
	return errmsg.SUCCESS, res
}

// NewService 新增一个顾客的服务 根据顾问的ID直接新建一套服务
func NewService(advisorId int64, tx *sql.Tx) (code int) {
	var data []map[string]interface{}
	// 三种服务！
	for k, v := range model.ServiceKind {
		data = append(data,
			structs.Map(model.Service{
				AdvisorId:     advisorId,
				ServiceName:   v,
				ServiceNameId: k,
				Price:         tools.ConvertCoinF2I(1.0),
				Status:        0,
			}),
		)
	}
	code, _ = InsertTableItem(SERVICETABLE, data, tx)
	return code
}

func GetAdvisorServiceData(advisorId int64) (code int, res []*model.Service) {
	where := map[string]interface{}{
		"advisor_id": advisorId,
	}
	var rows *sql.Rows
	if code, rows = GetTableRows(SERVICETABLE, where, "*"); code != errmsg.SUCCESS {
		return
	}
	if err := scanner.Scan(rows, &res); err != nil {
		return errmsg.ErrorSqlScanner, res
	}
	return errmsg.SUCCESS, res
}
