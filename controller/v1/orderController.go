package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"service/model"
	"service/service"
	"service/utils"
	"service/utils/cronjob"
	"service/utils/errmsg"
	"service/utils/logger"
	"service/utils/validator"
	"strconv"
	"time"
)

// NewOrderController 新建订单
func NewOrderController(ctx *gin.Context) {
	var jsonData model.Order
	var code int
	var msg string
	// 数据绑定
	err := ctx.ShouldBindJSON(&jsonData)
	if err != nil {
		ginBindError(ctx, err, "NewOrderController", jsonData)
		return
	}
	// 只接受这些字段，其他字段不接受。
	data := model.Order{
		UserId:        jsonData.UserId,
		ServiceId:     jsonData.ServiceId,
		AdvisorId:     jsonData.AdvisorId,
		Situation:     jsonData.Situation,
		Question:      jsonData.Question,
		Coin:          jsonData.Coin,
		Status:        model.Pending,
		CommentStatus: model.NotComment,
	}
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		} else {
			logger.Log.Info("用户新建订单", zap.Int64("order_id", data.Id))
		}
		commonReturn(ctx, code, msg, data)
	}()
	// 数据基本校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// ------- 输入数据检查 -------
	data.UserId = ctx.GetInt64("id")

	// serviceId 跟顾问的Id是否绑定正确
	var advisorIdInSQL interface{}
	if code, advisorIdInSQL = service.GetTableItem(service.SERVICETABLE, data.ServiceId, "advisor_id"); code != errmsg.SUCCESS {
		return
	}
	if advisorIdInSQL.(int64) != data.AdvisorId {
		code = errmsg.ErrorServiceIdNotMatchWithAdvisorID
		return
	}
	// serviceId 是否还是open的
	var serviceStatus interface{}
	if code, serviceStatus = service.GetTableItem(service.SERVICETABLE, data.ServiceId, "status"); code != errmsg.SUCCESS {
		return
	}
	if serviceStatus.(int64) != 1 {
		code = errmsg.ErrorServiceNotOpen
		return
	}
	// coin如果==0或者前端没传 自己查 如果传了但是不一致，返回可能顾问可能修改了价格
	var coinInSQL interface{}
	if code, coinInSQL = service.GetTableItem(service.SERVICETABLE, data.ServiceId, "price"); code != errmsg.SUCCESS {
		return
	}
	if data.Coin == 0 {
		data.Coin = coinInSQL.(float32)
	} else if data.Coin != coinInSQL || data.Coin < -1 {
		code = errmsg.ErrorPriceNotMatch
		return
	}

	// ------- 输入数据检查结束 -------
	data.Status = 0
	data.CreateTime = time.Now().Unix()
	// 加急订单的价格 只做记录，等到用户加急的时候安装这个去扣钱
	data.RushCoin = data.Coin * utils.RushOrderCost
	// 提交到service层的事务
	if code, data.Id = service.NewOrderAndCostTrans(&data); code != errmsg.SUCCESS {
		return
	}

	// 订单状态24h后过期
	job := cronjob.CronJob{
		OrderId:    data.Id,
		UserId:     data.UserId,
		CreateTime: data.CreateTime,
		CronId:     -1,
		CronType:   cronjob.PendingOrderType,
	}
	if code = cronjob.AddJob(&job); code != errmsg.SUCCESS {
		logger.Log.Error("用户订单的定时任务创建失败", zap.Int64("order_id", data.Id))
		return
	}
	return
}

// GetOrderListController 获取顾问的订单列表
func GetOrderListController(ctx *gin.Context) {
	code, data := service.GetOrderList(ctx.GetInt64("id"))
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		} else {
			logger.Log.Info("顾问查看订单列表", zap.Int64("advisor_id", ctx.GetInt64("id")))
		}
		commonReturn(ctx, code, "", TransformDataSlice(data))
	}()
	if code == errmsg.SUCCESS {
		// 附加信息 用户名、时间格式、服务类型、
		for _, v := range data {
			code, userNameUint8 := service.GetTableItem(service.USERTABLE, v["user_id"].(int64), "name")
			if code != errmsg.SUCCESS {
				return
			}
			v["user_name"] = fmt.Sprintf("%s", userNameUint8)
			v["show_time"] = time.Unix(v["create_time"].(int64), 0).Format("Jan 02,2006")
			code, v["service_name_id"] = service.GetTableItem(service.SERVICETABLE, v["service_id"].(int64), "service_name_id")
			if code != errmsg.SUCCESS {
				return
			}
			code, v["service_name"] = model.GetServiceNameById(int(v["service_name_id"].(int64)))
			if code != errmsg.SUCCESS {
				return
			}
			v["status_name"] = model.GetOrderStatusNameById(int(v["status"].(int64)))
		}
	}
	return
}

// GetOrderDetailController 获取订单详情
func GetOrderDetailController(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.Atoi(idString)
	var code int
	var msg string
	var data map[string]interface{}
	// return
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		} else {
			logger.Log.Info("顾问查看订单详情", zap.Int("order_id", id))
		}
		commonReturn(ctx, code, msg, TransformData(data))
	}()
	if err != nil || id < 1 {
		code = errmsg.ErrorInput
		data = map[string]interface{}{"id": idString}
		return
	}
	// 是不是你的订单
	code, advisorIdInSQL := service.GetTableItem(service.ORDERTABLE, int64(id), "advisor_id")
	if code != errmsg.SUCCESS {
		return
	}
	if advisorIdInSQL.(int64) != ctx.GetInt64("id") {
		code = errmsg.ErrorOrderIdNotMatchWithAdvisorID
		return
	}
	// 获取基础的订单信息
	if code, data = service.GetManyTableItemsById(service.ORDERTABLE, int64(id), []string{"*"}); code != errmsg.SUCCESS {
		return
	}
	//在基础的信息上扩充用户的姓名、出生日期、和性别等相关信息
	code, userInfo := service.GetUserInfo(data["user_id"].(int64))
	// 转化生日格式
	if code == errmsg.SUCCESS {
		birth := userInfo["birth"].(string)
		birthTime, err := time.Parse("02-01-2006", birth)
		if err == nil {
			userInfo["birthShow"] = birthTime.Format("Jan 02,2006")
		}
	}
	data["userInfo"] = TransformData(userInfo)
	return
}

// OrderReplyController 顾问回复订单
func OrderReplyController(ctx *gin.Context) {

	var data model.OrderReply
	var code int
	var msg string
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ginBindError(ctx, err, "orderController.orderReplyController", data)
		return
	}
	//基础校验 回复长度
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		commonReturn(ctx, code, msg, data)
		return
	}

	// return
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		} else {
			logger.Log.Info("顾问回复订单", zap.Int64("order_id", data.Id))
		}
		commonReturn(ctx, code, "", data)
	}()

	data.AdvisorId = ctx.GetInt64("id")
	// 检查顾问的ID和库里的订单上的id是否一致
	code, advisorIdInSQL := service.GetTableItem(service.ORDERTABLE, data.Id, "advisor_id")
	if code != errmsg.SUCCESS || advisorIdInSQL.(int64) != data.AdvisorId {
		code = errmsg.ErrorOrderIdNotMatchWithAdvisorID
		return
	}
	//// 检测订单是什么状态 只有pending,rush可以回复 放到事务里去了。
	//获取订单创建时的金币价格
	code, coin := service.GetTableItem(service.ORDERTABLE, data.Id, "coin")
	if code != errmsg.SUCCESS {
		return
	}
	data.Coin = coin.(float32)
	//加急的订单价格
	code, rushCoin := service.GetTableItem(service.ORDERTABLE, data.Id, "rush_coin")
	if code == errmsg.SUCCESS {
		data.RushCoin = rushCoin.(float32)
		// 提交到service层
		code = service.ReplyOrderServiceTrans(&data)
	}
	return
}

// RushOrderController 订单加急
func RushOrderController(ctx *gin.Context) {

	var data model.OrderRush
	var code int
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ginBindError(ctx, err, "RushOrderController", data)
		return
	}
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		} else {
			logger.Log.Info("用户加急订单", zap.Int64("order_id", data.Id))
		}
		commonReturn(ctx, code, "", data)
	}()

	data.UserId = ctx.GetInt64("id")
	// 是不是自己的订单
	code, userIdInSQL := service.GetTableItem(service.ORDERTABLE, data.Id, "user_id")
	if code != errmsg.SUCCESS || userIdInSQL.(int64) != ctx.GetInt64("id") {
		code = errmsg.ErrorOrderIdNotMatchWithUserID
		return
	}
	// 最后一个小时不能加急了
	code, orderCreateTime := service.GetTableItem(service.ORDERTABLE, data.Id, "create_time")
	if orderCreateTime.(int64)-time.Now().Unix() > 23*60*60 {
		code = errmsg.ErrorOrderCantRush
		return
	}
	data.RushTime = time.Now().Unix()
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

// CommentOrderController 用户回复订单 week3
func CommentOrderController(ctx *gin.Context) {
	var comment model.CommentStruct
	var data model.OrderComment
	var code int
	var msg string
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ginBindError(ctx, err, "CommentOrderController", comment)
		return
	}
	// defer return
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		} else {
			logger.Log.Info("用户评论订单", zap.Int64("order_id", data.Id))
		}
		commonReturn(ctx, code, msg, data)
	}()
	// 数据基本校验
	msg, code = validator.Validate(comment)
	if code != errmsg.SUCCESS {
		return
	}
	// 构造
	data = model.OrderComment{
		CommentStruct: comment,
		UserId:        ctx.GetInt64("id"),
		CommentTime:   time.Now().Unix(),
	}
	// 检查订单与用户ID是否对应
	var userIdInSQL interface{}
	if code, userIdInSQL = service.GetTableItem(service.ORDERTABLE, data.Id, "user_id"); code != errmsg.SUCCESS {
		return
	}
	if userIdInSQL.(int64) != data.UserId {
		code = errmsg.ErrorOrderIdNotMatchWithUserID
		return
	}
	// 订单完成才能回复
	var orderStatus interface{}
	if code, orderStatus = service.GetTableItem(service.ORDERTABLE, data.Id, "status"); code != errmsg.SUCCESS {
		return
	}
	if orderStatus.(int64) != model.Completed {
		// 没完成你回复啥
		code = errmsg.ErrorOrderCantComment
		return
	}
	// 订单是否已经回复过一次
	var commentStatus interface{}
	if code, commentStatus = service.GetTableItem(service.ORDERTABLE, data.Id, "comment_status"); code != errmsg.SUCCESS {
		return
	}
	if commentStatus.(int64) != model.NotComment {
		// 订单不可以在回复了
		code = errmsg.ErrorOrderCantComment
		return
	}
	// 更新数据
	code = service.UpdateTableItem(service.ORDERTABLE, data.Id,
		map[string]interface{}{
			"comment_time": data.CommentTime,
			"comment":      data.Comment,
		},
	)
	if code == errmsg.SUCCESS {
		// 更新评论状态
		code = service.UpdateTableItem(service.ORDERTABLE, data.Id, map[string]interface{}{
			"comment_status": model.Commented,
		})
		if code != errmsg.SUCCESS {
			logger.Log.Warn("用户的评论状态更新失败，实际已评论!", zap.Int64("order_id", data.Id))
		}
	}
	return
}
