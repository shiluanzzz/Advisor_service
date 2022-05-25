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

func NewAdvisorController(ctx *gin.Context) {
	var data model.Advisor
	_ = ctx.ShouldBindJSON(&data)
	// 数据校验
	msg, code := validator.Validate(data)
	if code != errmsg.SUCCESS {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		return
	}
	// 加密
	data.Password = service.GetPwd(data.Password)
	// 检查重复
	// question: 同一个手机号是否可同时注册为顾问和顾客
	code = service.CheckPhoneExist(service.ADVISORTABLE, data.Phone)
	if code == errmsg.SUCCESS {
		code = service.NewAdvisor(&data)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code),
		"data": data,
	})
}
func UpdateAdvisorController(ctx *gin.Context) {
	var data model.Advisor
	var code int
	_ = ctx.ShouldBindJSON(&data)
	// 跳过校验
	data.Password = "*********"
	data.Phone = ctx.GetString("phone")
	msg, code := validator.Validate(data)
	if code != errmsg.SUCCESS {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
		return
	}
	code = service.UpdateAdvisor(&data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code),
		"data": data,
	})
}
func UpdateAdvisorPwd(ctx *gin.Context) {
	phone := ctx.GetString("phone")
	oldPwd := ctx.PostForm("oldPwd")
	newPwd := ctx.PostForm("newPwd")
	var code int
	if oldPwd == newPwd || newPwd == "" {
		code = errmsg.ERROR_INPUT
		ctx.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  errmsg.GetErrMsg(code),
		})
		return
	}
	// 检查新旧密码是否匹配
	code = service.CheckRolePwd(service.ADVISORTABLE, phone, oldPwd)
	if code == errmsg.SUCCESS {
		// update
		code = service.ChangePWD(service.ADVISORTABLE, phone, newPwd)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code),
	})
}

// 获取顾问的列表
func GetAdvisorList(ctx *gin.Context) {

}

// 顾问登录
func AdvisorLogin(ctx *gin.Context) {
	phone := ctx.Query("phone")
	pwd := ctx.Query("password")
	var code int
	var token string
	if phone == "" || pwd == "" {
		code = errmsg.ERROR_INPUT
	} else {
		code = service.CheckRolePwd(service.ADVISORTABLE, phone, pwd)
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

// 获取顾问的详细信息
func GetAdvisorInfo(ctx *gin.Context) {
	phone := ctx.Query("phone")
	code, data := service.GetAdvisorInfo(phone)
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code),
		"data": data,
	})
}
