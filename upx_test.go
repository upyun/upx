package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

var (
	ROOT     = fmt.Sprintf("/upx-test/%s", time.Now())
	BUCKET_1 = "bigfile"
	BUCKET_2 = "prog-test"
	USERNAME = "myworker"
	PASSWORD = "TYGHBNtyghbn"
)

func SetUp() {
	Upx("login", BUCKET_1, USERNAME, PASSWORD)
}

func TearDown() {
	for {
		b, err := Upx("logout")
		if err != nil || string(b) == "Nothing to do ~\n" {
			break
		}
	}
}

func Equal(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\tnexp: %#v\n\n\tgot:  %#v\033[39m\n\n",
			filepath.Base(file), line, expected, actual)
		t.FailNow()
	}
}

func NotEqual(t *testing.T, actual, expected interface{}) {
	if reflect.DeepEqual(actual, expected) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\tnexp: %#v\n\n\tgot:  %#v\033[39m\n\n",
			filepath.Base(file), line, expected, actual)
		t.FailNow()
	}
}
func Nil(t *testing.T, object interface{}) {
	if !isNil(object) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\t   <nil> (expected)\n\n\t!= %#v (actual)\033[39m\n\n",
			filepath.Base(file), line, object)
		t.FailNow()
	}
}

func NotNil(t *testing.T, object interface{}) {
	if isNil(object) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\tExpected value not to be <nil>\033[39m\n\n",
			filepath.Base(file), line, object)
		t.FailNow()
	}
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false

}

func CreateFile(fpath string) {
	os.MkdirAll(filepath.Dir(fpath), 0755)
	fd, _ := os.Create(fpath)
	fd.WriteString("UPX")
	fd.Close()
}

func Upx(args ...string) ([]byte, error) {
	cmd := exec.Command("upx", args...)
	var obuf, ebuf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &obuf, &ebuf
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	err := cmd.Wait()
	ob, _ := ioutil.ReadAll(&obuf)
	eb, _ := ioutil.ReadAll(&ebuf)
	if err != nil {
		return ob, fmt.Errorf("%s", string(eb))
	}
	return ob, nil
}

func TestMain(m *testing.M) {
	pwd, _ := os.Getwd()
	os.Setenv("PATH", pwd)
	flag.Parse()
	code := m.Run()
	os.Exit(code)
}
