package service

import (
	"database/sql"
	"fmt"
	qb "github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/fatih/structs"
	"go.uber.org/zap"
	"service/model"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
)

const ORDERTABLE = "user_order"

func NewOrderAndCostTrans(model *model.Order) (int, int64) {
	begin, err := utils.DbConn.Begin()
	if err != nil {
		_ = begin.Rollback()
		logger.Log.Error("事务创建失败", zap.Error(err))
		return errmsg.ErrorSqlTransError, -1
	}
	//检查金币够不够
	code, userCoin := GetTableItem(USERTABLE, model.UserId, "coin", begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code, -1
	} else if userCoin.(float32) < model.Coin {
		_ = begin.Rollback()
		return errmsg.ErrorOrderMoneyInsufficient, -1
	}
	// 扣掉金币
	code = CostUserCoin(model, begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code, -1
	}
	// 新建订单
	code, id := NewOrder(model, begin)
	if code == errmsg.SUCCESS {
		err := begin.Commit()
		if err != nil {
			logger.Log.Error("事务提交失败", zap.Error(err))
			return errmsg.ErrorSqlTransCommitError, -1
		}
	} else {
		_ = begin.Rollback()
	}
	return code, id
}

func NewOrder(model *model.Order, tx *sql.Tx) (int, int64) {
	// 转化数据并生成sql语句
	var table = ORDERTABLE
	var data []map[string]interface{}
	// 去出token字段
	userMap := structs.Map(model)
	delete(userMap, "token")
	data = append(data, userMap)
	cond, values, err := qb.BuildInsert(table, data)
	if err != nil {
		logger.Log.Error("新增订单错误，编译SQL错误", zap.Error(err))
		return errmsg.ErrorSqlBuild, -1
	}
	var row sql.Result
	// 执行sql语句
	row, err = tx.Exec(cond, values...)
	if err != nil {
		logger.SqlError("NewOrder", "insert", err, "cond", cond, "values", values)
		return errmsg.ErrorMysql, -1
	}
	// 获取用户的主键ID
	Id, err := row.LastInsertId()
	if err != nil {
		logger.Log.Error("获取数据库主键错误", zap.Error(err))
	}
	return errmsg.SUCCESS, Id
}

func CostUserCoin(model *model.Order, tx *sql.Tx) int {
	//where := map[string]interface{}{
	//	"id": model.UserId,
	//}
	//updates := map[string]interface{}{
	//	"coin": fmt.Sprintf("coin-%f", model.Coin),
	//}
	//cond, values, err := qb.BuildUpdate(USERTABLE, where, updates)
	//if err != nil {
	//	logger.GendryBuildError("orderService.CostUserCoin", err, "cond", cond, "update", updates)
	//	return errmsg.ErrorSqlBuild
	//}
	cond := "update `user` set coin = coin - ? where id = ?"
	row, err := tx.Exec(cond, model.Coin, model.UserId)
	if err != nil {
		logger.SqlError("orderService.CostUserCoin", "update", err, "cond", cond)
		return errmsg.ErrorMysql
	}
	affects, _ := row.RowsAffected()
	if affects != 1 {
		logger.Log.Error("用户金币修改设计到多个行列", zap.Int64("userId", model.UserId))
		return errmsg.ErrorAffectsNotOne
	}
	return errmsg.SUCCESS
}

func GetOrderList(advisorId int64) (int, []map[string]interface{}) {
	where := map[string]interface{}{
		"advisor_id": advisorId,
	}
	selects := []string{"user_id", "id", "question", "create_time", "service_id", "status"}
	cond, values, err := qb.BuildSelect(ORDERTABLE, where, selects)
	if err != nil {
		logger.GendryBuildError("orderService.GetOrderList", err, "cond", cond, "values", values)
		return errmsg.ErrorSqlBuild, nil
	}
	rows, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		logger.SqlError("orderService.GetOrderList", "select", err, "cond", cond, "values", values)
		return errmsg.ErrorMysql, nil
	}
	res, err := scanner.ScanMapDecodeClose(rows)
	if err != nil {
		logger.GendryScannerError("orderService.GetOrderList", err, "cond", cond, "values", values)
		return errmsg.ErrorSqlScanner, nil
	}
	return errmsg.SUCCESS, res
}

func GetOrderInfo(orderId int) (int, map[string]interface{}) {
	where := map[string]interface{}{
		"id": orderId,
	}
	selects := []string{"*"}
	cond, values, err := qb.BuildSelect(ORDERTABLE, where, selects)
	if err != nil {
		logger.GendryBuildError("orderService.GetOrderInfo", err, "cond", cond, "values", values)
		return errmsg.ErrorSqlBuild, nil
	}
	rows, err := utils.DbConn.Query(cond, values...)
	if err != nil {
		logger.SqlError("orderService.GetOrderInfo", "select", err, "cond", cond, "values", values)
		return errmsg.ErrorMysql, nil
	}
	ress, err := scanner.ScanMapDecodeClose(rows)
	if err != nil {
		logger.GendryScannerError("orderService.GetOrderInfo", err, "cond", cond, "values", values)
		return errmsg.ErrorSqlScanner, nil
	}
	if len(ress) > 1 {
		return errmsg.ErrorRowNotExpect, nil
	} else if len(ress) == 0 {
		return errmsg.ErrorNoResult, nil
	}
	return errmsg.SUCCESS, ress[0]
}

// ReplyOrderServiceTrans 事务提交 订单回复服务
func ReplyOrderServiceTrans(data *model.OrderReply) int {
	begin, err := utils.DbConn.Begin()
	if err != nil {
		_ = begin.Rollback()
		logger.Log.Error("事务创建失败", zap.Error(err))
		return errmsg.ErrorSqlTransError
	}

	// 获取订单的状态 检测订单是什么状态 只有pending,rush可以回复
	code, OrderStatusInSQL := GetTableItem(ORDERTABLE, data.Id, "status", begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return errmsg.ErrorSqlSingleSelectError
	}
	canReply := false
	for _, i := range model.GetOrderEnableReplyId() {
		if int(OrderStatusInSQL.(int64)) == i {
			data.Status = OrderStatusInSQL.(int64)
			canReply = true
			break
		}
	}
	if !canReply {
		_ = begin.Rollback()
		return errmsg.ErrorOrderHasCompleted
	}

	// 回复订单
	code = replyOrder(data, begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	// 增加金币
	code = AddCoin2Advisor(data, begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	// 加急订单还有额外的金币
	rushId := model.Rush
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	if rushId == int(data.Status) {
		code = AddRushCoin2Advisor(data, begin)
		if code != errmsg.SUCCESS {
			_ = begin.Rollback()
			return code
		}
	}
	// 标记订单为完成状态
	id := model.Completed
	if code == errmsg.SUCCESS {
		code = ModifyOrderStatus(data, id, begin)
	} else {
		_ = begin.Rollback()
		return code
	}
	// 事务终于结束了ho
	err = begin.Commit()
	if err != nil {
		logger.Log.Error("事务最终提交失败", zap.Error(err))
		return errmsg.ErrorSqlTransCommitError
	}
	return errmsg.SUCCESS
}

// AddCoin2Advisor 在顾问回复订单后为顾问增加金币
func AddCoin2Advisor(data *model.OrderReply, tx *sql.Tx) int {
	cond := fmt.Sprintf("update %s set coin=coin + ? where id= ?", ADVISORTABLE)
	row, err := tx.Exec(cond, data.Coin, data.AdvisorId)
	if err != nil {
		logger.SqlError("AddCoin2Advisor", "update", err, "cond", cond)
		return errmsg.ErrorMysql
	}
	affects, _ := row.RowsAffected()
	if affects != 1 {
		logger.Log.Error("用户金币修改设计到多个行列", zap.Int64("advisor_id", data.AdvisorId))
		return errmsg.ErrorAffectsNotOne
	}
	return errmsg.SUCCESS
}

// AddRushCoin2Advisor 为顾问增加加急订单的金币
func AddRushCoin2Advisor(data *model.OrderReply, tx *sql.Tx) int {
	cond := fmt.Sprintf("update %s set coin=coin + ? where id= ?", ADVISORTABLE)
	row, err := tx.Exec(cond, data.RushCoin, data.AdvisorId)
	if err != nil {
		logger.SqlError("AddCoin2Advisor", "update", err, "cond", cond)
		return errmsg.ErrorMysql
	}
	affects, _ := row.RowsAffected()
	if affects != 1 {
		logger.Log.Error("用户金币修改设计到多个行列", zap.Int64("advisor_id", data.AdvisorId))
		return errmsg.ErrorAffectsNotOne
	}
	return errmsg.SUCCESS
}

// replyOrder 回复订单
func replyOrder(data *model.OrderReply, tx *sql.Tx) int {
	// 回复订单后为顾问增加金币
	where := map[string]interface{}{
		"id": data.Id,
	}
	update := map[string]interface{}{
		"reply": data.Reply,
	}
	cond, values, err := qb.BuildUpdate(ORDERTABLE, where, update)
	if err != nil {
		logger.GendryBuildError("replyOrder", err, "cond", cond, "values", values)
		return errmsg.ErrorSqlBuild
	}
	row, err := tx.Exec(cond, values...)
	if err != nil {
		logger.SqlError("replyOrder", "update", err, "cond", cond)
		return errmsg.ErrorMysql
	}
	affects, _ := row.RowsAffected()
	if affects != 1 {
		logger.Log.Error("修改订单回复", zap.Int64("order_id", data.Id))
		return errmsg.ErrorAffectsNotOne
	}
	return errmsg.SUCCESS
}

// ModifyOrderStatus 修改订单的状态
func ModifyOrderStatus(data *model.OrderReply, status int, tx *sql.Tx) int {
	where := map[string]interface{}{
		"id": data.Id,
	}
	update := map[string]interface{}{
		"status": status,
	}
	cond, values, err := qb.BuildUpdate(ORDERTABLE, where, update)
	if err != nil {
		logger.GendryBuildError("ModifyOrderStatus", err, "cond", cond, "values", values)
		return errmsg.ErrorSqlBuild
	}
	var row sql.Result
	if tx != nil {
		row, err = tx.Exec(cond, values...)
	} else {
		row, err = utils.DbConn.Exec(cond, values...)
	}
	if err != nil {
		logger.SqlError("ModifyOrderStatus", "update", err, "cond", cond)
		return errmsg.ErrorMysql
	}
	affects, _ := row.RowsAffected()
	if affects != 1 {
		logger.Log.Error("修改订单状态涉及到多条数据", zap.Int64("order_id", data.Id))
		return errmsg.ErrorAffectsNotOne
	}
	return errmsg.SUCCESS
}

func RushOrderTrans(data *model.OrderRush) int {
	begin, err := utils.DbConn.Begin()
	if err != nil {
		_ = begin.Rollback()
		logger.Log.Error("事务创建失败", zap.Error(err))
		return errmsg.ErrorSqlTransError
	}
	// 检测用户的金币是否足够
	code, userMoney := GetTableItem(USERTABLE, data.UserId, "coin", begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	//加急需要的钱
	code, orderRushMoney := GetTableItem(ORDERTABLE, data.Id, "rush_coin", begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	if userMoney.(float32) < orderRushMoney.(float32) {
		return errmsg.ErrorOrderMoneyInsufficient
	}
	// 状态对吗 获取订单的状态 检测订单是什么状态 只有pending可以加急
	code, OrderStatusInSQL := GetTableItem(ORDERTABLE, data.Id, "status", begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return errmsg.ErrorSqlSingleSelectError
	}
	PendingId := model.Pending
	if PendingId != int(OrderStatusInSQL.(int64)) {
		_ = begin.Rollback()
		return errmsg.ErrorOrderCantRush
	}
	// 修改订单状态
	code = RushOrder(data.Id)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	// 修改用户金币
	code = UpdateTableItem(USERTABLE, data.UserId, map[string]interface{}{
		"coin": userMoney.(float32) - orderRushMoney.(float32),
	}, begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	err = begin.Commit()
	if err != nil {
		_ = begin.Rollback()
		return errmsg.ErrorSqlTransCommitError
	}
	return errmsg.SUCCESS
}

// RushOrder 订单加急
func RushOrder(id int64) int {
	orderId := model.Rush
	code := ModifyOrderStatus(&model.OrderReply{Id: id}, orderId, nil)
	return code
}

// ChangeOrderStatus 修改订单状态 加急->普通,普通->过期
func ChangeOrderStatus(orderId int64, userId int64, originStatus int, targetStatus int) int {
	begin, err := utils.DbConn.Begin()
	if err != nil {
		_ = begin.Rollback()
		return errmsg.ErrorSqlTransError
	}
	// 获取原始状态
	code, statusInSQL := GetTableItem(ORDERTABLE, orderId, "status", begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	// 订单已经完成了 这个状态转移也就结束了，也不需要为用户退回金币
	if int(statusInSQL.(int64)) == model.Completed {
		_ = begin.Rollback()
		return errmsg.ErrorOrderHasCompleted
	}
	// 订单与预期状态不符合 TODO?
	if int(statusInSQL.(int64)) != originStatus {
		_ = begin.Rollback()
		return errmsg.ErrorJobStatusNotExpect
	}
	// 退回金币的逻辑
	code, orderMoney := GetTableItem(ORDERTABLE, orderId, "coin", begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	code, orderRushMoney := GetTableItem(ORDERTABLE, orderId, "rush_coin", begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	code, userMoney := GetTableItem(USERTABLE, userId, "coin", begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	// 用户的金币增加逻辑
	endCoin := userMoney.(float32)
	if originStatus == model.Rush && targetStatus == model.Pending {
		// 加急到普通
		endCoin += orderRushMoney.(float32)
	} else if originStatus == model.Pending && targetStatus == model.Expired {
		// 加急到过期
		endCoin += orderRushMoney.(float32)
		endCoin += orderMoney.(float32)
	} else if originStatus == model.Pending && targetStatus == model.Expired {
		// 普通到过期
		endCoin += orderMoney.(float32)
	}
	// 提交用户的金币修改
	code = UpdateTableItem(USERTABLE, userId, map[string]interface{}{
		"coin": endCoin,
	}, begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	// 提交状态修改
	code = UpdateTableItem(ORDERTABLE, orderId, map[string]interface{}{
		"status": targetStatus,
	}, begin)
	if code != errmsg.SUCCESS {
		_ = begin.Rollback()
		return code
	}
	// 提交事务
	err = begin.Commit()
	if err != nil {
		_ = begin.Rollback()
		return errmsg.ErrorSqlTransCommitError
	}
	return errmsg.SUCCESS
}
