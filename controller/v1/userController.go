package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service/middleware"
	"service/model"
	"service/service"
	"service/utils/errmsg"
	"service/utils/validator"
)

func NewUserController(ctx *gin.Context) {
	var data model.User
	_ = ctx.ShouldBindJSON(&data)
	// 数据校验
	msg, code := validator.Validate(data)
	// 数据不合法
	if code != errmsg.SUCCESS {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		return
	}
	// 用户密码加密存储
	data.Password = service.GetPwd(data.Password)

	// 调用service层 检查手机号是否重复
	code = service.CheckPhoneExist(service.USERTABLE, data.Phone)
	if code == errmsg.SUCCESS {
		code = service.NewUser(&data)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code),
		"data": data,
	})
}
func UpdateUserInfoController(ctx *gin.Context) {
	var data model.User
	var code int
	_ = ctx.ShouldBindJSON(&data)
	// 跳过validate的一些校验,实际更新的时候也不会涉及这两个字段
	data.Password = "*********"
	data.Phone = ctx.GetString("phone")
	msg, code := validator.Validate(data)
	// 数据不符合要求
	if code != errmsg.SUCCESS {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		return
	}
	code = service.UpdateUser(&data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code),
		"data": data,
	})
}

func UpdateUserPwd(ctx *gin.Context) {
	phone := ctx.GetString("phone") // 从token里拿
	oldPwd := ctx.PostForm("oldPwd")
	newPwd := ctx.PostForm("newPwd")
	var code int
	// 输入基本检查
	if oldPwd == newPwd || newPwd == "" {
		code = errmsg.ERROR_INPUT
		ctx.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  errmsg.GetErrMsg(code),
		})
		return
	}
	// 检查旧密码是否正确
	code = service.CheckRolePwd(service.USERTABLE, phone, oldPwd)

	if code == errmsg.SUCCESS {
		// update
		code = service.ChangePWD(service.USERTABLE, phone, newPwd)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code),
	})
}
func UserLogin(ctx *gin.Context) {
	phone := ctx.Query("phone")
	Pwd := ctx.Query("password")
	var code int
	var token string
	if phone == "" || Pwd == "" {
		code = errmsg.ERROR_INPUT
	} else {
		// 检查用户密码是否正确 里面包含了用户不存在的情况
		code = service.CheckRolePwd(service.USERTABLE, phone, Pwd)
		if code == errmsg.SUCCESS {
			token, code = middleware.NewToken(phone)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   errmsg.GetErrMsg(code),
		"token": token,
	})
}

func GetUserInfo(ctx *gin.Context) {
	//这个username 是token鉴权成功后写入到请求中的
	phone := ctx.GetString("phone")
	code, data := service.GetUser(phone)
	ctx.JSON(http.StatusOK, gin.H{
		"code":        code,
		"msg":         errmsg.GetErrMsg(code),
		"token_phone": phone,
		"data":        data,
	})
}
