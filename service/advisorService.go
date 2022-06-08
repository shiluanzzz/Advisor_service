package service

import (
	"database/sql"
	"fmt"
	"github.com/didi/gendry/scanner"
	"service-backend/model"
	"service-backend/utils"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
	"service-backend/utils/tools"
	"time"
)

var ADVISORTABLE = "advisor"

func GetAdvisor(advisorId int64) (code int, res model.Advisor) {
	var err error
	defer logger.CommonServiceLog(&code, advisorId, "err", err)
	where := map[string]interface{}{
		"id": advisorId,
	}
	code, rows := GetTableRows(ADVISORTABLE, where, "*")
	if err = scanner.Scan(rows, &res); err != nil {
		return errmsg.ErrorSqlScanner, res
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
	return GetManyTableItemsByWhere(ADVISORTABLE, where, selects)
}

func NewAdvisorAndService(data *model.Login) (code int, id int64) {
	id = -1
	begin, err := utils.DbConn.Begin()
	defer CommonTranDefer(&code, begin)
	if err != nil {
		return errmsg.ErrorSqlTransError, -1
	}
	// 新建用户
	if code, id = NewUser(ADVISORTABLE, data, begin); code != errmsg.SUCCESS {
		return errmsg.ErrorSqlTransError, -1
	}
	// 顾问的服务项创建失败
	if code = NewService(id, begin); code != errmsg.SUCCESS {
		return errmsg.ErrorSqlTransError, -1
	}
	// commit
	err = begin.Commit()
	if err != nil {
		return errmsg.ErrorSqlTransCommitError, -1
	}
	return errmsg.SUCCESS, id
}

// GetAdvisorCommentData 获取顾问的订单评论数据
func GetAdvisorCommentData(id int64) (code int, res []*model.OrderComment) {
	defer logger.CommonServiceLog(&code, id)
	where := map[string]interface{}{
		"advisor_id":     id,
		"status":         model.Completed,
		"comment_status": model.Commented,
	}
	var rows *sql.Rows
	if code, rows = GetTableRows(ORDERTABLE, where, "*"); code != errmsg.SUCCESS {
		return
	}
	if err := scanner.Scan(rows, &res); err != nil {
		return errmsg.ErrorSqlScanner, nil
	}
	// 扩充数据
	for _, v := range res {
		var userNameUint8 interface{}
		if code, userNameUint8 = GetTableItem(USERTABLE, v.UserId, "name"); code != errmsg.SUCCESS {
			return
		}
		v.UserName = fmt.Sprintf("%s", userNameUint8)
		v.CreateShowTime = time.Unix(v.OrderCreateTime, 0).Format("Jan 02,2006 15:04:05")
		v.CommentShowTime = time.Unix(v.CommentTime, 0).Format("Jan 02,2006 15:04:05")
	}
	return
}

// UpdateAdvisorIndicators week3 更新用户的一些指标信息
func UpdateAdvisorIndicators(advisorId int64, tx ...*sql.Tx) (code int) {

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
