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
)

//func NewUser(ctx *gin.Context) {
//	NewUser(service.USERTABLE, ctx)
//}

func UpdateUserInfoController(ctx *gin.Context) {
	// 前端传什么后端就更新什么
	var data map[string]interface{}
	var code int

	// 数据绑定
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ginBindError(ctx, err, "UpdateUserInfoController", &data)
		return
	}
	// 数据校验 将不同的字段绑定到不同的校验函数中，使用反射做校验
	// 不存在的字段在函数中做了检验
	validateFunc := map[string]interface{}{
		"name":   validator.Name,
		"phone":  validator.Phone,
		"birth":  validator.Birth,
		"gender": validator.Gender,
		"bio":    validator.Bio,
		"about":  validator.About,
	}
	//var key string
	//var value interface{}
	for key, value := range data {
		// 判断是否传的都是字符类型 手机号码传数字会被识别为float不好处理
		if key == "phone" && reflect.TypeOf(value).Kind() != reflect.TypeOf("1").Kind() {
			commonReturn(ctx, errmsg.ErrorPhoneInput, "", data)
			return
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
	code = service.Update(service.USERTABLE, data)
	commonReturn(ctx, code, "", data)
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
	// TODO 顾问创建和顾问的服务创建用事务
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

// UserLogin 用户登录
func UserLogin(ctx *gin.Context) {
	Login(service.USERTABLE, ctx)
}

func GetUserInfo(ctx *gin.Context) {
	//这个username 是token鉴权成功后写入到请求中的
	id := ctx.GetInt64("id")
	code, data := service.GetUserInfo(id)
	commonReturn(ctx, code, "", data)
}
