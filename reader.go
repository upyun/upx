package main

import (
	"os"
)

type ProgressReader struct {
	fd    *os.File
	coyed int
}

func (r *ProgressReader) Len() int {
	fInfo, _ := r.fd.Stat()
	return int(fInfo.Size())
}

func (r *ProgressReader) MD5() string {
	return ""
}

func (r *ProgressReader) Read(b []byte) (n int, err error) {
	n, err = r.fd.Read(b)
	r.coyed += n
	return
}

func (r *ProgressReader) Copyed() int {
	return r.coyed
}
