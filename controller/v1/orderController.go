package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
		Id:        jsonData.Id,
		UserId:    jsonData.UserId,
		ServiceId: jsonData.ServiceId,
		AdvisorId: jsonData.AdvisorId,
		Situation: jsonData.Situation,
		Question:  jsonData.Question,
		Coin:      jsonData.Coin,
		Status:    model.Pending,
	}
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		}
		commonReturn(ctx, code, msg, data)
	}()
	// 数据基本校验
	msg, code = validator.Validate(data)
	if code != errmsg.SUCCESS {
		return
	}
	// ------- 输入数据检查 -------
	// user_id 跟 token里的id是否一致
	if data.UserId == 0 {
		data.UserId = ctx.GetInt64("id")
	} else if data.UserId != ctx.GetInt64("id") {
		code = errmsg.ErrorIdNotMatchWithToken
		return
	}
	// serviceId 跟顾问的Id是否绑定正确
	code, advisorIdInSQL := service.GetTableItem(service.SERVICETABLE, data.ServiceId, "advisor_id")
	if code != errmsg.SUCCESS {
		return
	}
	if advisorIdInSQL.(int64) != data.AdvisorId {
		code = errmsg.ErrorServiceIdNotMatchWithAdvisorID
		return
	}
	// serviceId 是否还是open的
	code, serviceStatus := service.GetTableItem(service.SERVICETABLE, data.ServiceId, "status")
	if code != errmsg.SUCCESS {
		return
	}
	if serviceStatus.(int64) != 1 {
		code = errmsg.ErrorServiceNotOpen
		return
	}
	// coin如果==0或者前端没传 自己查 如果传了但是不一致，返回可能顾问可能修改了价格
	code, coinInSQL := service.GetTableItem(service.SERVICETABLE, data.ServiceId, "price")
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
	// 订单状态24h后过期
	job := cronjob.CronJob{
		OrderId:    data.Id,
		UserId:     data.UserId,
		CreateTime: data.CreateTime,
		CronId:     -1,
		CronType:   cronjob.PendingOrderType,
	}
	code = cronjob.AddJob(&job)
	if code == errmsg.SUCCESS {
		code, data.Id = service.NewOrderAndCostTrans(&data)
		if code != errmsg.SUCCESS {
			cronjob.CloseJob(&job)
		}
	}
	return
}

func GetOrderListController(ctx *gin.Context) {
	code, data := service.GetOrderList(ctx.GetInt64("id"))
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
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
		code = errmsg.ErrorServiceIdNotMatchWithAdvisorID
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
		}
		commonReturn(ctx, code, "", data)
	}()

	data.AdvisorId = ctx.GetInt64("id")
	// 检查顾问的ID和库里的订单上的id是否一致
	code, advisorIdInSQL := service.GetTableItem(service.ORDERTABLE, data.Id, "advisor_id")
	if code != errmsg.SUCCESS || advisorIdInSQL.(int64) != data.AdvisorId {
		code = errmsg.ErrorServiceIdNotMatchWithAdvisorID
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
		}
		commonReturn(ctx, code, "", data)
	}()

	data.UserId = ctx.GetInt64("id")
	// 是不是自己的订单
	code, userIdInSQL := service.GetTableItem(service.ORDERTABLE, data.Id, "user_id")
	if code != errmsg.SUCCESS || userIdInSQL.(int64) != ctx.GetInt64("id") {
		code = errmsg.ErrorServiceIdNotMatchWithAdvisorID
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
