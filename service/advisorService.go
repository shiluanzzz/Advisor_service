package service

import (
	qb "github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"go.uber.org/zap"
	"service/model"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
)

var ADVISORTABLE = "advisor"

func GetAdvisorList(page int) (int, []map[string]interface{}) {
	uPage := uint(page)
	where := map[string]interface{}{
		"status": 1,
		"_limit": []uint{(uPage - 1) * 10, uPage * 10},
	}
	selects := []string{
		"id", "phone", "name", "bio",
	}
	cond, values, err := qb.BuildSelect(ADVISORTABLE, where, selects)
	if err != nil {
		logger.GendryBuildError("GetAdvisorList", err)
		return errmsg.ErrorSqlBuild, nil
	}
	rows, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		logger.SqlError("GetAdvisorList", "select", err)
		return errmsg.ERROR, nil
	}
	res, err := scanner.ScanMapDecodeClose(rows)
	if err != nil {
		logger.GendryScannerError("GetAdvisorList", err)
		return errmsg.ErrorSqlBuild, nil
	}
	return errmsg.SUCCESS, res
}

func NewAdvisorAndService(data *model.Login) (int, int64) {
	var code int
	var id int64
	begin, err := utils.DbConn.Begin()
	if err != nil {
		_ = begin.Rollback()
		logger.Log.Error("事务创建失败", zap.Error(err))
		return errmsg.ErrorSqlTransError, -1
	}
	code, id = NewUser(ADVISORTABLE, data, begin)
	if code != errmsg.SUCCESS {
		// 创建顾问失败
		_ = begin.Rollback()
		logger.Log.Warn("顾问数据创建错误,事务回滚")
		return errmsg.ErrorSqlTransError, -1
	}
	code = NewService(id, begin)
	if code != errmsg.SUCCESS {
		// 顾问的服务项创建失败
		_ = begin.Rollback()
		logger.Log.Warn("顾问服务创建失败，事务回滚")
		return errmsg.ErrorSqlTransError, -1
	}
	err = begin.Commit()
	if err != nil {
		_ = begin.Rollback()
		logger.Log.Error("事务提交失败", zap.Error(err))
		return errmsg.ErrorSqlTransCommitError, -1
	}
	return errmsg.SUCCESS, id
}
