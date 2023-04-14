package upx

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginAndLogout(t *testing.T) {

	b, err := Upx("login", BUCKET_1, USERNAME, PASSWORD)
	assert.NoError(t, err)
	assert.Equal(t, string(b), fmt.Sprintf("Welcome to %s, %s!\n", BUCKET_1, USERNAME))

	b, err = Upx("login", BUCKET_2, USERNAME, PASSWORD)
	assert.NoError(t, err)
	assert.Equal(t, string(b), fmt.Sprintf("Welcome to %s, %s!\n", BUCKET_2, USERNAME))

	b, err = Upx("logout")
	assert.NoError(t, err)
	assert.Equal(t, string(b), fmt.Sprintf("Goodbye %s/%s ~~\n", USERNAME, BUCKET_2))

	b, err = Upx("logout")
	assert.NoError(t, err)
	assert.Equal(t, string(b), fmt.Sprintf("Goodbye %s/%s ~~\n", USERNAME, BUCKET_1))
}

func TestGetInfo(t *testing.T) {
	SetUp()
	defer TearDown()
	pwd, _ := Upx("pwd")
	b, err := Upx("info")
	assert.NoError(t, err)
	s := []string{
		"ServiceName:   " + BUCKET_1,
		"Operator:      " + USERNAME,
		"CurrentDir:    " + strings.TrimRight(string(pwd), "\n"),
		"Usage:         ",
	}
	assert.Equal(t, strings.HasPrefix(string(b), strings.Join(s, "\n")), true)
}

func TestSessionsAndSwitch(t *testing.T) {
	SetUp()
	defer TearDown()
	b, err := Upx("sessions")
	assert.NoError(t, err)
	assert.Equal(t, string(b), fmt.Sprintf("> %s\n", BUCKET_1))

	Upx("login", BUCKET_2, USERNAME, PASSWORD)
	b, err = Upx("sessions")
	assert.NoError(t, err)
	assert.Equal(t, string(b), fmt.Sprintf("  %s\n> %s\n", BUCKET_1, BUCKET_2))

	Upx("switch", BUCKET_1)
	b, err = Upx("sessions")
	assert.NoError(t, err)
	assert.Equal(t, string(b), fmt.Sprintf("> %s\n  %s\n", BUCKET_1, BUCKET_2))

	pwd, _ := Upx("pwd")
	b, err = Upx("info")
	assert.NoError(t, err)
	s := []string{
		"ServiceName:   " + BUCKET_1,
		"Operator:      " + USERNAME,
		"CurrentDir:    " + strings.TrimRight(string(pwd), "\n"),
		"Usage:         ",
	}
	assert.Equal(t, strings.HasPrefix(string(b), strings.Join(s, "\n")), true)
}

// TODO
func TestAuth(t *testing.T) {
}

func TestPurge(t *testing.T) {
	SetUp()
	defer TearDown()
	b, err := Upx("purge", fmt.Sprintf("http://%s.b0.upaiyun.com/test.jpg", BUCKET_1))
	assert.NoError(t, err)
	assert.Equal(t, len(b), 0)

	_, err = Upx("purge", "http://www.baidu.com")
	if !assert.Error(t, err) {
		assert.Fail(t, "purge not has return error")
	}
	assert.Equal(t, err.Error(), "Purge failed urls:\nhttp://www.baidu.com\ntoo many fails\n")

	fd, _ := os.Create("list")
	fd.WriteString(fmt.Sprintf("http://%s.b0.upaiyun.com/test.jpg\n", BUCKET_1))
	fd.WriteString(fmt.Sprintf("http://%s.b0.upaiyun.com/测试.jpg\n", BUCKET_1))
	fd.WriteString(fmt.Sprintf("http://%s.b0.upaiyun.com/%%E5%%8F%%88%%E6%%8B%%8D%%E4%%BA%%91.jpg\n", BUCKET_1))
	fd.Close()

	b, err = Upx("purge", "--list", "list")
	assert.NoError(t, err)
	assert.Equal(t, len(b), 0)
}
