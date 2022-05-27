package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"reflect"
	"service/model"
	"service/service"
	"service/utils/errmsg"
	"service/utils/logger"
	"service/utils/validator"
	"strconv"
)

func NewAdvisorController(ctx *gin.Context) {
	NewUserOrAdvisor(service.ADVISORTABLE, ctx)
}
func UpdateAdvisorController(ctx *gin.Context) {

	var data map[string]interface{}
	//data := map[string]interface{}{}
	var code int
	// 数据绑定
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		GinBindError(ctx, err, "UpdateAdvisorController", data)
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
	if value := data["workExperience"]; value != nil {
		data["work_experience"] = int(value.(float64))
		value = int(value.(float64))
		delete(data, "workExperience")
	}
	for key, value := range data {
		// 判断是否传的都是字符类型 手机号码传数字会被识别为float不好处理
		// json中的字符被识别为float
		if key == "phone" && reflect.TypeOf(value).Kind() != reflect.TypeOf("1").Kind() {
			value = strconv.FormatFloat(value.(float64), 'f', 0, 64)
		}

		msg, code := validator.CallFunc(validateFunc, key, value)
		if code != errmsg.SUCCESS {
			logger.Log.Warn("数据校验非法", zap.Error(err))
			commonReturn(ctx, code, msg, data)
			return
		}
	}
	data["id"] = ctx.GetInt64("id")
	// 检查手机号是否重复
	if data["phone"] != nil {
		code = service.CheckPhoneExist(service.USERTABLE, data["phone"])
		if code != errmsg.SUCCESS {
			commonReturn(ctx, code, "", data)
			return
		}
	}
	// 更新
	code = service.Update(service.ADVISORTABLE, data)
	commonReturn(ctx, code, "", data)
}

func UpdateAdvisorPwd(ctx *gin.Context) {
	UpdatePwdControl(service.ADVISORTABLE, ctx)
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

// AdvisorLogin 顾问登录
func AdvisorLogin(ctx *gin.Context) {
	Login(service.ADVISORTABLE, ctx)
}

// GetAdvisorInfo 获取顾问的详细信息
func GetAdvisorInfo(ctx *gin.Context) {
	type Num struct {
		Id int64 `form:"id" validate:"required,min=0"`
	}
	var data Num
	err := ctx.ShouldBindQuery(&data)
	if err != nil {
		GinBindError(ctx, err, "GetAdvisorInfo", data)
	}
	code, res := service.GetAdvisorInfo(data.Id)

	var serviceData []map[string]interface{}
	if code == errmsg.SUCCESS {
		code, serviceData = service.GetAdvisorService(data.Id)
	}
	commonReturn(ctx, code, "",
		map[string]interface{}{
			"info":    TransformDataSlice(res),
			"service": TransformDataSlice(serviceData),
		},
	)
}

// ModifyAdvisorStatus 顾问修改自己的服务状态
func ModifyAdvisorStatus(ctx *gin.Context) {
	id := ctx.GetInt64("id")
	var newStatus model.ServiceState
	err := ctx.ShouldBind(&newStatus)
	returnData := map[string]interface{}{
		"id":     id,
		"status": newStatus.Status,
	}
	if err != nil {
		GinBindError(ctx, err, "ModifyAdvisorStatus", returnData)
		return
	}
	msg, code := validator.Validate(newStatus)
	if code == errmsg.SUCCESS {
		code = service.ModifyAdvisorStatus(id, newStatus.Status)
	}
	commonReturn(ctx, code, msg, returnData)
	return
}
