package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"service/model"
	"service/service"
	"service/utils/errmsg"
	"service/utils/logger"
	"service/utils/tools"
	"service/utils/validator"
)

func UpdateUserInfoController(ctx *gin.Context) {
	var data model.UserInfo
	var code int
	var msg string
	var res map[string]interface{}
	//mapData := make(map[string]interface{}, 8)
	// 数据绑定,通过结构体绑定数据,如果数据输入不对在这里就会报错return
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ginBindError(ctx, err, "UpdateUserInfoController", &data)
		return
	}
	defer func() {
		if code != errmsg.SUCCESS {
			logger.Log.Warn(errmsg.GetErrMsg(code))
		}
		commonReturn(ctx, code, msg, data)
	}()
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
	code = service.Update(service.USERTABLE, res)
	return
}
func NewUser(ctx *gin.Context) {
	var table = service.USERTABLE
	var data model.Login
	// 数据绑定
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ginBindError(ctx, err, "NewUser", data)
		return
	}
	// 数据校验
	msg, code := validator.Validate(data)
	// 数据不合法
	if code != errmsg.SUCCESS {
		logger.Log.Warn("数据校验非法", zap.Error(err))
		commonReturn(ctx, code, msg, data)
		return
	}
	// 调用service层 检查手机号是否重复
	code = service.CheckPhoneExist(table, data.Phone)
	if code == errmsg.SUCCESS {
		// 用户密码加密存储
		data.Password = service.GetPwd(data.Password)
		code, data.Id = service.NewUser(table, &data, nil)
		logger.Log.Info("新增用户", zap.String("phone", data.Phone))
	}
	// success
	commonReturn(ctx, code, msg, data)
	return
}

// UpdateUserPwd 修改用户密码
func UpdateUserPwd(ctx *gin.Context) {
	UpdatePwdController(service.USERTABLE, ctx)
}

// UserLoginController 用户登录
func UserLoginController(ctx *gin.Context) {
	Login(service.USERTABLE, ctx)
}

func GetUserInfoController(ctx *gin.Context) {
	id := ctx.GetInt64("id")
	code, data := service.GetUserInfo(id)
	commonReturn(ctx, code, "", data)
	return
}
