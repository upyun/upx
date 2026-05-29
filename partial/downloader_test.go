package partial

import (
	"bytes"
	"crypto/md5"
	"strings"
	"testing"
)

func TestDownload(t *testing.T) {
	var buffer bytes.Buffer

	filedata := []byte(strings.Repeat("hello world", 1024*100))
	download := NewMultiPartialDownloader(
		"myTestfile",
		int64(len(filedata)),
		1024,
		&buffer,
		3,
		func(start, end int64) ([]byte, error) {
			return filedata[start : end+1], nil
		},
	)

	err := download.Download()
	if err != nil {
		t.Fatal(err.Error())
	}
	if md5.Sum(buffer.Bytes()) != md5.Sum(filedata) {
		t.Fatal("download file has diff MD5")
	}
}
