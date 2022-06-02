package v1

import (
	"fmt"
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
	"strings"
	"unicode"
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
func ginBindError(ctx *gin.Context, err error, funcName string, data interface{}) {
	code := errmsg.ErrorGinBind
	logger.Log.Error("gin绑定json错误", zap.Error(err), zap.String("function", funcName))
	commonReturn(ctx, code, "", data)
	return
}

// Login 用户或者顾问登录
func Login(table string, ctx *gin.Context) {
	var data model.Login
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ginBindError(ctx, err, "Login", &data)
		return
	}

	// 数据校验
	msg, code := validator.Validate(data)
	// 数据不合法
	if code != errmsg.SUCCESS {
		logger.Log.Warn("数据校验非法", zap.String("msg", msg))
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

// UpdatePwdController  修改用户密码
func UpdatePwdController(table string, ctx *gin.Context) {
	// 拿数据
	var data model.ChangePwd
	var msg string
	var code int
	if err := ctx.ShouldBind(&data); err != nil {
		ginBindError(ctx, err, "UpdatePwdController", data)
		return
	}
	// 数据校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		commonReturn(ctx, code, msg, data)
		return
	}
	id := ctx.GetInt64("id")
	// 检查旧密码是否正确
	if code = service.CheckRolePwd(table, id, data.OldPassword); code == errmsg.SUCCESS {
		// 更新密码
		code = service.ChangePWD(table, id, data.NewPassword)
	}
	logger.Log.Info(fmt.Sprintf("%s 修改密码", table), zap.String("id", strconv.FormatInt(id, 10)))
	commonReturn(ctx, code, "", data)
}

// LowFirst 首字母小写 SomeThing->someThing
func LowFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// Case2CamelCase 蛇形转驼峰 some_thing -> someThing
func Case2CamelCase(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	str = strings.Replace(str, " ", "", -1)
	return LowFirst(str)
}

// TransformDataSlice 把数据转换为小驼峰返回
func TransformDataSlice(data []map[string]interface{}) []map[string]interface{} {
	var res []map[string]interface{}
	for _, each := range data {
		res = append(res, TransformData(each))
	}
	return res
}

// TransformData 数据的key转化为小驼峰返回
func TransformData(data map[string]interface{}) map[string]interface{} {
	t := map[string]interface{}{}
	for k, v := range data {
		t[Case2CamelCase(k)] = v
	}
	return t
}

// StructToMap 结构体转为Map[string]interface{},忽略nil指针
func StructToMap(in interface{}, tagName string) (map[string]interface{}, int) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, errmsg.ERROR
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			// 如果这个指向的是一个空指针就不用添加到map里去。
			if !v.Field(i).IsNil() {
				out[tagValue] = v.Field(i).Interface()
			}
		}
	}
	return out, errmsg.SUCCESS
}
