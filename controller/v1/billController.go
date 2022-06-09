package v1

import (
	"github.com/gin-gonic/gin"
	"service-backend/model"
	"service-backend/service"
	"service-backend/utils/logger"
)

func GetBillController(ctx *gin.Context) {
	var code int
	var response []*model.Bill
	defer func() {
		logger.CommonControllerLog(&code, nil, ctx.GetInt64("id"), response)
		commonReturn(ctx, code, "", response)
	}()
	code, response = service.GetBill(ctx.GetInt64("id"), ctx.GetString("role"))
	return
}
