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

}

func UpdateUserPwd(ctx *gin.Context) {

}
