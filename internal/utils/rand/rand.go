package rand

import (
	"math/rand"
	"time"
)

// RandStr 随机生成指定长度字符串
func RandStr(n int) string {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	random := rand.New(source)

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// 使用循环生成指定长度的字符串
	result := make([]rune, n)
	for i := range result {
		result[i] = letters[random.Intn(len(letters))]
	}

	return string(result)
}
