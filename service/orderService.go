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

var ORDERTABLE = "user_order"

func NewOrderAndCostTrans(model *model.Order) (code int, id int64) {
	begin, err := utils.DbConn.Begin()
	//  defer在函数返回后统一判断事务是否需要回滚
	defer func() {
		if code != errmsg.SUCCESS {
			err = begin.Rollback()
			logger.Log.Error("事务回滚失败!", zap.Error(err))
		}
	}()
	if err != nil {
		return errmsg.ErrorSqlTransError, -1
	}
	//检查金币够不够
	code, userCoin := GetTableItem(USERTABLE, model.UserId, "coin", begin)
	if code != errmsg.SUCCESS {
		return code, -1
	} else if userCoin.(int64) < model.Coin {
		return errmsg.ErrorOrderMoneyInsufficient, -1
	}
	// 扣掉金币
	code = CostUserCoin(model, begin)
	if code != errmsg.SUCCESS {
		return code, -1
	}
	// 新建订单
	code, id = NewOrder(model, begin)
	if code == errmsg.SUCCESS {
		err := begin.Commit()
		if err != nil {
			return errmsg.ErrorSqlTransCommitError, -1
		}
	}
	return code, id
}

// NewOrder 新建订单
func NewOrder(model *model.Order, tx *sql.Tx) (int, int64) {
	// 转化数据并生成sql语句
	var table = ORDERTABLE
	var userMap map[string]interface{}
	var data []map[string]interface{}
	var err error
	// 去出token字段
	userMap = structs.Map(model)

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

// CostUserCoin 扣掉用户的金币
func CostUserCoin(model *model.Order, tx *sql.Tx) int {

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

// ReplyOrderServiceTrans 事务提交 订单回复服务
func ReplyOrderServiceTrans(data *model.OrderReply) (code int) {
	begin, err := utils.DbConn.Begin()
	defer func() {
		if code != errmsg.SUCCESS {
			err = begin.Rollback()
			logger.Log.Error("事务回滚失败!", zap.Error(err))
		}
	}()
	if err != nil {
		code = errmsg.ErrorSqlTransError
		return
	}

	// 获取订单的状态 检测订单是什么状态 只有pending,rush可以回复
	code, OrderStatusInSQL := GetTableItem(ORDERTABLE, data.Id, "status", begin)
	if code != errmsg.SUCCESS {
		code = errmsg.ErrorSqlSingleSelectError
		return
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
		code = errmsg.ErrorOrderHasCompleted
		return
	}

	// 回复订单
	//code = replyOrder(data, begin)
	code = UpdateTableItem(ORDERTABLE, data.Id, map[string]interface{}{
		"reply": data.Reply,
	}, begin)
	if code != errmsg.SUCCESS {
		return code
	}
	// 增加金币
	code = AddCoin2Advisor(data, begin)
	if code != errmsg.SUCCESS {
		return code
	}
	// 加急订单还有额外的金币
	rushId := model.Rush
	if code != errmsg.SUCCESS {
		return code
	}
	if rushId == int(data.Status) {
		code = AddRushCoin2Advisor(data, begin)
		if code != errmsg.SUCCESS {
			return code
		}
	}
	// 标记订单为完成状态
	if code == errmsg.SUCCESS {
		//code = ModifyOrderStatus(data, id, begin)
		code = UpdateTableItem(ORDERTABLE, data.Id, map[string]interface{}{
			"status": model.Completed,
		}, begin)
	} else {
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

func RushOrderTrans(data *model.OrderRush) (code int) {
	begin, err := utils.DbConn.Begin()
	defer func() {
		if code != errmsg.SUCCESS {
			err = begin.Rollback()
			logger.Log.Error("事务回滚失败!", zap.Error(err))
		}
	}()

	if err != nil {
		code = errmsg.ErrorSqlTransError
		return
	}
	// 检测用户的金币是否足够
	code, userMoney := GetTableItem(USERTABLE, data.UserId, "coin", begin)
	if code != errmsg.SUCCESS {
		return code
	}
	//加急需要的钱
	code, orderRushMoney := GetTableItem(ORDERTABLE, data.Id, "rush_coin", begin)
	if code != errmsg.SUCCESS {
		return code
	}
	if userMoney.(int64) < orderRushMoney.(int64) {
		return errmsg.ErrorOrderMoneyInsufficient
	}

	// 状态对吗 获取订单的状态 检测订单是什么状态 只有pending可以加急
	code, OrderStatusInSQL := GetTableItem(ORDERTABLE, data.Id, "status", begin)
	if code != errmsg.SUCCESS {
		code = errmsg.ErrorSqlSingleSelectError
		return
	}
	PendingId := model.Pending
	if PendingId != int(OrderStatusInSQL.(int64)) {
		code = errmsg.ErrorOrderCantRush
		return
	}
	// 修改订单状态为加急
	//code = RushOrder(data.Id)
	code = UpdateTableItem(ORDERTABLE, data.Id, map[string]interface{}{
		"status": model.Rush,
	}, begin)
	if code != errmsg.SUCCESS {
		return code
	}
	// 修改用户金币
	code = UpdateTableItem(USERTABLE, data.UserId, map[string]interface{}{
		"coin": userMoney.(int64) - orderRushMoney.(int64),
	}, begin)
	if code != errmsg.SUCCESS {
		return code
	}
	// 把rush_time提交到表里面去
	code = UpdateTableItem(ORDERTABLE, data.Id, map[string]interface{}{
		"rush_time": data.RushTime,
	}, begin)
	if code != errmsg.SUCCESS {
		return code
	}

	// 事务结束 commit.
	err = begin.Commit()
	if err != nil {
		code = errmsg.ErrorSqlTransCommitError
	}
	return errmsg.SUCCESS
}

// ChangeOrderStatus 修改订单状态 加急->普通,普通->过期
func ChangeOrderStatus(orderId int64, userId int64, originStatus int, targetStatus int) (code int) {
	defer func() {
		m := fmt.Sprintf("用户 %d 订单 %d 状态变化 %s -> %s", userId, orderId,
			model.GetOrderStatusNameById(originStatus),
			model.GetOrderStatusNameById(targetStatus))
		if code == errmsg.SUCCESS {
			logger.Log.Info(m)
		} else {
			logger.Log.Error(m, zap.String("errorMsg", errmsg.GetErrMsg(code)))
		}
	}()
	begin, err := utils.DbConn.Begin()
	//  defer在函数返回后统一判断事务是否需要回滚
	defer func() {
		if code != errmsg.SUCCESS {
			err = begin.Rollback()
			logger.Log.Error("事务回滚失败!", zap.Error(err))
		}
	}()

	if err != nil {
		code = errmsg.ErrorSqlTransError
		return code
	}
	// 获取原始状态
	code, statusInSQL := GetTableItem(ORDERTABLE, orderId, "status", begin)
	if code != errmsg.SUCCESS {
		return code
	}
	// 订单已经完成了 这个状态转移也就结束了，也不需要为用户退回金币
	if int(statusInSQL.(int64)) == model.Completed {
		code = errmsg.ErrorOrderHasCompleted
		return code
	}
	// 加急订单恢复pending状态下，如果被提前标记为expired，可以不用往下执行了.
	if int(statusInSQL.(int64)) == model.Expired {
		code = errmsg.ErrorOrderHasCompleted
		return code
	}
	// 退回金币的逻辑
	code, orderMoney := GetTableItem(ORDERTABLE, orderId, "coin", begin)
	if code != errmsg.SUCCESS {
		return code
	}
	code, orderRushMoney := GetTableItem(ORDERTABLE, orderId, "rush_coin", begin)
	if code != errmsg.SUCCESS {
		return code
	}
	code, userMoney := GetTableItem(USERTABLE, userId, "coin", begin)
	if code != errmsg.SUCCESS {
		return code
	}
	// 用户的金币增加逻辑
	endCoin := userMoney.(float32)
	if originStatus == model.Rush && targetStatus == model.Pending {
		// 加急到普通
		endCoin += orderRushMoney.(float32)
	} else if originStatus == model.Pending && targetStatus == model.Expired {
		// 普通到过期
		endCoin += orderMoney.(float32)
		// 如果订单是在加急的状态下，也要把钱退回去。
		if int(statusInSQL.(int64)) == model.Rush {
			endCoin += orderRushMoney.(float32)
		}
	} else {
		code = errmsg.ErrorJobStatusConvert
		return code
	}

	// 提交用户的金币修改
	code = UpdateTableItem(USERTABLE, userId, map[string]interface{}{
		"coin": endCoin,
	}, begin)
	if code != errmsg.SUCCESS {
		return code
	}
	// 提交状态修改
	code = UpdateTableItem(ORDERTABLE, orderId, map[string]interface{}{
		"status": targetStatus,
	}, begin)
	if code != errmsg.SUCCESS {
		return code
	}
	// 提交事务
	err = begin.Commit()
	if err != nil {
		code = errmsg.ErrorSqlTransCommitError
		return code
	}
	return errmsg.SUCCESS
}
