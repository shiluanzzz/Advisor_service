package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"service/middleware"
	"service/model"
	"service/service"
	"service/utils/errmsg"
	"service/utils/logger"
	"service/utils/validator"
	"strconv"
	"strings"
	"unicode"
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
func NewUserOrAdvisor(table string, ctx *gin.Context) {
	var data model.Login
	// 数据绑定
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		GinBindError(ctx, err, "NewUserOrAdvisor", data)
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
		code, data.Id = service.NewUser(table, &data)
		logger.Log.Info("新增用户", zap.String("phone", data.Phone))
		// 如果是顾问 在为顾问创建一次服务
		if table == service.ADVISORTABLE {
			// TODO
			code = service.NewService(data.Id)
		}
	}
	// success
	commonReturn(ctx, code, msg, &data)
	return
}

// Login 用户或者顾问登录
func Login(table string, ctx *gin.Context) {
	var data model.Login
	err := ctx.ShouldBindQuery(&data)
	fmt.Println(data)
	// 数据绑定错误
	if err != nil {
		GinBindError(ctx, err, "Login", &data)
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
	data.Id, code = service.GetId(table, data.Phone)
	// 检查用户密码是否正确
	if code == errmsg.SUCCESS {
		code = service.CheckRolePwd(table, data.Id, data.Password)
		// 生成Token
		if code == errmsg.SUCCESS {
			data.Token, code = middleware.NewToken(data.Id, table)
		}
	}
	commonReturn(ctx, code, "", data)
}

// UpdateUserPwd 修改用户密码
func UpdatePwdControl(table string, ctx *gin.Context) {
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
	code = service.CheckRolePwd(table, id, data.OldPassword)

	if code == errmsg.SUCCESS {
		// update
		code = service.ChangePWD(table, id, data.NewPassword)
	}
	logger.Log.Info(fmt.Sprintf("%s 修改密码", table), zap.String("id", strconv.FormatInt(id, 10)))
	commonReturn(ctx, code, "", data)
}

func Case2CamelCase(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	str = strings.Replace(str, " ", "", -1)
	return LowFirst(str)
}

// 首字母小写
func LowFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// 把数据转换为小驼峰返回
func TransformDataSlice(data []map[string]interface{}) []map[string]interface{} {
	var res []map[string]interface{}
	for _, each := range data {
		res = append(res, TransformData(each))
	}
	return res
}
func TransformData(data map[string]interface{}) map[string]interface{} {
	t := map[string]interface{}{}
	for k, v := range data {
		t[Case2CamelCase(k)] = v
	}
	return t
}
