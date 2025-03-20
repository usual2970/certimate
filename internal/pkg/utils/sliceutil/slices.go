package sliceutil

// 创建给定切片一部分的浅拷贝，其包含通过所提供函数实现的测试的所有元素。
//
// 入参：
//   - slice: 原切片。
//   - filter: 为切片中的每个元素执行的函数。它应该返回一个布尔值以指使是否将元素保留在结果切片中。
//
// 出参：
//   - 新切片。
func Filter[T any](slice []T, filter func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if filter(item) {
			result = append(result, item)
		}
	}
	return result
}

// 创建一个新切片，这个新切片由原切片中的每个元素都调用一次提供的函数后的返回值组成。
//
// 入参：
//   - slice: 原切片。
//   - mapper: 为切片中的每个元素执行的函数。它的返回值作为一个元素被添加为新切片中。
//
// 出参：
//   - 新切片。
func Map[T1 any, T2 any](slice []T1, mapper func(T1) T2) []T2 {
	result := make([]T2, 0, len(slice))
	for _, item := range slice {
		result = append(result, mapper(item))
	}
	return result
}

// 测试切片中是否每个元素都通过了由提供的函数实现的测试。
//
// 入参：
//   - slice: 切片。
//   - condition: 为切片中的每个元素执行的函数。它应该返回一个布尔值以指示元素是否通过测试。
//
// 出参：
//   - 测试结果。
func Every[T any](slice []T, condition func(T) bool) bool {
	for _, item := range slice {
		if !condition(item) {
			return false
		}
	}
	return true
}

// 测试切片中是否至少有一个元素通过了由提供的函数实现的测试。
//
// 入参：
//   - slice: 切片。
//   - condition: 为切片中的每个元素执行的函数。它应该返回一个布尔值以指示元素是否通过测试。
//
// 出参：
//   - 测试结果。
func Some[T any](slice []T, condition func(T) bool) bool {
	for _, item := range slice {
		if condition(item) {
			return true
		}
	}
	return false
}
