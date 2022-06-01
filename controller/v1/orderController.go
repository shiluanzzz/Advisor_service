package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"service/model"
	"service/service"
	"service/utils/cronjob"
	"service/utils/errmsg"
	"service/utils/validator"
	"strconv"
	"time"
)

func NewOrderController(ctx *gin.Context) {
	var data model.Order
	// 数据绑定
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ginBindError(ctx, err, "NewOrderController", data)
		return
	}
	// 数据基本校验
	msg, code := validator.Validate(data)
	if code != errmsg.SUCCESS {
		commonReturn(ctx, code, msg, data)
		return
	}
	// ------- 输入数据检查 -------
	// token role角色校验
	if ctx.GetString("role") != service.USERTABLE {
		commonReturn(ctx, errmsg.ErrorTokenRoleNotMatch, ctx.GetString("role"), data)
		return
	}
	// user_id 跟 token里的id是否一致
	if data.UserId == 0 {
		data.UserId = ctx.GetInt64("id")
	} else if data.UserId != ctx.GetInt64("id") {
		commonReturn(ctx, errmsg.ErrorIdNotMatchWithToken, "", data)
		return
	}
	// serviceId 跟顾问的Id是否绑定正确
	code, advisorIdInSQL := service.GetTableItem(service.SERVICETABLE, data.ServiceId, "advisor_id")
	if code != errmsg.SUCCESS {
		commonReturn(ctx, code, "", data)
		return
	}
	if advisorIdInSQL.(int64) != data.AdvisorId {
		commonReturn(ctx, errmsg.ErrorServiceIdNotMatchWithAdvisorID, "", data)
		return
	}
	// serviceId 是否还是open的
	code, serviceStatus := service.GetTableItem(service.SERVICETABLE, data.ServiceId, "status")
	if code != errmsg.SUCCESS {
		commonReturn(ctx, code, "", data)
		return
	}
	if serviceStatus.(int64) != 1 {
		commonReturn(ctx, errmsg.ErrorServiceNotOpen, "", data)
		return
	}
	// coin如果==0或者前端没传 自己查 如果传了但是不一致，返回可能顾问可能修改了价格
	code, coinInSQL := service.GetTableItem(service.SERVICETABLE, data.ServiceId, "price")
	if data.Coin == 0 {
		data.Coin = coinInSQL.(float32)
	} else if data.Coin != coinInSQL || data.Coin < -1 {
		commonReturn(ctx, errmsg.ErrorPriceNotMatch, "", data)
		return
	}
	// 钱够吗你? 放到订单的创建事务中去了

	// ------- 输入数据检查 -------
	data.Status = 0
	data.CreateTime = time.Now().Unix()
	// 加急订单的价格 只做记录
	data.RushCoin = data.Coin / 2
	code, data.Id = service.NewOrderAndCostTrans(&data)
	//if code==errmsg.SUCCESS{
	//	// TODO 订单状态24h后过期
	//}
	commonReturn(ctx, code, "", data)
	return
}
func GetOrderListController(ctx *gin.Context) {
	if ctx.GetString("role") != service.ADVISORTABLE {
		commonReturn(ctx, errmsg.ErrorTokenRoleNotMatch, "", ctx.GetString("role"))
		return
	}
	code, data := service.GetOrderList(ctx.GetInt64("id"))
	if code != errmsg.SUCCESS {
		commonReturn(ctx, code, "", data)
	}
	// 附加信息 用户名、时间格式、服务类型、
	for _, v := range data {
		_, userNameUint8 := service.GetTableItem(service.USERTABLE, v["user_id"].(int64), "name")
		v["user_name"] = fmt.Sprintf("%s", userNameUint8)
		v["show_time"] = time.Unix(v["create_time"].(int64), 0).Format("Jan 02,2006")
		_, v["service_name_id"] = service.GetTableItem(service.SERVICETABLE, v["service_id"].(int64), "service_name_id")
		_, v["service_name"] = model.GetServiceNameById(int(v["service_name_id"].(int64)))
		v["status_name"] = model.GetStatusNameByCode(int(v["status"].(int64)))
	}
	commonReturn(ctx, code, "", TransformDataSlice(data))
	return
}
func GetOrderDetailController(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil || id < 1 {
		commonReturn(ctx, errmsg.ErrorInput, "", map[string]string{"id": idString})
		return
	}
	// 获取基础的订单信息
	code, base := service.GetOrderInfo(id)
	if code != errmsg.SUCCESS {
		commonReturn(ctx, code, "", base)
		return
	}
	//在基础的信息上扩充用户的姓名、出生日期、和性别等相关信息
	code, userInfo := service.GetUserInfo(base["user_id"].(int64))
	// 转化生日格式
	if code == errmsg.SUCCESS {
		birth := userInfo["birth"].(string)
		if birth != "" {
			birthTime, err := time.Parse("02-01-2006", birth)
			if err == nil {
				userInfo["birthShow"] = birthTime.Format("Jan 02,2006")
			}
		}
	}
	base["userInfo"] = TransformData(userInfo)
	commonReturn(ctx, code, "", TransformData(base))
}
func OrderReplyController(ctx *gin.Context) {
	// token角色
	if ctx.GetString("role") != service.ADVISORTABLE {
		commonReturn(ctx, errmsg.ErrorTokenRoleNotMatch, "", nil)
		return
	}
	var data model.OrderReply
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ginBindError(ctx, err, "orderController.orderReplyController", data)
		return
	}
	data.AdvisorId = ctx.GetInt64("id")
	// 检查顾问的ID和库里的订单上的id是否一致
	code, advisorIdInSQL := service.GetTableItem(service.ORDERTABLE, data.Id, "advisor_id")
	if code != errmsg.SUCCESS || advisorIdInSQL.(int64) != data.AdvisorId {
		commonReturn(ctx, errmsg.ErrorServiceIdNotMatchWithAdvisorID, "", data)
		return
	}
	//// 检测订单是什么状态 只有pending,rush可以回复 放到事务里去了。

	//基础校验 回复长度
	msg, code := validator.Validate(data)
	if code != errmsg.SUCCESS {
		commonReturn(ctx, code, msg, data)
		return
	}
	//获取订单创建时的金币价格
	code, coin := service.GetTableItem(service.ORDERTABLE, data.Id, "coin")
	if code == errmsg.SUCCESS {
		data.Coin = coin.(float32)
	} else {
		commonReturn(ctx, code, "", data)
	}
	//加急的订单价格
	code, rushCoin := service.GetTableItem(service.ORDERTABLE, data.Id, "rush_coin")
	if code == errmsg.SUCCESS {
		data.RushCoin = rushCoin.(float32)
		// 提交到service层
		code = service.ReplyOrderServiceTrans(&data)
	}
	commonReturn(ctx, code, "", data)
	return
}

// RushOrderController 订单加急
func RushOrderController(ctx *gin.Context) {

	var data model.OrderRush
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ginBindError(ctx, err, "RushOrderController", data)
		return
	}
	// token role角色校验
	if ctx.GetString("role") != service.USERTABLE {
		commonReturn(ctx, errmsg.ErrorTokenRoleNotMatch, ctx.GetString("role"), data)
		return
	}
	data.UserId = ctx.GetInt64("id")
	// 是不是自己的订单
	code, userIdInSQL := service.GetTableItem(service.ORDERTABLE, data.Id, "user_id")
	if code != errmsg.SUCCESS || userIdInSQL.(int64) != ctx.GetInt64("id") {
		commonReturn(ctx, errmsg.ErrorServiceIdNotMatchWithAdvisorID, "", data)
	}
	data.RushTime = time.Now().Unix()
	// 新建一个cron的定时job
	job := cronjob.CronJob{
		OrderId:  data.Id,
		UserId:   data.UserId,
		RushTime: data.RushTime,
		CronId:   -1,
	}
	code = cronjob.AddJob(&job)
	if code == errmsg.SUCCESS {
		// 提交到service的事务层
		code = service.RushOrderTrans(&data)
		if code != errmsg.SUCCESS {
			// 如果事务失败,停止定时job
			cronjob.CloseJob(&job)
		}
	}
	commonReturn(ctx, code, "", data)
}
