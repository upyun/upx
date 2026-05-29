//go:build windows

package fsutil

import (
	"io/fs"
	"path/filepath"
	"syscall"
)

// 判断文件是否是需要忽略的文件
func IsIgnoreFile(path string, fileInfo fs.FileInfo) bool {
	for hasDotPrefix(filepath.Base(path)) {
		return true
	}

	underlyingData := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	if underlyingData != nil {
		return underlyingData.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
	}

	return false
}
