package service

import (
	"service-backend/model"
	"service-backend/utils"
	"service-backend/utils/cache"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
	"service-backend/utils/tools"
)

var ADVISORTABLE = "advisor"

func GetAdvisor(advisorId int64) (code int, res *model.Advisor) {
	code = GetTableRows2StructByWhere(
		ADVISORTABLE,
		map[string]interface{}{"id": advisorId},
		[]string{"*"},
		&res,
	)
	return errmsg.SUCCESS, res
}
func GetAdvisorList(page int) (code int, res []*model.Advisor) {
	uPage := uint(page)
	where := map[string]interface{}{
		"status": 1,
		"_limit": []uint{(uPage - 1) * 10, uPage * 10},
	}
	selects := []string{
		"id", "phone", "name", "bio", "total_order_num", "total_comment_num", "rank", "on_time",
	}
	code = GetTableRows2StructByWhere(ADVISORTABLE, where, selects, &res)
	return
}

func NewAdvisorAndService(data *model.Login) (code int, id int64) {
	id = -1
	tran, err := utils.DbConn.Begin()
	defer CommonTranDefer(&code, tran)
	if err != nil {
		return errmsg.ErrorSqlTransError, -1
	}
	// 新建用户
	if code, id = NewRole(ADVISORTABLE, data, tran); code != errmsg.SUCCESS {
		return errmsg.ErrorSqlTransError, -1
	}
	// 顾问的服务项创建失败
	if code = NewService(id, tran); code != errmsg.SUCCESS {
		return errmsg.ErrorSqlTransError, -1
	}
	// commit
	if err = tran.Commit(); err != nil {
		return errmsg.ErrorSqlTransCommitError, -1
	}
	return errmsg.SUCCESS, id
}

// GetAdvisorCommentData 获取顾问的订单评论数据
func GetAdvisorCommentData(advisorId int64) (code int, res []*model.OrderComment) {
	// 查询缓存
	var cacheKey = cache.GetCommentKey(advisorId)
	if code = cache.GetCacheData(cacheKey, &res); code == errmsg.SUCCESS {
		return
	}
	defer cache.SetCacheData(cacheKey, &res)

	where := map[string]interface{}{
		"advisor_id":     advisorId,
		"status":         model.Completed,
		"comment_status": model.Commented,
	}
	if code = GetTableRows2StructByWhere(ORDERTABLE, where, []string{"*"}, &res); code != errmsg.SUCCESS {
		return
	}

	return
}

// UpdateAdvisorIndicators week3 更新用户的一些指标信息 TODO 优化
func UpdateAdvisorIndicators(advisorId int64) (code int) {

	var indicators model.AdvisorIndicators
	defer logger.CommonServiceLog(&code, advisorId)
	// 评分
	if code, indicators.Rank = GetAdvisorOrderScore(advisorId); code != errmsg.SUCCESS {
		return
	}
	// 总评论数
	var totalCommentNum interface{}
	where := map[string]interface{}{
		"status":         model.Completed,
		"comment_status": model.Commented,
		"advisor_id":     advisorId,
	}
	if code, totalCommentNum = GetTableItemByWhere(ORDERTABLE, where, "count(id)"); code != errmsg.SUCCESS {
		return
	}
	indicators.TotalCommentNum = int(totalCommentNum.(int64))

	// 总订单数(readings)
	var totalReadings interface{}
	where = map[string]interface{}{
		"advisor_id": advisorId,
		"_or": []map[string]interface{}{
			{"service_name_id": model.VideoReading},
			{"service_name_id": model.AudioReading},
			{"service_name_id": model.TextReading},
		},
	}
	if code, totalReadings = GetTableItemByWhere(ORDERTABLE, where, "count(id)"); code != errmsg.SUCCESS {
		return
	}
	indicators.TotalOrderNum = int(totalReadings.(int64))

	// on-time 订单完成数/总订单数

	// 完成的订单数
	var totalOrderCompleted, totalOrderNum interface{}
	where = map[string]interface{}{
		"advisor_id": advisorId,
		"status":     model.Completed,
	}
	if code, totalOrderCompleted = GetTableItemByWhere(ORDERTABLE, where, "count(id)"); code != errmsg.SUCCESS {
		return
	}
	// 总订单数
	if code, totalOrderNum = GetTableItemByWhere(ORDERTABLE, map[string]interface{}{
		"advisor_id": advisorId,
	}, "count(id)"); code != errmsg.SUCCESS {
		return
	}

	// 计算
	if totalOrderNum.(int64) != 0 {
		indicators.OnTime = float32(totalOrderCompleted.(int64)) / float32(totalOrderNum.(int64))
	} else {
		indicators.OnTime = 0.0
	}
	// 更新入库
	code = UpdateTableItemById(ADVISORTABLE, advisorId, tools.Structs2SQLTable(indicators))
	return
}
