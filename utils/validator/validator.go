package validator

import (
	"fmt"
	"github.com/go-playground/locales/zh_Hans_CN"
	uTran "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"service/utils/errmsg"
)

func Validate(data interface{}) (string, int) {
	validate := validator.New()
	// 错误语言翻译成中文
	uni := uTran.New(zh_Hans_CN.New())
	trans, _ := uni.GetTranslator("zh_Hans_Cn")
	err := zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println("err:", err)
	}
	// 验证
	err = validate.Struct(data)
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			return v.Translate(trans), errmsg.ERROR
		}
	}
	return "", errmsg.SUCCESS
}
