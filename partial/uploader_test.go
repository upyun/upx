package partial

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploaderError(t *testing.T) {
	filedata := []byte(strings.Repeat("hello world", 1024*100))
	uploader := NewMultiPartialUploader(
		1,
		bytes.NewReader(filedata),
		3,
		func(partId int64, body []byte) error {
			if partId > 20 {
				return errors.New("error")
			}
			return nil
		},
	)

	err := uploader.Upload()
	assert.Error(t, err)
	assert.Equal(t, err, errors.New("error"))
}

func TestUploader(t *testing.T) {
	filedata := []byte(strings.Repeat("hello world", 24))

	parts := make(map[int64][]byte)
	var mutex sync.RWMutex
	uploader := NewMultiPartialUploader(
		10,
		bytes.NewReader(filedata),
		3,
		func(partId int64, body []byte) error {
			mutex.Lock()
			parts[partId] = body
			mutex.Unlock()
			return nil
		},
	)

	err := uploader.Upload()
	assert.NoError(t, err)

	// 组合结果
	var buffer bytes.Buffer
	for i := 0; i <= len(filedata)/10; i++ {
		d, ok := parts[int64(i)]
		assert.Equal(t, ok, true)
		buffer.Write(d)
	}

	assert.Equal(t, filedata, buffer.Bytes())
}
