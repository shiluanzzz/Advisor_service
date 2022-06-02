package v1

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"service/model"
	"service/service"
	"service/utils/errmsg"
	"service/utils/logger"
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
		code, data.Id = service.NewAdvisorAndOrder(&data)
	}
	// success
	return
}

// TODO 修改
func UpdateAdvisorController(ctx *gin.Context) {
	var data map[string]interface{}
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
		commonReturn(ctx, code, msg, TransformData(data))
	}()

	// 数据校验 将不同的字段绑定到不同的校验函数中，使用反射做校验
	// 不存在的字段在函数中做了检验
	validateFunc := map[string]interface{}{
		"name":           validator.Name,
		"phone":          validator.Phone,
		"workExperience": validator.WorkExperience,
		"bio":            validator.Bio,
		"about":          validator.About,
	}
	for key, value := range data {
		// 判断是否传的都是字符类型 手机号码传数字会被识别为float不好处理
		// json中的字符被识别为float
		if key == "phone" {
			if reflect.TypeOf(value).Kind() == reflect.TypeOf(1.0).Kind() {
				value = strconv.FormatFloat(value.(float64), 'f', 0, 64)
				data[key] = value
			}
		}
		if key == "workExperience" {
			if reflect.TypeOf(value).Kind() != reflect.TypeOf(1.0).Kind() {
				code = errmsg.ErrorInput
				msg = "请检查workExperience字段的输入"
				return
			}
			value = int(value.(float64))
			data[key] = value
		}
		msg, code = validator.CallFunc(validateFunc, key, value)
		if code != errmsg.SUCCESS {
			return
		}
	}
	data["id"] = ctx.GetInt64("id")
	// 检查手机号是否重复
	if data["phone"] != nil {
		code = service.CheckPhoneExist(service.USERTABLE, data["phone"])
		if code != errmsg.SUCCESS {
			return
		}
	}
	// 更新
	data["work_experience"] = data["workExperience"]
	delete(data, "workExperience")
	code = service.Update(service.ADVISORTABLE, data)
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
	var serviceData []map[string]interface{}
	var infoData map[string]interface{}
	var code int
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ginBindError(ctx, err, "GetAdvisorInfo", data)
	}
	// 字段基础校验
	if msg, code := validator.Validate(data); code != errmsg.SUCCESS {
		commonReturn(ctx, code, msg, data)
		return
	}
	// 先拿用户的info
	if code, infoData = service.GetManyTableItemsById(service.ADVISORTABLE, data.Id, []string{"*"}); code == errmsg.SUCCESS {
		// 在拿用户的服务
		//code, serviceData = service.GetAdvisorService(data.Id)
		code, serviceData = service.GetManyTableItemsByWhere(service.ADVISORTABLE,
			map[string]interface{}{"advisor_id": data.Id, "status": 1},
			[]string{"*"},
		)
	}

	commonReturn(ctx, code, "",
		map[string]interface{}{
			"info":    TransformData(infoData),
			"service": TransformDataSlice(serviceData),
		},
	)
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
