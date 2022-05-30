package v1

import (
	"github.com/gin-gonic/gin"
	"service/service"
	"service/utils/errmsg"
	"service/utils/validator"
)

// ModifyServiceStatus 修改顾客的服务状态
func ModifyServiceStatus(ctx *gin.Context) {
	type serviceStatus struct {
		AdvisorId int64 `json:"advisorId"`
		ServiceID int   `form:"serviceId" json:"serviceId" validate:"required,number,lte=4"`
		Status    int   `form:"status" validate:"number,min=0,max=1"`
	}
	var data serviceStatus
	err := ctx.ShouldBind(&data)
	if err != nil {
		GinBindError(ctx, err, "ModifyServiceStatus", data)
		return
	}
	// 数据校验
	msg, code := validator.Validate(data)
	data.AdvisorId = ctx.GetInt64("id")
	if code == errmsg.SUCCESS {
		code = service.ModifyServiceStatus(data.AdvisorId, data.ServiceID, data.Status)
	}
	commonReturn(ctx, code, msg, data)
	return
}

// ModifyServicePrice 修改顾客服务的价格
func ModifyServicePrice(ctx *gin.Context) {
	type servicePrice struct {
		AdvisorId int64   `json:"advisorId"`
		ServiceID int     `form:"serviceId" json:"serviceId" validate:"required,number,lte=4"`
		Price     float32 `form:"price" validate:"required,number,gte=1,lte=36"`
	}
	var data servicePrice
	err := ctx.ShouldBind(&data)
	data.AdvisorId = ctx.GetInt64("id")
	if err != nil {
		GinBindError(ctx, err, "ModifyServiceStatus", data)
		return
	}
	// 数据校验
	msg, code := validator.Validate(data)
	if code == errmsg.SUCCESS {
		// 修改价格
		code = service.ModifyServicePrice(data.AdvisorId, data.ServiceID, data.Price)
	}
	commonReturn(ctx, code, msg, data)
	return
}
