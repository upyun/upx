package upx

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

var (
	ROOT     = fmt.Sprintf("/upx-test/%s", time.Now())
	BUCKET_1 = os.Getenv("UPYUN_BUCKET1")
	BUCKET_2 = os.Getenv("UPYUN_BUCKET2")
	USERNAME = os.Getenv("UPYUN_USERNAME")
	PASSWORD = os.Getenv("UPYUN_PASSWORD")
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
