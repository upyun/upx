package main

import (
	"testing"
	"fmt"
)

func Test_ResumePut(t *testing.T) {
	b, err := Upx("resume-put")
	Nil(t, err)
	Equal(t, string(b), fmt.Sprintf("Done %s/%s ~~\n"))
}
