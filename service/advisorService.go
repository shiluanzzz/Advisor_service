package service

import (
	qb "github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"go.uber.org/zap"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
)

var ADVISORTABLE = "advisor"

func GetAdvisorInfo(Id int64) (int, []map[string]interface{}) {
	where := map[string]interface{}{
		"id": Id,
	}
	selects := []string{
		"name", "phone", "coin", "total_order_num", "status",
		"rank", "rank_num", "work_experience", "bio", "about",
	}
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
		logger.Log.Error("gendry scanner赋值出错", zap.Error(err))
	}
	if res == nil {
		return errmsg.ErrorAdvisorNotExist, nil
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
		"phone", "name", "bio",
	}
	cond, values, err := qb.BuildSelect(ADVISORTABLE, where, selects)
	if err != nil {
		logger.GendryError("GetAdvisorList", err)
		return errmsg.ErrorSqlBuild, nil
	}
	rows, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		logger.SqlError("GetAdvisorList", "select", err)
		return errmsg.ERROR, nil
	}
	res, err := scanner.ScanMapDecodeClose(rows)
	if err != nil {
		logger.GendryError("GetAdvisorList", err)
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
		logger.GendryError("ModifyAdvisorStatus", err)
		return errmsg.ErrorSqlBuild
	}
	_, err = utils.DbConn.Exec(cond, values...)
	if err != nil {
		logger.SqlError("ModifyAdviosrStatus", "update", err)
		return errmsg.ErrorMysql
	}
	return errmsg.SUCCESS
}
