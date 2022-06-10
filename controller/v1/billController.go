package v1

import (
	"github.com/gin-gonic/gin"
	"service-backend/model"
	"service-backend/service"
	"service-backend/utils/tools"
	"time"
)

func GetBillController(ctx *gin.Context) {
	var code int
	var response []*model.Bill
	defer commonControllerDefer(ctx, &code, nil, ctx.GetInt64("id"), &response)
	code, response = service.GetBill(ctx.GetInt64("id"), ctx.GetString("role"))

	// 完善一些展示信息
	for _, v := range response {
		v.BillType = v.Type.Name()
		v.ShowTime = time.Unix(v.Time, 0).Format("Jan 02,2006 15:04:05")
		v.ShowAmount = tools.ConvertCoinI2F(v.Amount)
	}
	return
}
