package validator

import (
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"service/utils/errmsg"
	"service/utils/logger"
)

func CallFunc(m map[string]interface{}, name string, params ...interface{}) (string, int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			logger.Log.Error("反射校验字段panic", zap.String("errorMsg", fmt.Sprintf("%v", err)))
		}
	}()
	if m[name] == nil {
		return fmt.Sprintf("不存在字段%s校验函数", name), errmsg.ErrorInput
	}
	f := reflect.ValueOf(m[name])
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result := f.Call(in)
	return result[0].String(), int(result[1].Int())
}

// Name 校验函数模板
//func XXXX(value interface{}) (errMsg string, errCode int) {
//	type t struct {
//		XXXX string `validate:"XXXX"`
//	}
//	return Validate(t{XXX: value.(string)})
//}
func Name(value interface{}) (errMsg string, errCode int) {
	type t struct {
		Name string `validate:"min=4,max=20"`
	}
	return Validate(t{Name: value.(string)})
}
func Phone(value interface{}) (errMsg string, errCode int) {

	type t struct {
		Phone string `validate:"required,number,len=11"`
	}
	return Validate(t{Phone: value.(string)})
}
func Birth(value interface{}) (errMsg string, errCode int) {
	type t struct {
		Birth string `validate:"datetime=02-01-2006"`
	}
	return Validate(t{Birth: value.(string)})
}

func Gender(value interface{}) (errMsg string, errCode int) {
	type t struct {
		Gender string `validate:"required,oneof=Female Male 'Not Specified'"`
	}
	return Validate(t{Gender: value.(string)})
}
func Bio(value interface{}) (errMsg string, errCode int) {
	type t struct {
		Bio string `validate:"max=50"`
	}
	return Validate(t{Bio: value.(string)})
}
func About(value interface{}) (errMsg string, errCode int) {
	type t struct {
		About string `validate:""`
	}
	return Validate(t{About: value.(string)})
}
func WorkExperience(value interface{}) (string, int) {
	type t struct {
		Num int `validate:"number,gte=0"`
	}
	return Validate(t{Num: value.(int)})
}
func CoinFunc(value interface{}) (string, int) {
	type t struct {
		Coin float64 `validate:"required,gte=0"`
	}
	return Validate(t{Coin: value.(float64)})
}
