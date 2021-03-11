package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/gosuri/uiprogress"
)

var (
	isVerbose = true
	mu        = &sync.Mutex{}
)

type WrappedWriter struct {
	w      io.WriteCloser
	Copyed int
	bar    *uiprogress.Bar
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

func NewFileWrappedWriter(localPath string, bar *uiprogress.Bar) (*WrappedWriter, error) {
	fd, err := os.Create(localPath)
	if err != nil {
		return nil, err
	}

	return &WrappedWriter{
		w:      fd,
		Copyed: 0,
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
