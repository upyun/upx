package upyun

import (
	"fmt"
	"io"
	"os"
)

type UpYunPutReader interface {
	Len() (n int)
	MD5() (ret string)
	Read([]byte) (n int, err error)
	Copyed() (n int)
}

type fragmentFile struct {
	realFile *os.File
	offset   int64
	limit    int64
	cursor   int64
}

func (f *fragmentFile) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	case 0:
		f.cursor = offset
		ret, err = f.realFile.Seek(f.offset+f.cursor, 0)
		return ret - f.offset, err
	default:
		return 0, fmt.Errorf("whence must be 0")
	}
}

func (f *fragmentFile) Read(b []byte) (n int, err error) {
	if f.cursor >= f.limit {
		return 0, io.EOF
	}
	n, err = f.realFile.Read(b)
	if f.cursor+int64(n) > f.limit {
		n = int(f.limit - f.cursor)
	}
	f.cursor += int64(n)
	return n, err
}

func (f *fragmentFile) Stat() (fInfo os.FileInfo, err error) {
	return fInfo, fmt.Errorf("fragmentFile not implement Stat()")
}

func (f *fragmentFile) Close() error {
	return nil
}

func (f *fragmentFile) Copyed() int {
	return int(f.cursor - f.offset)
}

func (f *fragmentFile) Len() int {
	return int(f.limit - f.offset)
}

func (f *fragmentFile) MD5() string {
	s, _ := md5File(f)
	return s
}

func newFragmentFile(file *os.File, offset, limit int64) (*fragmentFile, error) {
	f := &fragmentFile{
		realFile: file,
		offset:   offset,
		limit:    limit,
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}
	return f, nil
}
