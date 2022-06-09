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
	"time"
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
	defer commonControllerDefer(ctx, &code, &msg, data, data)
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
	var data model.AdvisorInfo
	var res map[string]interface{}
	var code int
	var msg string
	// 数据绑定
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ginBindError(ctx, err, data)
		return
	}
	defer commonControllerDefer(ctx, &code, &msg, data, data)
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
	var data []*model.Advisor
	defer commonControllerDefer(ctx, &code, nil, map[string]interface{}{"page": page}, &data)
	if err != nil || page < 1 {
		code = errmsg.ErrorInput
		commonReturn(ctx, code, "", map[string]int{"page": page})
		return
	}
	code, data = service.GetAdvisorList(page)
	// 删除敏感信息
	for _, v := range data {
		v.Coin = 0
	}
	//commonReturn(ctx, code, "", data)
}

// GetAdvisorInfoController 获取顾问的详细信息
func GetAdvisorInfoController(ctx *gin.Context) {

	var data model.TableID
	var serviceData []*model.Service
	var comment []*model.OrderComment
	var info *model.Advisor
	var code int
	var msg string
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ginBindError(ctx, err, data)
	}
	defer commonControllerDefer(ctx, &code, &msg, data, map[string]interface{}{
		"info":     &info,
		"service":  &serviceData,
		"comments": &comment,
	})

	// 字段基础校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 顾问的信息
	if code, info = service.GetAdvisor(data.Id); code != errmsg.SUCCESS {
		return
	}
	// 相关订单的评论
	if code, comment = service.GetAdvisorCommentData(data.Id); code != errmsg.SUCCESS {
		return
	}
	// 开放的服务
	if code, serviceData = service.GetAdvisorServiceData(data.Id); code != errmsg.SUCCESS {
		return
	}

	// 完善相关展示字段
	// 扩充数据
	for _, v := range comment {
		var userNameUint8 interface{}
		if code, userNameUint8 = service.GetTableItem(service.USERTABLE, v.UserId, "name"); code != errmsg.SUCCESS {
			return
		}
		v.UserName = fmt.Sprintf("%s", userNameUint8)
		v.CreateShowTime = time.Unix(v.OrderCreateTime, 0).Format("Jan 02,2006 15:04:05")
		v.CommentShowTime = time.Unix(v.CommentTime, 0).Format("Jan 02,2006 15:04:05")
	}
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
	defer commonControllerDefer(ctx, &code, &msg, data, data)
	// 输入校验后执行业务逻辑
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 更新
	code = service.UpdateTableItemById(service.ADVISORTABLE, data.AdvisorId,
		map[string]interface{}{"status": data.Status})
	return
}
