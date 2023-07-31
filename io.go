package upx

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/upyun/upx/processbar"
)

var (
	IsVerbose = true
	mu        = &sync.Mutex{}
)

func NewFileWrappedWriter(localPath string, bar *processbar.UpxProcessBar, resume bool) (io.WriteCloser, error) {
	var fd *os.File
	var err error
	if resume {
		fd, err = os.OpenFile(localPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	} else {
		fd, err = os.Create(localPath)
	}
	if err != nil {
		return nil, err
	}

	fileinfo, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	if bar != nil {
		bar.SetCurrent(fileinfo.Size())
	}
	return bar.NewProxyWriter(fd), nil
}

func NewFileWrappedReader(bar *processbar.UpxProcessBar, fd io.Reader) io.ReadCloser {
	return bar.NewProxyReader(fd)
}

func Print(arg0 string, args ...interface{}) {
	s := arg0 //arg0 may include '%'
	if len(args) > 0 {
		s = fmt.Sprintf(arg0, args...)
	}
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	mu.Lock()
	os.Stdout.WriteString(s)
	mu.Unlock()
}

func PrintOnlyVerbose(arg0 string, args ...interface{}) {
	if IsVerbose {
		Print(arg0, args...)
	}
}

func PrintError(arg0 string, args ...interface{}) {
	s := fmt.Sprintf(arg0, args...)
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	mu.Lock()
	os.Stderr.WriteString(s)
	mu.Unlock()
}

func PrintErrorAndExit(arg0 string, args ...interface{}) {
	PrintError(arg0, args...)
	os.Exit(-1)
}
