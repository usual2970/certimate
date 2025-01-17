package types

import "reflect"

// 判断对象是否为 nil。
//
// 入参：
//   - value：待判断的对象。
//
// 出参：
//   - 如果对象值为 nil，则返回 true；否则返回 false。
func IsNil(obj any) bool {
	if obj == nil {
		return true
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	} else if v.Kind() == reflect.Interface {
		return v.Elem().IsNil()
	}

	return false
}
