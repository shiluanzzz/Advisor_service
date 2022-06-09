package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"service-backend/model"
	"service-backend/service"
	"service-backend/utils/errmsg"
	"service-backend/utils/tools"
	"service-backend/utils/validator"
)

func UpdateUserInfoController(ctx *gin.Context) {
	var data model.UserInfo
	var code int
	var msg string
	var res map[string]interface{}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ginBindError(ctx, err, &data)
		return
	}
	defer commonControllerDefer(ctx, &code, &msg, &data, &data)

	// 将结构体中非nil的字段提取到map中
	if res, code = tools.StructToMap(data, "structs"); code != errmsg.SUCCESS {
		return
	}
	validateFunc := map[string]interface{}{
		"name":   validator.Name,
		"phone":  validator.Phone,
		"birth":  validator.Birth,
		"gender": validator.Gender,
		"bio":    validator.Bio,
		"about":  validator.About,
		"coin":   validator.CoinFunc,
	}
	for key, value := range res {
		// 对更新的字段逐个做校验
		if msg, code = validator.CallFunc(validateFunc, key, value); code != errmsg.SUCCESS {
			return
		}
	}
	res["id"] = ctx.GetInt64("id")
	// 检查手机号是否重复
	if res["phone"] != nil {
		code, value := service.GetTableItem(service.USERTABLE, res["id"].(int64), "phone")
		// 电话号码有变动
		if fmt.Sprintf("%s", value) != *(res["phone"].(*string)) {
			code = service.CheckPhoneExist(service.USERTABLE, res["phone"])
			if code != errmsg.SUCCESS {
				return
			}
		}
	}
	// 用户金币乘base存储
	if res["coin"] != nil {
		res["coin"] = tools.ConvertCoinF2I(*(res["coin"].(*float32)))
	}
	// 更新
	code = service.UpdateTableItemById(service.USERTABLE, ctx.GetInt64("id"), res)
	return
}
func NewUserController(ctx *gin.Context) {
	var data model.Login
	var code int
	var msg string
	// 数据绑定
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ginBindError(ctx, err, data)
		return
	}
	defer commonControllerDefer(ctx, &code, &msg, data, &data)

	// 数据校验
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 调用service层 检查手机号是否重复
	if code = service.CheckPhoneExist(service.USERTABLE, data.Phone); code != errmsg.SUCCESS {
		return
	}
	// 用户密码加密存储
	data.Password = service.GetEncryptPwd(data.Password)
	// 入库
	insertMap := tools.Structs2SQLTable(data)
	code, data.Id = service.InsertTableItem(
		service.USERTABLE,
		[]map[string]interface{}{insertMap},
	)
	return
}

// UpdateUserPwd 修改用户密码
func UpdateUserPwd(ctx *gin.Context) {
	UpdatePwdController(service.USERTABLE, ctx)
}

// UserLoginController 用户登录
func UserLoginController(ctx *gin.Context) {
	LoginController(service.USERTABLE, ctx)
}

func GetUserInfoController(ctx *gin.Context) {
	id := ctx.GetInt64("id")
	var response *model.User
	var code int

	defer commonControllerDefer(ctx, &code, nil, &id, &response)

	code, response = service.GetUser(id)
	return
}
