package upx

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/arrebole/progressbar"
)

var (
	isVerbose = true
	mu        = &sync.Mutex{}
)

type WrappedWriter struct {
	w      io.WriteCloser
	Copyed int
	bar    *progressbar.ProgressBar
}

func (w *WrappedWriter) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	w.Copyed += n
	if w.bar != nil {
		w.bar.Set(w.Copyed)
	}
	return n, err
}

func (w *WrappedWriter) Close() error {
	return w.w.Close()
}

type WrappedReader struct {
	r      io.Reader
	Copyed int
	bar    *progressbar.ProgressBar
}

func (r *WrappedReader) Read(b []byte) (int, error) {
	n, err := r.r.Read(b)
	r.Copyed += n
	if r.bar != nil {
		r.bar.Set(r.Copyed)
	}
	return n, err
}

func NewFileWrappedWriter(localPath string, bar *progressbar.ProgressBar, resume bool) (*WrappedWriter, error) {
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
		bar.SetOffset64(fileinfo.Size())
	}
	return &WrappedWriter{
		w:      fd,
		Copyed: int(fileinfo.Size()),
		bar:    bar,
	}, nil
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
	if isVerbose {
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
