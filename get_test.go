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
	compareGet(t, src, correct, start, end)
	Nil(t, err)
}

func TestGet(t *testing.T) {
	tpath, _ := os.Getwd()
	testdir := filepath.Join(tpath, "test-get")
	base := ROOT + "/get/"
	start := base + "FILE"
	end := base + "putfile"
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

	// upx put localBase/FILE upBase/FILE
	CreateFile("FILE")
	putFile(t, filepath.Join(localBase, "FILE"), "", path.Join(base, "FILE"))

	// upx put ../put/FILE2
	CreateFile("FILE2")
	localPath := filepath.Join(localBase, "FILE2")
	putFile(t, localPath, "", path.Join(base, "FILE2"))

	// upx put /path/to/file /path/to/file
	putFile(t, "FILE", path.Join(base, "FILE4"), path.Join(base, "FILE4"))

	// upx put /path/to/file /path/to/dir
	CreateFile("FILE3")
	putFile(t, "FILE3", base, path.Join(base, "FILE3"))

	// upx put /path/to/file ../path/to/dir/
	putFile(t, "FILE", base+"/putfile/", path.Join(base, "putfile", "FILE"))

	// upx put /path/to/dir /path/to/dir/
	putDir(t, localBase, base+"/putdir/", base+"/putdir/")
	_, err = Upx("put", localBase, path.Join(base, "FILE"))
	NotNil(t, err)
	err = os.MkdirAll(testdir, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	GetStartBetweenEndFiles(t, base, testdir, testdir, start, end)
}

func compareGet(t *testing.T, up, local, start, end string) {
	locals := []string{}
	ups := []string{}

	lpath := local
	var localPath func(up string)
	localPath = func(local string) {
		fInfos, _ := ioutil.ReadDir(local + "/")
		for _, fInfo := range fInfos {
			fp := filepath.Join(local, fInfo.Name())
			locals = append(locals, up[:len(up)-1]+fp[len(lpath):])
			if IsDir(fp) {
				localPath(fp)
			}
		}
	}
	localPath(lpath)

	var upPath func(up string)
	upPath = func(up string) {
		//log.Println("upPath:", up)
		b, err := Upx("ls", up)
		Nil(t, err)
		output := strings.TrimRight(string(b), "\n")
		for _, line := range strings.Split(output, "\n") {
			//log.Println(line)
			items := strings.Split(line, " ")
			fp := filepath.Join(up, items[len(items)-1])
			if fp >= start && fp < end {
				ups = append(ups, fp)
				if items[0][0] == 'd' {
					upPath(fp)
				}
			} else if strings.HasPrefix(start, fp) {
				if items[0][0] == 'd' {
					upPath(fp)
				}
			}
		}
	}
	upPath(up)

	sort.Strings(locals)
	sort.Strings(ups)
	Equal(t, len(locals), len(ups))
	for i := 0; i < len(locals); i++ {
		log.Println("compare:", locals[i], " ", ups[i])
		Equal(t, locals[i], ups[i])
	}
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
