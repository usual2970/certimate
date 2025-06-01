package typeutil

import "reflect"

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

// 将非零值的对象转换为指针。
// 与 [ToPtr] 不同的是，如果对象的值为零值，则返回 nil。
//
// 入参：
//   - 待转换的对象。
//
// 出参：
//   - 返回对象的指针。
func ToPtrOrZeroNil[T any](v T) (p *T) {
	if reflect.ValueOf(v).IsZero() {
		return nil
	}

	return &v
}

// 将指针转换为对象。
//
// 入参：
//   - 待转换的指针。
//
// 出参：
//   - 返回指针指向的对象。如果指针为空，则返回对象的零值。
func ToVal[T any](p *T) (v T) {
	if IsNil(p) {
		return v
	}

	return *p
}
