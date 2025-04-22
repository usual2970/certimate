package fileutil

import (
	"os"
	"path/filepath"

	xerrors "github.com/pkg/errors"
)

// 与 [Write] 类似，但写入的是字符串内容。
//
// 入参:
//   - path: 文件路径。
//   - content: 文件内容。
//
// 出参:
//   - 错误。
func WriteString(path string, content string) error {
	return Write(path, []byte(content))
}

// 将数据写入指定路径的文件。
// 如果目录不存在，将会递归创建目录。
// 如果文件不存在，将会创建该文件；如果文件已存在，将会覆盖原有内容。
//
// 入参:
//   - path: 文件路径。
//   - data: 文件数据字节数组。
//
// 出参:
//   - 错误。
func Write(path string, data []byte) error {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return xerrors.Wrap(err, "failed to create directory")
	}

	file, err := os.Create(path)
	if err != nil {
		return xerrors.Wrap(err, "failed to create file")
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return xerrors.Wrap(err, "failed to write file")
	}

	return nil
}
