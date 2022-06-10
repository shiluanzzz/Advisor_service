package tools

import (
	"github.com/fatih/structs"
	"reflect"
	"runtime"
	"service-backend/utils/errmsg"
	"service-backend/utils/setting"
	"unicode"
)

// StructToMap 结构体转为Map[string]interface{},忽略nil指针
func StructToMap(in interface{}, tagName string) (map[string]interface{}, int) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, errmsg.ERROR
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			// 如果这个指向的是一个空指针就不用添加到map里去。
			if !v.Field(i).IsNil() {
				out[tagValue] = v.Field(i).Interface()
			}
		}
	}
	return out, errmsg.SUCCESS
}

// ConvertCoinF2I 转化金币浮点->INT 入库
func ConvertCoinF2I(coin float32) int64 {
	return int64(coin * float32(setting.ServiceCfg.CoinBase))
}

// ConvertCoinI2F  转化金币INT-> 浮点 展示
func ConvertCoinI2F(coin int64) float32 {
	return float32(coin) / float32(setting.ServiceCfg.CoinBase)
}

// WhoCallMe 显示上一层调用该函数的方法名
func WhoCallMe() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

// WhoAmI 我是谁往上跳一层
func WhoAmI() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

// Structs2SQLTable 将结构体实例中的带有structs tag字段的值提取为map
func Structs2SQLTable(s interface{}) map[string]interface{} {
	out := structs.Map(s)

	// key为大写开头，说明不包含structs这个tag，这个item不需要入库
	for k := range out {
		if unicode.IsUpper([]rune(k)[0]) {
			delete(out, k)
		}
	}
	return out
}
