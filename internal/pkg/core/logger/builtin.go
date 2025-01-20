package logger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/usual2970/certimate/internal/pkg/utils/types"
)

// 表示默认的日志记录器类型。
type DefaultLogger struct {
	records []string
}

var _ Logger = (*DefaultLogger)(nil)

func (l *DefaultLogger) Logt(tag string, data ...any) {
	l.ensureInitialized()

	temp := make([]string, len(data)+1)
	temp[0] = tag
	for i, v := range data {
		s := ""
		if types.IsNil(v) {
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

func (l *DefaultLogger) Logf(format string, args ...any) {
	l.ensureInitialized()

	l.records = append(l.records, fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) GetRecords() []string {
	l.ensureInitialized()

	temp := make([]string, len(l.records))
	copy(temp, l.records)
	return temp
}

func (l *DefaultLogger) FlushRecords() {
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

func (l *NilLogger) Logt(string, ...any) {}
func (l *NilLogger) Logf(string, ...any) {}
func (l *NilLogger) GetRecords() []string {
	return make([]string, 0)
}
func (l *NilLogger) FlushRecords() {}

func NewNilLogger() *NilLogger {
	return &NilLogger{}
}
