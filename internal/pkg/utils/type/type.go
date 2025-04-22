package typeutil

import "reflect"

// 判断对象是否为 nil。
// 与直接使用 `obj == nil` 不同，该函数会正确判断接口类型对象的真实值是否为空。
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

// 将对象转换为指针。
//
// 入参：
//   - 待转换的对象。
//
// 出参：
//   - 返回对象的指针。
func ToPtr[T any](v T) (p *T) {
	return &v
}

// 将指针转换为对象。
//
// 入参：
//   - 待转换的指针。
//
// 出参：
//   - 返回指针指向的对象。如果指针为空，则返回对象的零值。
func ToObj[T any](p *T) (v T) {
	if p == nil {
		return v
	}

	return *p
}
