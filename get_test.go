package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func GetStartBetweenEndFiles(t *testing.T, src, dst, correct, start, end string) {
	var err error
	src = AbsPath(src)

	if start != "" && start[0] != '/' {
		start = filepath.Join(src, start)
	}
	if end != "" && end[0] != '/' {
		end = filepath.Join(src, end)
	}
	if dst == "" {
		_, err = Upx("get", src, "--start="+start, "--end="+end)
	} else {
		_, err = Upx("get", src, dst, "--start="+start, "--end="+end)
	}
	Nil(t, err)
}

func TestGetStartBetweenEndFiles(t *testing.T) {
	tpath, _ := os.Getwd()
	base := ROOT + "/get/"
	pwd, err := ioutil.TempDir("", "test")
	Nil(t, err)
	localBase := filepath.Join(pwd, "get")

	func() {
		SetUp()
		err := os.MkdirAll(localBase, 0755)
		Nil(t, err)
	}()
	defer TearDown()

	err = os.Chdir(localBase)
	Nil(t, err)
	Upx("mkdir", base)
	Upx("cd", base)

	type uploadFiles []struct {
		name    string
		file    string
		dst     string
		correct string
	}
	type uploadDirs []struct {
		dir     string
		dst     string
		correct string
	}
	files := uploadFiles{
		{name: "111", file: filepath.Join(localBase, "111"), dst: "", correct: filepath.Join(base, "111")},
		{name: "333", file: filepath.Join(localBase, "333"), dst: "", correct: path.Join(base, "333")},
		{name: "333", file: "333", dst: path.Join(base, "333"), correct: path.Join(base, "333")},
		{name: "777", file: "777", dst: base, correct: path.Join(base, "777")},
		{name: "666", file: "666", dst: base + "/444/", correct: path.Join(base, "444", "666")},
	}
	for _, file := range files {
		CreateFile(file.name)
		putFile(t, file.file, file.dst, file.correct)
	}
	log.Println(122)

	dirs := uploadDirs{
		{dir: localBase, dst: base + "/666/", correct: base + "/666/"},
	}
	for _, dir := range dirs {
		putDir(t, dir.dir, dir.dst, dir.correct)
	}

	type list struct {
		start   string
		end     string
		testDir string
	}
	type test struct {
		input list
		real  []string
		want  []string
	}

	tests := []test{
		{input: list{start: "123", end: "999", testDir: filepath.Join(tpath, "test1")}, real: localFile("test1", base), want: upFile(t, base, "123", "999")},
		{input: list{start: "0", end: "999", testDir: filepath.Join(tpath, "test2")}, real: localFile("test2", base), want: upFile(t, base, "444", "666")},
	}
	for _, tc := range tests {
		input := tc.input

		err = os.MkdirAll(input.testDir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}

		GetStartBetweenEndFiles(t, base, input.testDir, input.testDir, input.start, input.end)

		sort.Strings(tc.real)
		sort.Strings(tc.want)
		Equal(t, len(tc.real), len(tc.want))

		for i := 0; i < len(tc.real); i++ {
			log.Println("compare:", tc.real[i], " ", tc.want[i])
			Equal(t, tc.real[i], tc.want[i])
		}
	}
}

func localFile(local, up string) []string {
	var locals []string
	localLen := len(local)
	fInfos, _ := ioutil.ReadDir(local + "/")
	for _, fInfo := range fInfos {
		fp := filepath.Join(local, fInfo.Name())
		locals = append(locals, up[:len(up)-1]+fp[localLen:])
		if IsDir(fp) {
			localFile(fp, up)
		}
	}
	return locals
}

func upFile(t *testing.T, up, start, end string) []string {
	b, err := Upx("ls", up)
	Nil(t, err)

	var ups []string
	output := strings.TrimRight(string(b), "\n")
	for _, line := range strings.Split(output, "\n") {
		items := strings.Split(line, " ")
		fp := filepath.Join(up, items[len(items)-1])
		if fp >= start && fp < end {
			ups = append(ups, fp)
			if items[0][0] == 'd' {
				upFile(t, fp, start, end)
			}
		} else if strings.HasPrefix(start, fp) {
			if items[0][0] == 'd' {
				upFile(t, fp, start, end)
			}
		}
	}
	return ups
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {

		return false
	}
	return s.IsDir()

}

func AbsPath(upPath string) (ret string) {
	if strings.HasPrefix(upPath, "/") {
		ret = path.Join(upPath)
	} else {
		ret = path.Join("/", upPath)
	}
	if strings.HasSuffix(upPath, "/") && ret != "/" {
		ret += "/"
	}
	return
}
