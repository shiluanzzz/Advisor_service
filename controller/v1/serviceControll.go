package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service/model"
	"service/service"
	"service/utils/errmsg"
	"service/utils/validator"
)

// NewService 新增一个顾客的服务类型
func NewService(ctx *gin.Context) {
	var data model.Service
	_ = ctx.ShouldBindJSON(&data)
	// token 鉴权的接口直接从token拿数据
	data.AdvisorPhone = ctx.GetString("phone")
	// 获取这个服务对应的ID
	data.ServiceId = service.GetServiceId(data.ServiceName)
	msg, code := validator.Validate(data)
	if code == errmsg.SUCCESS {
		//检查是否有重复的服务
		code = service.CheckService(data.ServiceId, data.AdvisorPhone)
		if code == errmsg.ERROR_SERVICE_NOT_EXIST {
			//不存在的话在去创建服务
			code = service.NewService(&data)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code) + " " + msg,
		"data": data,
	})
}

// ModifyServiceStatus 修改顾客的服务状态
func ModifyServiceStatus(ctx *gin.Context) {
	type serviceStatus struct {
		Phone  string `json:"phone"`
		ID     int    `json:"id" validate:"required,number"`
		Status int    `json:"status" validate:"number,min=0,max=1"`
	}
	var data serviceStatus
	_ = ctx.ShouldBindJSON(&data)
	// 数据校验
	msg, code := validator.Validate(data)
	data.Phone = ctx.GetString("phone")
	if code == errmsg.SUCCESS {
		// 检查顾客是否有这个服务
		code = service.CheckService(data.ID, data.Phone)
		if code == errmsg.ERROR_SERVICE_EXIST {
			code = service.ModifyServiceStatus(data.Phone, data.ID, data.Status)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code) + " " + msg,
		"data": data,
	})
}

// ModifyServicePrice 修改顾客服务的价格
func ModifyServicePrice(ctx *gin.Context) {
	type servicePrice struct {
		Phone string  `json:"phone"`
		ID    int     `json:"id" validate:"required,number"`
		Price float32 `json:"price" validate:"required,number,gte=1,lte=36"`
	}
	var data servicePrice
	_ = ctx.ShouldBindJSON(&data)
	data.Phone = ctx.GetString("phone")
	// 数据校验
	msg, code := validator.Validate(data)
	// 数据校验通过
	if code == errmsg.SUCCESS {
		// 检查顾客是否有这个服务
		code = service.CheckService(data.ID, data.Phone)
		if code == errmsg.ERROR_SERVICE_EXIST {
			// 修改价格
			code = service.ModifyServicePrice(data.Phone, data.ID, data.Price)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  errmsg.GetErrMsg(code) + " " + msg,
		"data": data,
	})
}
