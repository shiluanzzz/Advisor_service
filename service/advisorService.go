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

func GetAdvisorInfo(Id int64) (int, []map[string]interface{}) {
	where := map[string]interface{}{
		"id": Id,
	}
	//selects := []string{
	//	"id", "name", "phone", "coin", "total_order_num", "status",
	//	"rank", "rank_num", "work_experience", "bio", "about",
	//}
	selects := []string{"*"}
	cond, values, err := qb.BuildSelect(ADVISORTABLE, where, selects)
	if err != nil {
		logger.Log.Error("获取顾问信息错误，编译SQL错误", zap.Error(err))
		return errmsg.ErrorSqlBuild, nil
	}
	row, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		logger.Log.Error("数据库查询出错", zap.Error(err))
		return errmsg.ErrorMysql, nil
	}
	res, err := scanner.ScanMapDecodeClose(row)
	if err != nil {
		logger.GendryScannerError("GetAdvisorInfo", err)
		return errmsg.ErrorSqlScanner, nil
	}
	if res == nil {
		return errmsg.ErrorAdvisorNotExist, nil
	}
	// 不传回密码
	for _, each := range res {
		delete(each, "password")
	}
	return errmsg.SUCCESS, res
}

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

func ModifyAdvisorStatus(id int64, newStatus int) int {
	where := map[string]interface{}{
		"id": id,
	}
	updates := map[string]interface{}{
		"status": newStatus,
	}
	cond, values, err := qb.BuildUpdate(ADVISORTABLE, where, updates)
	if err != nil {
		logger.GendryBuildError("ModifyAdvisorStatus", err)
		return errmsg.ErrorSqlBuild
	}
	_, err = utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.SqlError("ModifyAdviosrStatus", "update", err)
		return errmsg.ErrorMysql
	}
	return errmsg.SUCCESS
}

func NewAdvisorAndOrder(data *model.Login) (int, int64) {
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
