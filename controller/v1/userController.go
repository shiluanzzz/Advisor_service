package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service/model"
	"service/service"
	"service/utils/errmsg"
)

func NewUserController(ctx *gin.Context) {
	var data model.User
	_ = ctx.ShouldBindJSON(&data)
	// 数据校验 govalidator TODO

	// 用户密码加密存储
	data.Password = service.GetPwd(data.Password)

	// 调用service层
	code := service.CheckUserName(data.Name)
	if code == errmsg.SUCCESS {
		code = service.NewUser(&data)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  errmsg.GetErrMsg(code),
		"data": data,
	})
}
func UpdateUserInfoController(ctx *gin.Context) {
	var data model.User
	_ = ctx.ShouldBindJSON(&data)
	// 校验用户密码
	code := service.CheckRolePwd(service.USERTABLE, data.Name, data.Password)
	// 用户密码验证通过 执行下一步
	if code == errmsg.SUCCESS {
		code = service.UpdateUser(&data)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  errmsg.GetErrMsg(code),
		"data": data,
	})
}

func UpdateUserPwd(ctx *gin.Context) {
	username := ctx.PostForm("username")
	oldPwd := ctx.PostForm("oldPwd")
	newPwd := ctx.PostForm("newPwd")
	code := service.CheckRolePwd(service.USERTABLE, username, oldPwd)
	if code == errmsg.SUCCESS {
		code = service.ChangeUserPWD(username, newPwd)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": errmsg.GetErrMsg(code),
	})
}
