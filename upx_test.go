package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var (
	username string
	password string
	bucket   string
	tmpPath  string
)

func init() {
	username = os.Getenv("username")
	password = os.Getenv("password")
	bucket = os.Getenv("bucket")
	tmpPath = fmt.Sprintf("/upx/%d", time.Now().Unix())
	path := os.Getenv("PATH")
	pwd, _ := os.Getwd()

	fmt.Println(username, password, bucket)

	os.Setenv("PATH", path+":"+pwd)
}

func upx(cmd string, args ...string) ([]byte, error) {
	args = append([]string{cmd}, args...)
	return exec.Command("./upx", args...).Output()
}

func TestLogin(t *testing.T) {
	_, err := upx("login", bucket, username, password)
	if err != nil {
		t.Errorf("failed to upx")
		t.FailNow()
	}
}

func TestMkdir(t *testing.T) {
	_, err := upx("mkdir", tmpPath)
	if err != nil {
		t.Errorf("failed to upx")
		t.FailNow()
	}
}

func TestCd(t *testing.T) {
	_, err := upx("cd", tmpPath)
	if err != nil {
		t.Errorf("failed to upx")
		t.Fail()
	}
}

func TestPwd(t *testing.T) {
	b, err := upx("pwd")
	if err != nil {
		t.Errorf("failed to upx")
		t.Fail()
	}
	if string(b) != tmpPath+"\n" {
		t.Errorf("%s != %s\n", string(b), tmpPath)
		t.Fail()
	}
}

func TestPut(t *testing.T) {
	_, err := upx("put", "upx.go")
	if err != nil {
		t.Errorf("failed to upx")
		t.Fail()
	}
}

func TestLs(t *testing.T) {
	b, err := upx("ls")
	if err != nil {
		t.Errorf("failed to upx")
		t.FailNow()
	}

	if !strings.Contains(string(b), "upx.go") {
		t.FailNow()
	}

	b1, err := upx("ls", tmpPath)
	if err != nil {
		t.Errorf("failed to upx")
		t.FailNow()
	}

	if string(b) != string(b1) {
		t.FailNow()
	}
}

func TestGet(t *testing.T) {
	_, err := upx("get", "upx.go", "upx.go.2")
	if err != nil {
		t.Errorf("failed to upx")
		t.FailNow()
	}

	_, err = os.Lstat("upx.go.2")
	if err != nil {
		t.Errorf("failed to upx")
		t.FailNow()
	}
}

func TestRm(t *testing.T) {
	_, err := upx("rm", "*.go")
	if err != nil {
		t.Errorf("failed to upx")
		t.FailNow()
	}
	b, _ := upx("ls", tmpPath)
	if string(b) != "" {
		t.FailNow()
	}
}

func TestServices(t *testing.T) {
	b, err := upx("services")
	if err != nil {
		t.Errorf("failed to upx")
		t.FailNow()
	}
	if !strings.Contains(string(b), bucket) {
		t.FailNow()
	}

	b1, err := upx("sv")
	if err != nil {
		t.Errorf("failed to upx")
		t.FailNow()
	}
	if string(b) != string(b1) {
		t.Errorf("%s != %s\n", string(b), string(b1))
		t.Fail()
	}
}

func TestSwitch(t *testing.T) {
	_, err := upx("switch", bucket)
	if err != nil {
		t.Errorf("failed to upx")
		t.Fail()
	}
}

func TestLogout(t *testing.T) {
	_, err := upx("logout")
	if err != nil {
		t.Errorf("failed to upx")
		t.Fail()
	}
}
