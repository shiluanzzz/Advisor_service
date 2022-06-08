package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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

// ginBindError gin绑定数据的error 返回
func ginBindError(ctx *gin.Context, err error, data interface{}) {
	code := errmsg.ErrorGinBind
	logger.Log.Error("gin绑定json错误", zap.Error(err), zap.String("function", tools.WhoCallMe()))
	commonReturn(ctx, code, "", data)
	return
}

// Login 用户或者顾问登录
func Login(table string, ctx *gin.Context) {
	var data model.Login
	var code int
	var msg string
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ginBindError(ctx, err, &data)
		return
	}
	defer func() {
		logger.CommonControllerLog(&code, &msg, data, data)
		commonReturn(ctx, code, msg, data)
	}()
	// 数据校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 获取用户的ID
	//data.Id, code = service.GetId(table, data.Phone)
	code, id := service.GetTableItemByWhere(
		table, map[string]interface{}{"phone": data.Phone}, "id")
	if code != errmsg.SUCCESS {
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
	defer func() {
		logger.CommonControllerLog(&code, &msg, data, data)
		commonReturn(ctx, code, msg, data)
	}()
	// 数据校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	id := ctx.GetInt64("id")
	data.Id = id
	// 检查旧密码是否正确
	if code = service.CheckRolePwd(table, id, data.OldPassword); code != errmsg.SUCCESS {
		return
	}
	// 更新密码
	if code = service.ChangePWD(table, id, data.NewPassword); code != errmsg.SUCCESS {
		return
	}
	return
}
