package xtime

import (
	"time"
)

func BeijingTimeStr() string {
	location, _ := time.LoadLocation("Asia/Shanghai")

	// 获取当前时间
	now := time.Now().In(location)

	// 格式化为字符串
	return now.Format("2006-01-02 15:04:05")
}

func GetTimeAfter(d time.Duration) string {

	t := time.Now().UTC()

	return t.Add(d).Format("2006-01-02 15:04:05")
}
