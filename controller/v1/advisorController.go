package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"service/model"
	"service/service"
	"service/utils/errmsg"
	"service/utils/logger"
	"service/utils/tools"
	"service/utils/validator"
	"strconv"
)

func NewAdvisorController(ctx *gin.Context) {
	var table = service.ADVISORTABLE
	var data model.Login
	var code int
	var msg string
	// 数据绑定
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ginBindError(ctx, err, "NewUser", data)
		return
	}
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		}
		commonReturn(ctx, code, msg, data)
	}()

	// 数据校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		// 数据不合法
		return
	}
	// 调用service层 检查手机号是否重复
	if code = service.CheckPhoneExist(table, data.Phone); code == errmsg.SUCCESS {
		// 用户密码加密存储
		data.Password = service.GetPwd(data.Password)
		// 顾问的创建和服务的创建使用事务统一提交
		code, data.Id = service.NewAdvisorAndService(&data)
	}
	// success
	return
}

// TODO 修改
func UpdateAdvisorController(ctx *gin.Context) {
	var data model.Advisor
	var res map[string]interface{}
	var code int
	var msg string
	// 数据绑定
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ginBindError(ctx, err, "UpdateAdvisorController", data)
		return
	}
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		}
		commonReturn(ctx, code, msg, tools.TransformData(res))
	}()
	// 将结构体中非nil的字段提取到map中
	if res, code = tools.StructToMap(data, "structs"); code != errmsg.SUCCESS {
		return
	}
	// 数据校验 将不同的字段绑定到不同的校验函数中，使用反射做校验
	// 不存在的字段在函数中做了检验
	validateFunc := map[string]interface{}{
		"name":            validator.Name,
		"phone":           validator.Phone,
		"work_experience": validator.WorkExperience,
		"bio":             validator.Bio,
		"about":           validator.About,
	}
	for key, value := range res {
		msg, code = validator.CallFunc(validateFunc, key, value)
		if code != errmsg.SUCCESS {
			return
		}
	}
	res["id"] = ctx.GetInt64("id")
	// 检查手机号是否重复
	if res["phone"] != nil {
		code, value := service.GetTableItem(service.ADVISORTABLE, res["id"].(int64), "phone")
		// 电话号码有变动
		if fmt.Sprintf("%s", value) != *(res["phone"].(*string)) {
			code = service.CheckPhoneExist(service.ADVISORTABLE, res["phone"])
			if code != errmsg.SUCCESS {
				return
			}
		}
	}
	// 更新
	code = service.Update(service.ADVISORTABLE, res)
	return
}

func UpdateAdvisorPwd(ctx *gin.Context) {
	// 跟user修改密码共用一个接口，只不过涉及的表名不一样
	UpdatePwdController(service.ADVISORTABLE, ctx)
}

// AdvisorLogin 顾问登录
func AdvisorLogin(ctx *gin.Context) {
	Login(service.ADVISORTABLE, ctx)
}

// GetAdvisorList 获取顾问的列表
func GetAdvisorList(ctx *gin.Context) {
	pageString := ctx.Param("page")
	page, err := strconv.Atoi(pageString)
	var code int
	if err != nil || page < 1 {
		code = errmsg.ErrorInput
		commonReturn(ctx, code, "", map[string]int{"page": page})
		return
	}
	code, data := service.GetAdvisorList(page)
	commonReturn(ctx, code, "", data)
}

// GetAdvisorInfo 获取顾问的详细信息
func GetAdvisorInfo(ctx *gin.Context) {
	type Num struct {
		Id int64 `form:"id" validate:"required,min=0"`
	}
	var data Num
	var serviceData, commentsData []map[string]interface{}
	var infoData, showInfo map[string]interface{}
	var code int
	var msg string
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ginBindError(ctx, err, "GetAdvisorInfo", data)
	}
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		}
		commonReturn(ctx, code, msg,
			map[string]interface{}{
				"info":        tools.TransformData(infoData),
				"service":     tools.TransformDataSlice(serviceData),
				"showing":     tools.TransformData(showInfo),
				"commentData": tools.TransformDataSlice(commentsData),
			},
		)
	}()
	// 字段基础校验
	if msg, code := validator.Validate(data); code != errmsg.SUCCESS {
		commonReturn(ctx, code, msg, data)
		return
	}
	// 先拿顾问的info
	if code, infoData = service.GetManyTableItemsById(service.ADVISORTABLE, data.Id, []string{"*"}); code == errmsg.SUCCESS {
		// 在拿用户的服务
		code, serviceData = service.GetAdvisorService(data.Id)
	}
	delete(infoData, "password")
	infoData["coin"] = tools.ConvertCoinI2F(infoData["coin"].(int64))
	// week3 更详细的信息
	// 评分
	showInfo = map[string]interface{}{}
	if code, showInfo["score"] = service.GetAdvisorScore(data.Id); code != errmsg.SUCCESS {
		return
	}
	// 总评论数
	if code, showInfo["total_comment"] = service.GetTableItemByWhere(service.ORDERTABLE, map[string]interface{}{
		"status":         model.Completed,
		"comment_status": model.Commented,
		"advisor_id":     data.Id,
	}, "count(id)"); code != errmsg.SUCCESS {
		return
	}
	// 总订单数(readings)
	if code, showInfo["total_order_num"] = service.GetTableItemByWhere(service.ORDERTABLE, map[string]interface{}{
		"advisor_id": data.Id,
		"_or": []map[string]interface{}{
			{"service_name_id": model.VideoReading},
			{"service_name_id": model.AudioReading},
			{"service_name_id": model.TextReading},
		},
	}, "count(id)"); code != errmsg.SUCCESS {
		return
	}
	// on-time 订单完成数/总订单数
	var totalOrderCompleted, totalOrderNum interface{}
	if code, totalOrderCompleted = service.GetTableItemByWhere(service.ORDERTABLE, map[string]interface{}{
		"advisor_id": data.Id,
		"status":     model.Completed,
	}, "count(id)"); code != errmsg.SUCCESS {
		return
	}
	if code, totalOrderNum = service.GetTableItemByWhere(service.ORDERTABLE, map[string]interface{}{
		"advisor_id": data.Id,
	}, "count(id)"); code != errmsg.SUCCESS {
		return
	}
	if totalOrderNum.(int64) != 0 {
		showInfo["on-time"] = float32(totalOrderCompleted.(int64)) / float32(totalOrderNum.(int64))
	} else {
		showInfo["on-time"] = 0
	}
	// 评论数据
	code, commentsData = service.GetAdvisorCommentData(data.Id)
	return
}

// ModifyAdvisorStatus 顾问修改自己的服务状态
func ModifyAdvisorStatus(ctx *gin.Context) {
	var code int
	var msg string
	var data model.ServiceState

	if err := ctx.ShouldBind(&data); err != nil {
		ginBindError(ctx, err, "ModifyAdvisorStatus", data)
		return
	}
	data.Id = ctx.GetInt64("id")

	// return func
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		}
		commonReturn(ctx, code, msg, data)
	}()
	// 输入校验后执行业务逻辑
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 更新
	code = service.UpdateTableItem(service.ADVISORTABLE, data.Id,
		map[string]interface{}{"status": data.Status})
	return
}
