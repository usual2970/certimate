package typeutil

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
func ToVal[T any](p *T) (v T) {
	if IsNil(p) {
		return v
	}

	return *p
}
