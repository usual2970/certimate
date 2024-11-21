package deployer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// 表示定义证书部署器的日志记录器的抽象类型接口。
type Logger interface {
	// 追加一条日志记录。
	// 该方法会将 `data` 以 JSON 序列化后拼接到 `tag` 结尾。
	//
	// 入参：
	//   - tag：标签。
	//   - data：数据。
	Appendt(tag string, data ...any)

	// 追加一条日志记录。
	// 该方法会将 `args` 以 `format` 格式化。
	//
	// 入参：
	//   - format：格式化字符串。
	//   - args：格式化参数。
	Appendf(format string, args ...any)

	// 获取所有日志记录。
	GetRecords() []string

	// 清空。
	Flush()
}

// 表示默认的日志记录器类型。
type DefaultLogger struct {
	records []string
}

var _ Logger = (*DefaultLogger)(nil)

func (l *DefaultLogger) Appendt(tag string, data ...any) {
	l.ensureInitialized()

	temp := make([]string, len(data)+1)
	temp[0] = tag
	for i, v := range data {
		s := ""
		if v == nil {
			s = "<nil>"
		} else {
			switch reflect.ValueOf(v).Kind() {
			case reflect.String:
				s = v.(string)
			case reflect.Bool,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				s = fmt.Sprintf("%v", v)
			default:
				jsonData, _ := json.Marshal(v)
				s = string(jsonData)
			}
		}

		temp[i+1] = s
	}

	l.records = append(l.records, strings.Join(temp, ": "))
}

func (l *DefaultLogger) Appendf(format string, args ...any) {
	l.ensureInitialized()

	l.records = append(l.records, fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) GetRecords() []string {
	l.ensureInitialized()

	temp := make([]string, len(l.records))
	copy(temp, l.records)
	return temp
}

func (l *DefaultLogger) Flush() {
	l.records = make([]string, 0)
}

func (l *DefaultLogger) ensureInitialized() {
	if l.records == nil {
		l.records = make([]string, 0)
	}
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		records: make([]string, 0),
	}
}

// 表示空的日志记录器类型。
// 该日志记录器不会执行任何操作。
type NilLogger struct{}

var _ Logger = (*NilLogger)(nil)

func (l *NilLogger) Appendt(string, ...any) {}
func (l *NilLogger) Appendf(string, ...any) {}
func (l *NilLogger) GetRecords() []string {
	return make([]string, 0)
}
func (l *NilLogger) Flush() {}

func NewNilLogger() *NilLogger {
	return &NilLogger{}
}
