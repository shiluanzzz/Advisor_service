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

	// 调用service层
	code = service.CheckUserPhone(data.Phone)
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
	code = service.UpdateUser(&data)
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code),
		"data": data,
	})
}

func UpdateUserPwd(ctx *gin.Context) {
	username := ctx.PostForm("username")
	oldPwd := ctx.PostForm("oldPwd")
	newPwd := ctx.PostForm("newPwd")
	// 检查密码是否正确
	code := service.CheckRolePwd(service.USERTABLE, username, oldPwd)

	if code == errmsg.SUCCESS {
		// update
		code = service.ChangeUserPWD(username, newPwd)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code),
	})
}
func UserLogin(ctx *gin.Context) {
	username := ctx.Query("username")
	Pwd := ctx.Query("password")
	// 检查用户密码是否正确 里面包含了用户不存在的情况
	code := service.CheckRolePwd(service.USERTABLE, username, Pwd)
	var token string
	if code == errmsg.SUCCESS {
		token, code = middleware.NewToken(username)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   errmsg.GetErrMsg(code),
		"token": token,
	})
}

func GetUserInfo(ctx *gin.Context) {
	//这个username 是token鉴权成功后写入到请求中的
	username := ctx.GetString("username")
	code, data := service.GetUser(username)
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code),
		"data": data,
	})
}
