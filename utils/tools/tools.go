package tools

import (
	"encoding/json"
	"github.com/fatih/structs"
	"reflect"
	"runtime"
	"service-backend/utils"
	"service-backend/utils/errmsg"
	"strings"
	"unicode"
)

// LowFirst 首字母小写 SomeThing->someThing
func LowFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// Case2CamelCase 蛇形转驼峰 some_thing -> someThing
func Case2CamelCase(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	str = strings.Replace(str, " ", "", -1)
	return LowFirst(str)
}

// TransformDataSlice 把数据转换为小驼峰返回
func TransformDataSlice(data []map[string]interface{}) []map[string]interface{} {
	var res []map[string]interface{}
	for _, each := range data {
		res = append(res, TransformData(each))
	}
	return res
}

// TransformData 数据的key转化为小驼峰返回
func TransformData(data map[string]interface{}) map[string]interface{} {
	t := map[string]interface{}{}
	for k, v := range data {
		t[Case2CamelCase(k)] = v
	}
	return t
}
func TransformStruct(data interface{}) map[string]interface{} {
	err, res := Structs2Map(data)
	if err != nil {
		return nil
	}
	return TransformData(res)
}

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

func Structs2Map(data interface{}) (err error, res map[string]interface{}) {
	defer func() {
		if err != nil {
			//logger.Log.Error(" 结构体转化为map失败", zap.String("data", fmt.Sprintf("%v", data)))
		}
	}()
	var b []byte
	if b, err = json.Marshal(data); err != nil {
		return err, nil
	}
	err = json.Unmarshal(b, &res)
	return err, res
}

// ConvertCoinF2I 转化金币浮点->INT 入库
func ConvertCoinF2I(coin float32) int64 {
	return int64(coin * float32(utils.CoinBase))
}

// ConvertCoinI2F  转化金币INT-> 浮点 展示
func ConvertCoinI2F(coin int64) float32 {
	return float32(coin) / float32(utils.CoinBase)
}

// WhoCallMe 显示上一层调用该函数的方法名
func WhoCallMe() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

// Structs2SQLTable 将结构体实例中的带有structs tag字段的值提取为map
func Structs2SQLTable(s interface{}) map[string]interface{} {
	out := structs.Map(s)
	for k := range out {
		if unicode.IsUpper([]rune(k)[0]) {
			delete(out, k)
		}
	}
	return out
}
