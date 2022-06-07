package service

import (
	"database/sql"
	qb "github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
	"service-backend/model"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
	"service-backend/utils/tools"
)

var SERVICETABLE = "service"

func GetService(serviceId int64) (code int, res model.Service) {
	where := map[string]interface{}{
		"id": serviceId,
	}
	code, rows := GetTableRows(SERVICETABLE, where, "*")
	err := scanner.Scan(rows, &res)
	if err != nil {
		return errmsg.ErrorSqlScanner, res
	}
	return errmsg.SUCCESS, res
}

// NewService 新增一个顾客的服务 根据顾问的ID直接新建一套服务
func NewService(advisorId int64, tx *sql.Tx) int {
	var data []map[string]interface{}
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
	cond, values, err := qb.BuildInsert(SERVICETABLE, data)
	if err != nil {
		logger.GendryBuildError(err, "cond", cond, "values", values)
		return errmsg.ErrorSqlBuild
	}
	// 执行sql语句
	_, err = tx.Exec(cond, values...)
	if err != nil {
		logger.SqlError(err, "cond", cond, "values", values)
		return errmsg.ErrorSqlBuild
	}
	return errmsg.SUCCESS
}

func GetAdvisorService(id int64) (code int, res []map[string]interface{}) {

	code, res = GetManyTableItemsByWhere(SERVICETABLE,
		map[string]interface{}{"advisor_id": id, "status": 1},
		[]string{"*"},
	)
	if code != errmsg.SUCCESS {
		return
	}
	for _, v := range res {
		v["price"] = tools.ConvertCoinI2F(v["price"].(int64))
	}
	return
}
