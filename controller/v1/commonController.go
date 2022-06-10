package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"service-backend/middleware"
	"service-backend/model"
	"service-backend/service"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
	"service-backend/utils/tools"
	"service-backend/utils/validator"
)

// 统一的gin 数据返回格式
func commonReturn(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  errmsg.GetErrMsg(code) + " " + msg,
		"code": code,
		"data": data,
	})
	return
}

// defer
func commonControllerDefer(ctx *gin.Context, code *int, msg *string, request interface{}, data interface{}) {
	if msg == nil {
		ss := ""
		msg = &ss
	}
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		logger.Log.Warn("接口应当传递指针", zap.String("function", tools.WhoCallMe()))
	}
	if err := recover(); err != nil {
		logger.Log.Error("controller层Panic", zap.String("err", fmt.Sprintf("%v", err)), zap.String("function", tools.WhoCallMe()))
	}
	logger.CommonControllerLog(code, msg, request, data, "function", tools.WhoAmI())
	commonReturn(ctx, *code, *msg, data)
}

// ginBindError gin绑定数据的error 返回
func ginBindError(ctx *gin.Context, err error, data interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Log.Error("gin框架Panic", zap.String("err", fmt.Sprintf("%v", err)), zap.String("function", tools.WhoCallMe()))
		}
	}()
	code := errmsg.ErrorGinBind
	logger.Log.Error("gin绑定json错误", zap.Error(err), zap.String("function", tools.WhoCallMe()))
	commonReturn(ctx, code, "", data)
	return
}

// LoginController 用户或者顾问登录
func LoginController(table string, ctx *gin.Context) {
	var data model.Login
	var code int
	var msg string
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ginBindError(ctx, err, &data)
		return
	}
	defer commonControllerDefer(ctx, &code, &msg, &data, &data)
	// 数据校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 获取用户的ID
	var id interface{}
	if code, id = service.GetTableItemByWhere(table, map[string]interface{}{"phone": data.Phone}, "id"); code != errmsg.SUCCESS {
		return
	}
	data.Id = id.(int64)
	// 检查用户密码是否正确
	if code = service.CheckRolePwd(table, data.Id, data.Password); code != errmsg.SUCCESS {
		return
	}
	// 生成Token
	data.Token, code = middleware.NewToken(data.Id, table)
	return
}

// UpdatePwdController  修改用户密码
func UpdatePwdController(table string, ctx *gin.Context) {
	// 拿数据
	var data model.ChangePwd
	var msg string
	var code int
	if err := ctx.ShouldBind(&data); err != nil {
		ginBindError(ctx, err, data)
		return
	}
	defer commonControllerDefer(ctx, &code, &msg, data, data)
	// 数据校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	data.Id = ctx.GetInt64("id")
	// 检查旧密码是否正确
	if code = service.CheckRolePwd(table, data.Id, data.OldPassword); code != errmsg.SUCCESS {
		return
	}
	// 更新密码
	if code = service.ChangePWD(table, data.Id, data.NewPassword); code != errmsg.SUCCESS {
		return
	}
	return
}
