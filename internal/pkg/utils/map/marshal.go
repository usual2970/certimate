package maputil

import (
	mapstructure "github.com/go-viper/mapstructure/v2"
)

// 将字典填充到指定类型的结构体。
// 与 [json.Unmarshal] 类似，但传入的是一个 [map[string]any] 对象而非 JSON 格式的字符串。
//
// 入参：
//   - dict: 字典。
//   - output: 结构体指针。
//
// 出参：
//   - 错误信息。如果填充失败，则返回错误信息。
func Populate(dict map[string]any, output any) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		TagName:          "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(dict)
}
