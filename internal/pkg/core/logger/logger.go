package logger

// 表示定义日志记录器的抽象类型接口。
type Logger interface {
	// 追加一条日志记录。
	// 该方法会将 `data` 以 JSON 序列化后拼接到 `tag` 结尾。
	//
	// 入参：
	//   - tag：标签。
	//   - data：数据。
	Logt(tag string, data ...any)

	// 追加一条日志记录。
	// 该方法会将 `args` 以 `format` 格式化。
	//
	// 入参：
	//   - format：格式化字符串。
	//   - args：格式化参数。
	Logf(format string, args ...any)

	// 获取所有日志记录。
	// TODO: 记录时间
	GetRecords() []string

	// 清空所有日志记录。
	FlushRecords()
}
