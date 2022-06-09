package service

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"service-backend/model"
	"service-backend/utils"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
	"service-backend/utils/tools"
	"time"
)

var ORDERTABLE = "user_order"

func NewOrderAndCostTrans(data *model.Order) (code int, id int64) {
	Tran, err := utils.DbConn.Begin()
	defer CommonTranDefer(&code, Tran)
	if err != nil {
		return errmsg.ErrorSqlTransError, -1
	}

	// 扣掉金币
	if code = CostUserCoin(data, Tran); code != errmsg.SUCCESS {
		return code, -1
	}
	// 新建订单
	if code, id = NewOrder(data, Tran); code != errmsg.SUCCESS {
		return code, -1
	}
	// 添加流水
	bill := model.Bill{
		OrderId: id,
		UserId:  data.UserId,
		Amount:  data.Coin,
		Type:    model.ORDERCOST,
	}
	if code = NewBill(&bill, Tran); code != errmsg.SUCCESS {
		return
	}
	// 提交事务
	if err := Tran.Commit(); err != nil {
		return errmsg.ErrorSqlTransCommitError, -1
	}
	return code, id
}

// NewOrder 新建订单
func NewOrder(model *model.Order, tx *sql.Tx) (code int, id int64) {
	// 转化数据并生成sql语句AAA
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
func CostUserCoin(model *model.Order, tx *sql.Tx) (code int) {
	cond := "update `user` set coin = coin - ? where id = ?"
	code, _ = SQLExec(cond, []interface{}{model.Coin, model.UserId}, tx)
	return
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
	// 用户新增流水
	bill := model.Bill{
		OrderId:   data.Id,
		AdvisorId: data.AdvisorId,
		Amount:    reward,
		Type:      model.ORDERINCOME,
	}
	if code = NewBill(&bill, begin); code != errmsg.SUCCESS {
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
func AddCoin2Advisor(data *model.OrderReply, tx *sql.Tx) (code int) {
	cond := fmt.Sprintf("update %s set coin=coin + ? where id= ?", ADVISORTABLE)
	code, _ = SQLExec(cond, []interface{}{data.Coin, data.AdvisorId}, tx)
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
	//提交流水
	bill := model.Bill{
		UserId:  data.UserId,
		OrderId: data.Id,
		Amount:  data.RushMoney,
		Type:    model.ORDERRUSHCOST,
	}
	// 加急订单支出
	if code = NewBill(&bill, begin); code != errmsg.SUCCESS {
		return
	}

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
	originCoin := userMoney.(int64)
	var backCoin int64
	if originStatus == model.Rush && targetStatus == model.Pending {
		// 加急到普通
		backCoin += orderInSql.RushCoin
	} else if originStatus == model.Pending && targetStatus == model.Expired {
		// 普通到过期
		backCoin += orderInSql.Coin
		// 如果订单是在加急的状态下，也要把钱退回去。
		if orderInSql.Status == model.Rush {
			backCoin += orderInSql.RushCoin
		}
	} else {
		code = errmsg.ErrorJobStatusConvert
		return code
	}

	// 提交用户的金币修改
	if code = UpdateTableItemById(USERTABLE, userId, map[string]interface{}{
		"coin": originCoin + backCoin,
	}, Tran); code != errmsg.SUCCESS {
		return code
	}
	// 金币修改提交到流水表
	bill := model.Bill{
		UserId:  userId,
		OrderId: orderId,
		Amount:  backCoin,
	}
	switch originStatus {
	case model.Rush:
		bill.Type = model.ORDERRUSABACK
	case model.Pending:
		bill.Type = model.ORDERBACK
	default:
		return errmsg.ErrorJobStatusNotExpect
	}
	// 将用户的金币流水加入到账单中
	if code = NewBill(&bill, Tran); code != errmsg.SUCCESS {
		return
	}
	// 提交状态修改
	if code = UpdateTableItemById(ORDERTABLE, orderId, map[string]interface{}{
		"status": targetStatus,
	}, Tran); code != errmsg.SUCCESS {
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
	if code = GetTableRows2StructByWhere(ORDERTABLE, where, selects, &res); code != errmsg.SUCCESS {
		return
	}
	// 附加信息:用户名、时间格式、服务类型
	for _, v := range res {
		v.ShowTime = time.Unix(v.CreateTime, 0).Format("Jan 02,2006")
		_, v.UserName = GetUserName(v.UserId)
		v.ServiceName = model.ServiceKind[v.ServiceNameId]
		v.ServiceStatusName = v.Status.StatusName()
	}
	return code, res
}
func GetOrder(orderId int64) (code int, res model.Order) {
	where := map[string]interface{}{
		"id": orderId,
	}
	if code = GetTableRows2StructByWhere(ORDERTABLE, where, []string{"*"}, &res); code != errmsg.SUCCESS {
		return
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
