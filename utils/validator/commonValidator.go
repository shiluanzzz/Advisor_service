package validator

import (
	"fmt"
	"reflect"
	"service/utils/errmsg"
)

func CallFunc(m map[string]interface{}, name string, params ...interface{}) (string, int) {
	if m[name] == nil {
		return fmt.Sprintf("不存在字段%s校验函数", name), errmsg.ERROR
	}
	f := reflect.ValueOf(m[name])
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result := f.Call(in)
	return result[0].String(), int(result[1].Int())
}

// 校验函数模板
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
	if len(value.(string)) != 11 {
		return "手机号的长度不等于11", errmsg.ERROR
	}
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

// TODO 失效
func Gender(value interface{}) (errMsg string, errCode int) {
	fmt.Println(value)
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
		about string `validate:""`
	}
	return Validate(t{about: value.(string)})
}
func WorkExperience(value interface{}) (string, int) {
	type t struct {
		num int `validate:"number,get=0"`
	}
	return Validate(t{num: value.(int)})
}
