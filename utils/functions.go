package utils

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

const (
	//常用的时间格式串
	TimeFormat = "2006-01-02 15:04:05"
)

//算术相加，返回float64
func Add(params ...interface{}) (result float64) {
	if len(params) > 0 {
		result = 0
		for _, d := range params {
			result += AsFloat64(d)
		}
	}
	return
}

//算术相减，返回float64
func Div(params ...interface{}) (result float64) {
	if len(params) > 0 {
		result = AsFloat64(params[0])
		for i := 1; i < len(params); i++ {
			result -= AsFloat64(params[i])
		}
	}
	return
}

//最传入参数中的最大值，返回float64
func Max(params ...interface{}) (result float64) {
	if len(params) > 0 {
		result = AsFloat64(params[0])
		for i := 1; i < len(params); i++ {
			if d := AsFloat64(params[i]); d > result {
				result = d
			}
		}
	}
	return
}

//取传入参数中的最小值，返回float64
func Min(params ...interface{}) (result float64) {
	if len(params) > 0 {
		result = AsFloat64(params[0])
		for i := 1; i < len(params); i++ {
			if d := AsFloat64(params[i]); d < result {
				result = d
			}
		}
	}
	return
}

//将无指定类型的参数格式化为字符串，时间格式为标准格式，见常定义TimeFormat = "2006-01-02 15:04:05"
func AsString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int8, int16, int32, int64, int, uint, uint16, uint32, uint64, uint8:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case time.Time:
		return src.(time.Time).Format(TimeFormat)
	}
	return fmt.Sprintf("%v", src)
}

//将无指定类型的参数转化为float64
func AsFloat64(src interface{}) (result float64) {
	result = 0
	switch src.(type) {
	case float64:
		result = src.(float64)
	case float32:
		result = float64(src.(float32))
	case int8:
		result = float64(src.(int8))
	case int16:
		result = float64(src.(int16))
	case int32:
		result = float64(src.(int32))
	case int64:
		result = float64(src.(int64))
	case uint8:
		result = float64(src.(uint8))
	case uint16:
		result = float64(src.(uint16))
	case uint32:
		result = float64(src.(uint32))
	case uint64:
		result = float64(src.(uint64))
	default:
		if dc, err := strconv.ParseFloat(AsString(src), 10); err == nil {
			result = dc
		}
	}
	return
}

//将无指定类型的参数转化为int64
func AsInt64(src interface{}) (result int64) {
	result = 0
	switch src.(type) {
	case int8:
		result = int64(src.(int8))
	case int16:
		result = int64(src.(int16))
	case int32:
		result = int64(src.(int32))
	case int64:
		result = int64(src.(int64))
	case uint8:
		result = int64(src.(uint8))
	case uint16:
		result = int64(src.(uint16))
	case uint32:
		result = int64(src.(uint32))
	case uint64:
		result = int64(src.(uint64))
	default:
		if dc, err := strconv.ParseInt(AsString(src), 10, 64); err == nil {
			result = dc
		}
	}
	return
}

//将无指定类型的参数转化为int
func AsInt(src interface{}) (result int) {
	result = 0
	switch v := src.(type) {
	case int8:
		result = int(v)
	case int16:
		result = int(v)
	case int32:
		result = int(v)
	case int64:
		result = int(v)
	case uint8:
		result = int(v)
	case uint16:
		result = int(v)
	case uint32:
		result = int(v)
	case uint64:
		result = int(v)
	case float32:
		result = int(v)
	case float64:
		result = int(v)
	default:
		if dc, err := strconv.ParseInt(AsString(src), 10, 64); err == nil {
			result = int(dc)
		}
	}
	return
}

//对float64的数字截取小数，四舍五入
func Round(val float64, places int) float64 {
	var t float64
	f := math.Pow10(places)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return val
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f

	if !math.IsInf(x, 0) {
		return x
	}

	return t
}

//判断参数是否可转化为数字
func IsNumber(src interface{}) bool {
	switch v := src.(type) {
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		word := fmt.Sprint(v)
		if word[0] >= 48 && word[0] <= 57 {
			return true
		}
	}
	return false
}

//截断字符串
func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}
