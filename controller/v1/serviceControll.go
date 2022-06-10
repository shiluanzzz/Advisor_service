package v1

import (
	"github.com/gin-gonic/gin"
	"service-backend/model"
	"service-backend/service"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
	"service-backend/utils/tools"
	"service-backend/utils/validator"
)

// ModifyServiceStatus 修改顾客的服务状态
func ModifyServiceStatus(ctx *gin.Context) {

	var data model.ServiceState
	var msg string
	var code int
	if err := ctx.ShouldBind(&data); err != nil {
		ginBindError(ctx, err, data)
		return
	}
	defer func() {
		logger.CommonControllerLog(&code, &msg, data, data)
		commonReturn(ctx, code, msg, data)
	}()

	// 数据校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 修改状态
	data.AdvisorId = ctx.GetInt64("id")
	where := map[string]interface{}{
		"service_name_id": data.ServiceNameId,
		"advisor_id":      data.AdvisorId,
	}
	updates := map[string]interface{}{
		"status": data.Status,
	}
	code = service.UpDateTableItemByWhere(service.SERVICETABLE, where, updates)
	return
}

// ModifyServicePrice 修改顾客服务的价格
func ModifyServicePrice(ctx *gin.Context) {

	var data model.ServicePrice
	var msg string
	var code int
	if err := ctx.ShouldBind(&data); err != nil {
		ginBindError(ctx, err, data)
		return
	}
	defer commonControllerDefer(ctx, &code, &msg, &data, &data)

	data.AdvisorId = ctx.GetInt64("id")
	// 数据校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 修改价格
	where := map[string]interface{}{
		"service_name_id": data.ServiceNameId,
		"advisor_id":      data.AdvisorId,
	}
	updates := map[string]interface{}{
		"price": tools.ConvertCoinF2I(data.Price),
	}
	code = service.UpDateTableItemByWhere(service.SERVICETABLE, where, updates)
	return
}
