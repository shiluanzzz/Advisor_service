package v1

import (
	"fmt"
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

func UpdateUserInfoController(ctx *gin.Context) {
	// 前端传什么后端就更新什么
	var data map[string]interface{}
	var code int
	var msg string
	// 数据绑定
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ginBindError(ctx, err, "UpdateUserInfoController", &data)
		return
	}
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		}
		commonReturn(ctx, code, msg, data)
	}()
	// 数据校验 将不同的字段绑定到不同的校验函数中，使用反射做校验
	// 不存在的字段在函数中做了检验
	validateFunc := map[string]interface{}{
		"name":   validator.Name,
		"phone":  validator.Phone,
		"birth":  validator.Birth,
		"gender": validator.Gender,
		"bio":    validator.Bio,
		"about":  validator.About,
		"coin":   validator.CoinFunc,
	}
	//var key string
	//var value interface{}
	for key, value := range data {
		// 判断是否传的都是字符类型 手机号码传数字会被识别为float不好处理
		if key == "phone" {
			if reflect.TypeOf(value).Kind() == reflect.TypeOf(1.0).Kind() {
				value = strconv.FormatFloat(value.(float64), 'f', 0, 64)
				data[key] = value
			}
		}
		msg, code = validator.CallFunc(validateFunc, key, value)
		if code != errmsg.SUCCESS {
			msg = "数据校验非法" + msg
			return
		}
	}
	data["id"] = ctx.GetInt64("id")
	// 检查手机号是否重复
	if data["phone"] != nil {
		code, value := service.GetTableItem(service.USERTABLE, data["id"].(int64), "phone")
		// 电话号码有变动
		if fmt.Sprintf("%s", value) != data["phone"].(string) {
			code = service.CheckPhoneExist(service.USERTABLE, data["phone"])
			if code != errmsg.SUCCESS {
				return
			}
		}
	}
	// 更新
	code = service.Update(service.USERTABLE, data)
	return
}
func NewUser(ctx *gin.Context) {
	var table = service.USERTABLE
	var data model.Login
	// 数据绑定
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ginBindError(ctx, err, "NewUser", data)
		return
	}
	// 数据校验
	msg, code := validator.Validate(data)
	// 数据不合法
	if code != errmsg.SUCCESS {
		logger.Log.Warn("数据校验非法", zap.Error(err))
		commonReturn(ctx, code, msg, data)
		return
	}
	// 调用service层 检查手机号是否重复
	code = service.CheckPhoneExist(table, data.Phone)
	if code == errmsg.SUCCESS {
		// 用户密码加密存储
		data.Password = service.GetPwd(data.Password)
		code, data.Id = service.NewUser(table, &data, nil)
		logger.Log.Info("新增用户", zap.String("phone", data.Phone))
	}
	// success
	commonReturn(ctx, code, msg, data)
	return
}

// UpdateUserPwd 修改用户密码
func UpdateUserPwd(ctx *gin.Context) {
	UpdatePwdController(service.USERTABLE, ctx)
}

// UserLoginController 用户登录
func UserLoginController(ctx *gin.Context) {
	Login(service.USERTABLE, ctx)
}

func GetUserInfoController(ctx *gin.Context) {
	id := ctx.GetInt64("id")
	code, data := service.GetUserInfo(id)
	commonReturn(ctx, code, "", data)
	return
}
