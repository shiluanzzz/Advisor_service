package v1

import (
	"github.com/gin-gonic/gin"
	"service/model"
	"service/service"
	"service/utils/errmsg"
	"service/utils/validator"
)

// ModifyServiceStatus 修改顾客的服务状态
func ModifyServiceStatus(ctx *gin.Context) {
	type serviceStatus struct {
		AdvisorId     int64 `json:"advisorId"`
		ServiceNameId int   `form:"serviceNameId" json:"serviceNameId" validate:"required,number,lte=4"`
		Status        int   `form:"status" json:"status" validate:"number,min=0,max=1"`
	}
	var data serviceStatus
	var msg string
	var code int
	if err := ctx.ShouldBind(&data); err != nil {
		ginBindError(ctx, err, "ModifyServiceStatus", data)
		return
	}
	// 数据校验
	if msg, code = validator.Validate(data); code == errmsg.SUCCESS {
		// 修改状态
		data.AdvisorId = ctx.GetInt64("id")
		code = service.ModifyServiceStatus(data.AdvisorId, data.ServiceNameId, data.Status)
	}
	commonReturn(ctx, code, msg, data)
	return
}

// ModifyServicePrice 修改顾客服务的价格
func ModifyServicePrice(ctx *gin.Context) {

	var data model.ServicePrice
	var msg string
	var code int
	if err := ctx.ShouldBind(&data); err != nil {
		ginBindError(ctx, err, "ModifyServiceStatus", data)
		return
	}
	data.AdvisorId = ctx.GetInt64("id")
	// 数据校验
	if msg, code = validator.Validate(data); code == errmsg.SUCCESS {
		// 修改价格
		code = service.ModifyServicePrice(data.AdvisorId, data.ServiceID, data.Price)
	}
	commonReturn(ctx, code, msg, data)
	return
}
