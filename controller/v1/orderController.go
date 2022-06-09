package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"service-backend/model"
	"service-backend/service"
	"service-backend/utils/cronjob"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
	"service-backend/utils/setting"
	"service-backend/utils/validator"
	"time"
)

// NewOrderController 新建订单
func NewOrderController(ctx *gin.Context) {
	var orderInfo model.OrderInitInfo
	var response model.Order
	var code int
	var msg string
	// 数据绑定
	err := ctx.ShouldBindJSON(&orderInfo)
	if err != nil {
		ginBindError(ctx, err, orderInfo)
		return
	}
	// init model
	response = model.Order{
		UserId:        orderInfo.UserId,
		ServiceId:     orderInfo.ServiceId,
		AdvisorId:     orderInfo.AdvisorId,
		Situation:     orderInfo.Situation,
		Question:      orderInfo.Question,
		Status:        model.Pending,
		CommentStatus: model.NotComment,
	}
	defer func() {
		orderInfo.OrderId = response.Id
		orderInfo.UserId = response.UserId
		logger.CommonControllerLog(&code, &msg, orderInfo, response)
		commonReturn(ctx, code, msg, orderInfo)
	}()
	// 数据基本校验
	if msg, code = validator.Validate(orderInfo); code != errmsg.SUCCESS {
		return
	}
	// 取sql数据做校验
	response.UserId = ctx.GetInt64("id")
	var serviceInSQL model.Service
	if code, serviceInSQL = service.GetService(response.ServiceId); code != errmsg.SUCCESS {
		return
	}
	var UserInSQL model.User
	if code, UserInSQL = service.GetUser(response.UserId); code != errmsg.SUCCESS {
		return
	}
	// ------- 输入数据检查 -------
	// serviceId 跟顾问的Id是否绑定正确
	if serviceInSQL.AdvisorId != response.AdvisorId {
		code = errmsg.ErrorServiceIdNotMatchWithAdvisorID
		return
	}
	// serviceId 是否还是open的
	if serviceInSQL.Status != model.AdvisorServiceOpen {
		code = errmsg.ErrorServiceNotOpen
		return
	}
	// 金币检查
	if serviceInSQL.Price > UserInSQL.Coin {
		code = errmsg.ErrorOrderMoneyInsufficient
		return
	}
	// ------- 输入数据检查结束 -------
	response.Coin = serviceInSQL.Price
	response.CreateTime = time.Now().Unix()
	response.ServiceNameId = serviceInSQL.ServiceNameId
	// 加急订单的价格 只做记录，等到用户加急的时候安装这个去扣钱
	response.RushCoin = int64(float32(response.Coin) * setting.ServiceCfg.RushOrderCost)

	// 提交到service层的事务
	if code, response.Id = service.NewOrderAndCostTrans(&response); code != errmsg.SUCCESS {
		return
	}
	// 新建订单后更新顾问信息的指标
	if code = service.UpdateAdvisorIndicators(response.AdvisorId); code != errmsg.SUCCESS {
		return
	}

	// 订单状态24h后过期 新建一个监控事务
	job := cronjob.CronJob{
		OrderId:    response.Id,
		UserId:     response.UserId,
		CreateTime: response.CreateTime,
		CronId:     -1,
		CronType:   cronjob.PendingOrderType,
	}
	if code = cronjob.AddJob(&job); code != errmsg.SUCCESS {
		logger.Log.Error("用户订单的定时任务创建失败", zap.Int64("order_id", response.Id))
		return
	}
	return
}

// GetOrderListController 获取顾问的订单列表
func GetOrderListController(ctx *gin.Context) {
	var response []*model.Order
	var code int
	var msg string
	defer func() {
		logger.CommonControllerLog(&code, &msg, ctx.GetInt64("id"), response)
		commonReturn(ctx, code, "", response)
	}()
	code, response = service.GetAdvisorOrderList(ctx.GetInt64("id"))
	return
}

// GetOrderDetailController 获取订单详情
func GetOrderDetailController(ctx *gin.Context) {
	var request model.TableID
	var order model.Order
	var user model.User
	response := map[string]interface{}{
		"orderInfo": &order,
		"userInfo":  &user,
	}
	var code int
	var msg string
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ginBindError(ctx, err, request)
	}
	// return
	defer func() {
		logger.CommonControllerLog(&code, &msg, request, response)
		commonReturn(ctx, code, "", response)
	}()
	if msg, code = validator.Validate(request); code != errmsg.SUCCESS {
		return
	}
	// 逻辑校验 直接拿数据然后
	if code, order = service.GetOrder(request.Id); code != errmsg.SUCCESS {
		return
	}
	// 订单是不是你的
	if order.AdvisorId != ctx.GetInt64("id") {
		code = errmsg.ErrorOrderIdNotMatchWithAdvisorID
		return
	}
	//在基础的信息上扩充用户的相关信息
	if code, user = service.GetUser(order.UserId); code != errmsg.SUCCESS {
		return
	}
	// 业务修正
	user.Coin = 0
	user.CoinShow = 0.0
	user.UpdateShow("Jan 02,2006")

	return
}

// OrderReplyController 顾问回复订单
func OrderReplyController(ctx *gin.Context) {

	var data model.OrderReply
	var response model.OrderReply
	var code int
	var msg string
	// 数据绑定
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ginBindError(ctx, err, data)
		return
	}
	//基础校验 回复长度
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		commonReturn(ctx, code, msg, data)
		return
	}
	// return
	defer func() {
		logger.CommonControllerLog(&code, &msg, data, response)
		commonReturn(ctx, code, "", response)
	}()
	// 逻辑校验+service层提交
	code, response = func() (code int, response model.OrderReply) {

		data.AdvisorId = ctx.GetInt64("id")
		var orderInSql model.Order
		if code, orderInSql = service.GetOrder(data.Id); code != errmsg.SUCCESS {
			return code, data
		}
		// 检查顾问的ID和库里的订单上的id是否一致
		if orderInSql.AdvisorId != data.AdvisorId {
			return errmsg.ErrorOrderIdNotMatchWithAdvisorID, data
		}
		//检测订单是什么状态 只有pending,rush可以回复
		if !orderInSql.Status.CanReply() {
			return errmsg.ErrorOrderHasCompleted, data
		}

		// 校验全部通过后初始化response
		response = model.OrderReply{
			Id:        data.Id,
			AdvisorId: data.AdvisorId,
			Reply:     data.Reply,
			Coin:      orderInSql.Coin,
			RushCoin:  orderInSql.RushCoin,
			Status:    orderInSql.Status,
		}
		// 提交到service层
		if code = service.ReplyOrderServiceTrans(&response); code != errmsg.SUCCESS {
			return
		}
		// 回复订单后更新指标
		if code = service.UpdateAdvisorIndicators(response.AdvisorId); code != errmsg.SUCCESS {
			return
		}
		return
	}()
	return
}

// RushOrderController 订单加急
func RushOrderController(ctx *gin.Context) {

	var data model.OrderRush
	var code int
	var msg string
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ginBindError(ctx, err, data)
		return
	}
	defer func() {
		logger.CommonControllerLog(&code, &msg, data, data)
		commonReturn(ctx, code, "", data)
	}()

	data.UserId = ctx.GetInt64("id")
	var orderInSql model.Order
	var userMoney interface{}
	if code, orderInSql = service.GetOrder(data.Id); code != errmsg.SUCCESS {
		return
	}
	if code, userMoney = service.GetTableItem(service.USERTABLE, data.UserId, "coin"); code != errmsg.SUCCESS {
		return
	}
	/*  ---  逻辑校验  ---  */
	// 是不是自己的订单
	if orderInSql.UserId != data.UserId {
		code = errmsg.ErrorOrderIdNotMatchWithUserID
		return
	}
	// rush和expired下不能加急
	if !orderInSql.Status.CanRush() {
		code = errmsg.ErrorOrderCantRush
		return
	}
	// 最后一个小时不能加急了
	if orderInSql.CreateTime-time.Now().Unix() > 23*60*60 {
		code = errmsg.ErrorOrderCantRush
		return
	}
	// 钱够不够
	if orderInSql.RushCoin > userMoney.(int64) {
		code = errmsg.ErrorOrderMoneyInsufficient
		return
	}
	/*  ---  逻辑校验  ---  */
	data.RushTime = time.Now().Unix()
	data.RushMoney = orderInSql.RushCoin
	data.UserMoney = userMoney.(int64)
	// 新建一个cron的定时job
	job := cronjob.CronJob{
		OrderId:  data.Id,
		UserId:   data.UserId,
		RushTime: data.RushTime,
		CronType: cronjob.RushOrderType,
	}
	if code = cronjob.AddJob(&job); code != errmsg.SUCCESS {
		return
	}
	// 提交到service的事务层
	if code = service.RushOrderTrans(&data); code != errmsg.SUCCESS {
		// 如果事务失败,停止定时job
		cronjob.CloseJob(&job)
	}
	return
}

// CommentOrderController 用户评论订单 week3
func CommentOrderController(ctx *gin.Context) {
	var comment model.CommentStruct
	var data model.OrderComment
	var code int
	var msg string
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ginBindError(ctx, err, comment)
		return
	}
	// defer return
	defer func() {
		logger.CommonControllerLog(&code, &msg, comment, data)
		commonReturn(ctx, code, msg, data)
	}()
	// 数据基本校验
	msg, code = validator.Validate(comment)
	if code != errmsg.SUCCESS {
		return
	}
	// 构造
	data = model.OrderComment{
		Id:      comment.Id,
		Comment: comment.Comment,
		Rate:    comment.Rate,
		//CommentStruct: comment,
		UserId:      ctx.GetInt64("id"),
		CommentTime: time.Now().Unix(),
	}
	/*  ---  逻辑校验  ---  */
	var orderInSql model.Order
	if code, orderInSql = service.GetOrder(data.Id); code != errmsg.SUCCESS {
		return
	}
	// 检查订单与用户ID是否对应
	if orderInSql.UserId != data.UserId {
		code = errmsg.ErrorOrderIdNotMatchWithUserID
		return
	}
	// 订单完成才能回复
	if orderInSql.Status != model.Completed {
		code = errmsg.ErrorOrderCantComment
		return
	}
	// 订单是否已经回复过一次
	if orderInSql.CommentStatus != model.NotComment {
		code = errmsg.ErrorOrderCantComment
		return
	}
	/*  ---  逻辑校验  ---  */

	// 更新数据
	newData := map[string]interface{}{
		"comment_time":   data.CommentTime,
		"comment":        data.Comment,
		"rate":           data.Rate,
		"comment_status": model.Commented,
	}
	if code = service.UpdateTableItemById(service.ORDERTABLE, data.Id, newData); code != errmsg.SUCCESS {
		return
	}
	// 评论订单后更新指标
	if code = service.UpdateAdvisorIndicators(orderInSql.AdvisorId); code != errmsg.SUCCESS {
		return
	}
	return
}
