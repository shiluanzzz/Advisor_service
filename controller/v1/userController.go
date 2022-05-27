package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"service/middleware"
	"service/model"
	"service/service"
	"service/utils/errmsg"
	"service/utils/logger"
	"service/utils/validator"
	"strconv"
)

func commonReturn(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  errmsg.GetErrMsg(code) + " " + msg,
		"code": code,
		"data": data,
	})
	return
}
func GinBindError(ctx *gin.Context, err error, funcName string, data interface{}) {
	code := errmsg.ErrorGinBind
	logger.Log.Error("gin绑定json错误", zap.Error(err), zap.String("function", funcName))
	commonReturn(ctx, code, "", data)
	return

}
func NewUser(ctx *gin.Context) {
	var data model.UserLogin
	var code int

	// 数据绑定
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		GinBindError(ctx, err, "NewUser", data)
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
	// 用户密码加密存储
	data.Password = service.GetPwd(data.Password)

	// 调用service层 检查手机号是否重复
	code = service.CheckPhoneExist(service.USERTABLE, data.Phone)
	if code == errmsg.SUCCESS {
		code = service.NewUser(&data)
	}
	// success
	logger.Log.Info("新增用户", zap.String("phone", data.Phone))
	commonReturn(ctx, code, msg, &data)
	return
}

func UpdateUserInfoController(ctx *gin.Context) {
	// 前端传什么后端就更新什么
	var data map[string]interface{}
	var code int

	// 数据绑定
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		GinBindError(ctx, err, "NewUser", &data)
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
	for key, value := range data {
		// 判断是否传的都是字符类型 手机号码传数字会被识别为float不好处理
		if key == "phone" && reflect.TypeOf(value).Kind() != reflect.TypeOf("1").Kind() {
			commonReturn(ctx, errmsg.ERROR_PHONE_INPUT, "", &data)
			return
		}
		msg, code := validator.CallFunc(validateFunc, key, value)
		if code != errmsg.SUCCESS {
			logger.Log.Warn("数据校验非法", zap.Error(err))
			commonReturn(ctx, code, msg, &data)
			return
		}
	}
	data["id"] = ctx.GetInt64("id")
	// 检查手机号是否重复
	if data["phone"] != nil {
		code = service.CheckPhoneExist(service.USERTABLE, data["phone"])
		if code != errmsg.SUCCESS {
			commonReturn(ctx, code, "", &data)
			return
		}
	}
	// 更新
	code = service.UpdateUser(data)
	commonReturn(ctx, code, "", &data)
	return
}

// UpdateUserPwd 修改用户密码
func UpdateUserPwd(ctx *gin.Context) {
	// 拿数据
	var data model.ChangePwd
	err := ctx.ShouldBind(&data)
	// 数据校验
	msg, code := validator.Validate(data)
	id := ctx.GetInt64("id")
	// 数据不合法
	if code != errmsg.SUCCESS {
		logger.Log.Warn("数据校验非法", zap.Error(err))
		commonReturn(ctx, code, msg, data)
		return
	}

	// 检查旧密码是否正确
	code = service.CheckRolePwd(service.USERTABLE, id, data.OldPwd)

	if code == errmsg.SUCCESS {
		// update
		code = service.ChangePWD(service.USERTABLE, id, data.NewPwd)
	}
	logger.Log.Info("用户修改密码", zap.String("id", strconv.FormatInt(id, 10)))
	commonReturn(ctx, code, "", data)
}

// UserLogin 用户登录
func UserLogin(ctx *gin.Context) {
	var data model.UserLogin
	err := ctx.ShouldBindQuery(&data)
	// 数据绑定错误
	if err != nil {
		GinBindError(ctx, err, "UserLogin", &data)
		return
	}

	// 数据校验
	msg, code := validator.Validate(data)
	// 数据不合法
	if code != errmsg.SUCCESS {
		logger.Log.Warn("数据校验非法", zap.Error(err), zap.String("msg", msg))
		commonReturn(ctx, code, msg, data)
		return
	}
	// 获取用户的ID
	data.Id, code = service.GetUserId(data.Phone)
	// 检查用户密码是否正确
	if code == errmsg.SUCCESS {
		code = service.CheckRolePwd(service.USERTABLE, data.Id, data.Password)
		// 生成Token
		if code == errmsg.SUCCESS {
			data.Token, code = middleware.NewToken(data.Id)
		}
	}
	commonReturn(ctx, code, "", data)
}

func GetUserInfo(ctx *gin.Context) {
	//这个username 是token鉴权成功后写入到请求中的
	id := ctx.GetInt64("id")
	code, data := service.GetUser(id)
	commonReturn(ctx, code, "", data)
}
