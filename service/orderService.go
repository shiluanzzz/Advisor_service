package service

import (
	"database/sql"
	"fmt"
	"github.com/didi/gendry/scanner"
	"go.uber.org/zap"
	"service-backend/model"
	"service-backend/utils"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
	"service-backend/utils/tools"
	"time"
)

var ORDERTABLE = "user_order"

func NewOrderAndCostTrans(model *model.Order) (code int, id int64) {
	begin, err := utils.DbConn.Begin()
	defer CommonTranDefer(&code, begin)
	if err != nil {
		return errmsg.ErrorSqlTransError, -1
	}

	// 扣掉金币
	if code = CostUserCoin(model, begin); code != errmsg.SUCCESS {
		return code, -1
	}
	// 新建订单
	if code, id = NewOrder(model, begin); code != errmsg.SUCCESS {
		return code, -1
	}

	// 提交事务
	if err := begin.Commit(); err != nil {
		return errmsg.ErrorSqlTransCommitError, -1
	}
	return code, id
}

// NewOrder 新建订单
func NewOrder(model *model.Order, tx *sql.Tx) (code int, id int64) {
	// 转化数据并生成sql语句
	var userMap map[string]interface{}
	var data []map[string]interface{}
	// 转化为map
	userMap = tools.Structs2SQLTable(model)
	delete(userMap, "token")
	data = append(data, userMap)
	code, id = InsertTableItem(ORDERTABLE, data, tx)
	return errmsg.SUCCESS, id
}

// CostUserCoin 扣掉用户的金币
func CostUserCoin(model *model.Order, tx *sql.Tx) int {

	cond := "update `user` set coin = coin - ? where id = ?"
	row, err := tx.Exec(cond, model.Coin, model.UserId)
	if err != nil {
		logger.SqlError(err, "cond", cond)
		return errmsg.ErrorMysql
	}
	affects, _ := row.RowsAffected()
	if affects != 1 {
		logger.Log.Error("用户金币修改设计到多个行列", zap.Int64("userId", model.UserId))
		return errmsg.ErrorAffectsNotOne
	}
	return errmsg.SUCCESS
}

// ReplyOrderServiceTrans 事务提交 订单回复服务
func ReplyOrderServiceTrans(data *model.OrderReply) (code int) {
	begin, err := utils.DbConn.Begin()
	defer CommonTranDefer(&code, begin)
	if err != nil {
		code = errmsg.ErrorSqlTransError
		return
	}

	// 回复订单并标记订单为完成状态
	code = UpdateTableItemById(ORDERTABLE, data.Id, map[string]interface{}{
		"reply":  data.Reply,
		"status": model.Completed,
	}, begin)
	if code != errmsg.SUCCESS {
		return
	}
	// 增加顾问的金币
	reward := data.Coin
	if data.Status == model.Rush {
		reward += data.RushCoin
	}
	if code = AddCoin2Advisor(data, begin); code != errmsg.SUCCESS {
		return
	}

	// 事务终于结束了ho
	err = begin.Commit()
	if err != nil {
		return errmsg.ErrorSqlTransCommitError
	}
	return errmsg.SUCCESS
}

// AddCoin2Advisor 在顾问回复订单后为顾问增加金币
func AddCoin2Advisor(data *model.OrderReply, tx *sql.Tx) int {
	cond := fmt.Sprintf("update %s set coin=coin + ? where id= ?", ADVISORTABLE)
	row, err := tx.Exec(cond, data.Coin, data.AdvisorId)
	if err != nil {
		logger.SqlError(err, "cond", cond)
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
	defer CommonTranDefer(&code, begin)
	if err != nil {
		code = errmsg.ErrorSqlTransError
		return
	}

	// 修改订单状态为加急
	code = UpdateTableItemById(ORDERTABLE, data.Id, map[string]interface{}{
		"status":    model.Rush,
		"rush_time": data.RushTime,
	}, begin)
	if code != errmsg.SUCCESS {
		return code
	}
	// 修改用户金币
	code = UpdateTableItemById(USERTABLE, data.UserId, map[string]interface{}{
		"coin": data.UserMoney - data.RushMoney,
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
func ChangeOrderStatus(orderId int64, userId int64, originStatus model.OrderStatus, targetStatus model.OrderStatus) (code int) {
	// defer log
	defer func() {
		m := fmt.Sprintf("用户 %d 订单 %d 状态变化 %s -> %s", userId, orderId,
			originStatus.StatusName(),
			targetStatus.StatusName())
		if code == errmsg.SUCCESS {
			logger.Log.Info(m)
		} else {
			logger.Log.Error(m, zap.String("errorMsg", errmsg.GetErrMsg(code)))
		}
	}()

	Tran, err := utils.DbConn.Begin()
	defer CommonTranDefer(&code, Tran)
	if err != nil {
		code = errmsg.ErrorSqlTransError
		return code
	}
	// 获取原始状态
	var orderInSql model.Order
	if code, orderInSql = GetOrder(orderId); code != errmsg.SUCCESS {
		return
	}
	if orderInSql.Status == model.Completed {
		// 订单已经完成了 这个状态转移也就结束了，也不需要为用户退回金币
		code = errmsg.ErrorOrderHasCompleted
		return code
	}
	if orderInSql.Status == model.Expired {
		// 加急订单恢复pending状态下，如果被提前标记为expired，可以不用往下执行了.
		code = errmsg.ErrorOrderHasCompleted
		return code
	}

	code, userMoney := GetTableItem(USERTABLE, userId, "coin", Tran)
	if code != errmsg.SUCCESS {
		return code
	}
	// 退回用户的金币增加逻辑
	// 退回用户金币的逻辑
	endCoin := userMoney.(int64)
	if originStatus == model.Rush && targetStatus == model.Pending {
		// 加急到普通
		endCoin += orderInSql.RushCoin
	} else if originStatus == model.Pending && targetStatus == model.Expired {
		// 普通到过期
		endCoin += orderInSql.Coin
		// 如果订单是在加急的状态下，也要把钱退回去。
		if orderInSql.Status == model.Rush {
			endCoin += orderInSql.RushCoin
		}
	} else {
		code = errmsg.ErrorJobStatusConvert
		return code
	}

	// 提交用户的金币修改
	code = UpdateTableItemById(USERTABLE, userId, map[string]interface{}{
		"coin": endCoin,
	}, Tran)
	if code != errmsg.SUCCESS {
		return code
	}
	// 提交状态修改
	code = UpdateTableItemById(ORDERTABLE, orderId, map[string]interface{}{
		"status": targetStatus,
	}, Tran)
	if code != errmsg.SUCCESS {
		return code
	}
	// 提交事务
	err = Tran.Commit()
	if err != nil {
		code = errmsg.ErrorSqlTransCommitError
		return code
	}
	return errmsg.SUCCESS
}

func GetAdvisorOrderList(advisorId int64) (code int, res []*model.Order) {
	where := map[string]interface{}{
		"advisor_id": advisorId,
	}
	selects := []string{"id", "user_id", "service_id", "service_name_id",
		"status", "question", "situation", "advisor_id", "create_time"}
	code, rows := GetTableRows(ORDERTABLE, where, selects...)
	err := scanner.Scan(rows, &res)
	if err != nil {
		return errmsg.ErrorSqlScanner, nil
	}
	// 附加信息:用户名、时间格式、服务类型
	for _, v := range res {
		v.ShowTime = time.Unix(v.CreateTime, 0).Format("Jan 02,2006")
		_, UserName := GetTableItem(USERTABLE, v.UserId, "name")
		v.UserName = fmt.Sprintf("%s", UserName)
		v.ServiceName = model.ServiceKind[v.ServiceNameId]
		v.ServiceStatusName = v.Status.StatusName()
	}
	return code, res
}
func GetOrder(orderId int64) (code int, res model.Order) {
	var err error
	defer logger.CommonServiceLog(&code, orderId, "err", err)
	where := map[string]interface{}{
		"id": orderId,
	}
	var rows *sql.Rows
	if code, rows = GetTableRows(ORDERTABLE, where, "*"); code != errmsg.SUCCESS {
		return code, res
	}
	if err = scanner.Scan(rows, &res); err != nil {
		// TODO 更新到其他位置
		if err == scanner.ErrNilRows || err == scanner.ErrEmptyResult {
			return errmsg.ErrorNoResult, res
		}
		return errmsg.ErrorSqlScanner, res
	}
	return errmsg.SUCCESS, res
}

// GetAdvisorOrderScore 获取顾问的订单的评分
func GetAdvisorOrderScore(id int64) (code int, score float32) {
	score = 0.0
	where := map[string]interface{}{
		"advisor_id":     id,
		"status":         model.Completed,
		"comment_status": model.Commented,
	}
	selects := []string{"rate"}
	code, data := GetManyTableItemsByWhere(ORDERTABLE, where, selects)
	if code != errmsg.SUCCESS {
		return
	}
	if len(data) != 0 {
		for _, v := range data {
			score += float32(v["rate"].(int64))
		}
		score /= float32(len(data))
	}
	return errmsg.SUCCESS, score
}
