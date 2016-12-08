package main

import (
	"fmt"
	"os"
	"strings"
)

var DEBUG int = 0

func LogI(arg0 string, args ...interface{}) {
	s := fmt.Sprintf(arg0, args...)
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	os.Stdout.WriteString(s)
}

func LogD(arg0 string, args ...interface{}) {
	if DEBUG == 1 {
		LogI(arg0, args...)
	}
}

func LogE(arg0 string, args ...interface{}) {
	s := fmt.Sprintf(arg0, args...)
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	os.Stderr.WriteString(s)
}

func LogC(arg0 string, args ...interface{}) {
	s := fmt.Sprintf(arg0, args...)
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	os.Stderr.WriteString(s)
	os.Exit(-1)
}

func SetLogDebug() {
	DEBUG = 1
}
