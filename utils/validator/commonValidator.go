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

const msg = "数据类型不为 %s"

func Name(value interface{}) (errMsg string, errCode int) {
	type t struct {
		Name string `validate:"min=4,max=20"`
	}
	if name, ok := value.(string); ok {
		return Validate(t{Name: name})
	} else {
		return fmt.Sprintf(msg, "string"), errmsg.ErrorInput
	}
}
func Phone(value interface{}) (errMsg string, errCode int) {

	type t struct {
		Phone string `validate:"required,number,len=11"`
	}
	if phone, ok := value.(string); ok {
		return Validate(t{Phone: phone})
	} else {
		return fmt.Sprintf(msg, "string"), errmsg.ErrorInput
	}
}
func Birth(value interface{}) (errMsg string, errCode int) {
	type t struct {
		Birth string `validate:"datetime=02-01-2006"`
	}
	if valueTrue, ok := value.(string); ok {
		return Validate(t{Birth: valueTrue})
	} else {
		return fmt.Sprintf(msg, "string"), errmsg.ErrorInput
	}
}

func Gender(value interface{}) (errMsg string, errCode int) {
	type t struct {
		Gender int `validate:"required,number,min=1,max=3"`
	}
	if valueTrue, ok := value.(float64); ok {
		return Validate(t{Gender: int(valueTrue)})
	} else {
		return fmt.Sprintf(msg, "int"), errmsg.ErrorInput
	}
}
func Bio(value interface{}) (errMsg string, errCode int) {
	type t struct {
		Bio string `validate:"max=50"`
	}
	if valueTrue, ok := value.(string); ok {
		return Validate(t{Bio: valueTrue})
	} else {
		return fmt.Sprintf(msg, "string"), errmsg.ErrorInput
	}
}
func About(value interface{}) (errMsg string, errCode int) {
	type t struct {
		About string `validate:""`
	}
	if valueTrue, ok := value.(string); ok {
		return Validate(t{About: valueTrue})
	} else {
		return fmt.Sprintf(msg, "string"), errmsg.ErrorInput
	}
}
func WorkExperience(value interface{}) (string, int) {
	type t struct {
		Num int `validate:"number,gte=0"`
	}
	if valueTrue, ok := value.(int); ok {
		return Validate(t{Num: valueTrue})
	} else {
		return fmt.Sprintf(msg, "int"), errmsg.ErrorInput
	}
}
func CoinFunc(value interface{}) (string, int) {
	type t struct {
		Coin float64 `validate:"required,gte=0"`
	}
	if valueTrue, ok := value.(float64); ok {
		return Validate(t{Coin: valueTrue})
	} else {
		return fmt.Sprintf(msg, "int"), errmsg.ErrorInput
	}
}
