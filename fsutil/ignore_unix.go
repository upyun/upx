//go:build linux || darwin

package fsutil

import (
	"io/fs"
	"path/filepath"
)

// 判断文件是否是需要忽略的文件
func IsIgnoreFile(path string, fileInfo fs.FileInfo) bool {
	return hasDotPrefix(filepath.Base(path))
}
