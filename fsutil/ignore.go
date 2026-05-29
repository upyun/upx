package fsutil

import (
	"strings"
)

// 判断文件是否是 . 开头的
func hasDotPrefix(filename string) bool {
	return strings.HasPrefix(filename, ".")
}
