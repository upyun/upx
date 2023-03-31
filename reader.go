package main

import (
	"io"
)

type ProgressReader struct {
	reader io.Reader
	coyed  int
}

func (r *ProgressReader) Read(b []byte) (n int, err error) {
	n, err = r.reader.Read(b)
	r.coyed += n
	return
}

func (r *ProgressReader) Copyed() int {
	return r.coyed
}
