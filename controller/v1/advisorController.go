package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"service-backend/model"
	"service-backend/service"
	"service-backend/utils/errmsg"
	"service-backend/utils/tools"
	"service-backend/utils/validator"
	"strconv"
)

func NewAdvisorController(ctx *gin.Context) {
	var table = service.ADVISORTABLE
	var data model.Login
	var code int
	var msg string
	// 数据绑定
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ginBindError(ctx, err, data)
		return
	}
	defer func() {
		commonControllerLog(&code, &msg, data, data)
		commonReturn(ctx, code, msg, data)
	}()
	// 数据校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 调用service层 检查手机号是否重复
	code, _ = service.GetTableItemByWhere(table, map[string]interface{}{
		"phone": data.Phone,
	}, "phone")
	if code != errmsg.ErrorMysqlNoRows {
		code = errmsg.ErrorUserPhoneUsed
		return
	}

	if code = service.CheckPhoneExist(table, data.Phone); code != errmsg.SUCCESS {
		return
	}
	// 用户密码加密存储
	data.Password = service.GetPwd(data.Password)
	// 顾问的创建和服务的创建使用事务统一提交
	code, data.Id = service.NewAdvisorAndService(&data)
	return
}

func UpdateAdvisorController(ctx *gin.Context) {
	var data model.Advisor
	var res map[string]interface{}
	var code int
	var msg string
	// 数据绑定
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ginBindError(ctx, err, data)
		return
	}
	defer func() {
		commonControllerLog(&code, &msg, data, res)
		commonReturn(ctx, code, msg, res)
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
		if msg, code = validator.CallFunc(validateFunc, key, value); code != errmsg.SUCCESS {
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
	code = service.UpdateTableItemById(service.ADVISORTABLE, ctx.GetInt64("id"), res)
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

	var data model.TableID
	var serviceData, commentsData []map[string]interface{}
	var infoData, showInfo map[string]interface{}
	var code int
	var msg string
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ginBindError(ctx, err, data)
	}
	defer func() {
		commonControllerLog(&code, &msg, data, data)
		commonReturn(ctx, code, msg,
			// TODO
			map[string]interface{}{
				"info":        tools.TransformData(infoData),
				"service":     tools.TransformDataSlice(serviceData),
				"showing":     tools.TransformData(showInfo),
				"commentData": tools.TransformDataSlice(commentsData),
			},
		)
	}()
	// 字段基础校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 先拿顾问的info
	if code, infoData = service.GetTableItemsById(service.ADVISORTABLE, data.Id, []string{"*"}); code == errmsg.SUCCESS {
		// 在拿用户的服务
		code, serviceData = service.GetAdvisorService(data.Id)
	}
	delete(infoData, "password")
	infoData["coin"] = tools.ConvertCoinI2F(infoData["coin"].(int64))
	// week3 更详细的信息 TODO
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
		ginBindError(ctx, err, data)
		return
	}
	data.AdvisorId = ctx.GetInt64("id")
	// return func
	defer func() {
		commonControllerLog(&code, &msg, data, data)
		commonReturn(ctx, code, msg, data)
	}()
	// 输入校验后执行业务逻辑
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 更新
	code = service.UpdateTableItemById(service.ADVISORTABLE, data.AdvisorId,
		map[string]interface{}{"status": data.Status})
	return
}
