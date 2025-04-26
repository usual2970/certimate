package httputil

import (
	"bufio"
	"net/http"
	"net/textproto"
	"strings"
)

// 从表示 HTTP 标头的字符串解析并返回一个 http.Header 对象。
//
// 入参:
//   - headers: 表示 HTTP 标头的字符串。
//
// 出参:
//   - header: http.Header 对象。
//   - err: 错误。
func ParseHeaders(headers string) (http.Header, error) {
	str := strings.TrimSpace(headers) + "\r\n\r\n"
	if len(str) == 4 {
		return make(http.Header), nil
	}

	br := bufio.NewReader(strings.NewReader(str))
	tp := textproto.NewReader(br)

	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}

	return http.Header(mimeHeader), err
}
