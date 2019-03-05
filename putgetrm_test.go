package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func putFile(t *testing.T, src, dst, correct string) {
	var err error
	if dst == "" {
		_, err = Upx("put", src)
	} else {
		_, err = Upx("put", src, dst)
	}
	Nil(t, err)

	b, err := Upx("ls", correct)
	Nil(t, err)
	Equal(t, strings.HasSuffix(string(b), " "+correct+"\n"), true)
}

func compare(t *testing.T, local, up string) {
	locals := []string{}
	ups := []string{}
	fInfos, _ := ioutil.ReadDir(local)
	for _, fInfo := range fInfos {
		locals = append(locals, fInfo.Name())
	}

	b, err := Upx("ls", up)
	Nil(t, err)
	output := strings.TrimRight(string(b), "\n")
	for _, line := range strings.Split(output, "\n") {
		items := strings.Split(line, " ")
		ups = append(ups, items[len(items)-1])
	}

	sort.Strings(locals)
	sort.Strings(ups)

	Equal(t, len(locals), len(ups))
	for i := 0; i < len(locals); i++ {
		Equal(t, locals[i], ups[i])
	}
}

func putDir(t *testing.T, src, dst, correct string) {
	var err error
	if dst == "" {
		_, err = Upx("put", src)
	} else {
		_, err = Upx("put", src, dst)
	}
	Nil(t, err)

	compare(t, src, correct)
}

func getFile(t *testing.T, src, dst, correct string) {
	var err error
	if dst == "" {
		_, err = Upx("get", src)
	} else {
		_, err = Upx("get", src, dst)
	}
	Nil(t, err)

	_, err = os.Stat(correct)
	Nil(t, err)
}

func getDir(t *testing.T, src, dst, correct string) {
	var err error
	if dst == "" {
		_, err = Upx("get", src)
	} else {
		_, err = Upx("get", src, dst)
	}
	Nil(t, err)

	compare(t, correct, src)
}

func TestPutAndGet(t *testing.T) {
	base := ROOT + "/put/"
	pwd := "/tmp"
	localBase := filepath.Join(pwd, "put")
	func() {
		SetUp()
		err := os.MkdirAll(localBase, 0755)
		Nil(t, err)
	}()
	defer TearDown()

	err := os.Chdir(localBase)
	Nil(t, err)
	Upx("mkdir", base)
	Upx("cd", base)

	// upx put localBase/FILE upBase/FILE
	CreateFile("FILE")
	putFile(t, filepath.Join(localBase, "FILE"), "", path.Join(base, "FILE"))

	// upx put ../put/FILE2
	CreateFile("FILE2")
	localPath := ".." + string(filepath.Separator) + filepath.Join("put", "FILE2")
	putFile(t, localPath, "", path.Join(base, "FILE2"))

	// upx put /path/to/file /path/to/file
	putFile(t, "FILE", path.Join(base, "FILE4"), path.Join(base, "FILE4"))

	// upx put /path/to/file /path/to/dir
	CreateFile("FILE3")
	putFile(t, "FILE3", base, path.Join(base, "FILE3"))

	// upx put /path/to/file ../path/to/dir/
	putFile(t, "FILE", base+"/putfile/", path.Join(base, "putfile", "FILE"))

	// upx put ../path/to/dir
	localPath = ".." + string(filepath.Separator) + "put"
	putDir(t, localPath, "", path.Join(base, "put"))

	// upx put /path/to/dir /path/to/dir/
	putDir(t, localBase, base+"/putdir/", base+"/putdir/")

	_, err = Upx("put", localBase, path.Join(base, "FILE"))
	NotNil(t, err)

	localBase = filepath.Join(pwd, "get")
	os.MkdirAll(localBase, 0755)
	err = os.Chdir(localBase)
	Nil(t, err)

	// upx get /path/to/file
	getFile(t, path.Join(base, "FILE"), "", filepath.Join(localBase, "FILE"))

	// upx get ../path/to/file
	getFile(t, "../put/FILE2", "", filepath.Join(localBase, "FILE2"))

	// upx get /path/to/file /path/to/file
	getFile(t, "FILE4", filepath.Join(localBase, "FILE5"), filepath.Join(localBase, "FILE5"))

	// upx get /path/to/file /path/to/dir
	getFile(t, "FILE3", localBase, filepath.Join(localBase, "FILE3"))

	// upx get /path/to/file /path/to/dir/
	localPath = filepath.Join(localBase, "getfile") + string(filepath.Separator)
	getFile(t, "FILE", localPath, filepath.Join(localPath, "FILE"))

	// upx get ../path/to/dir
	getDir(t, "../put", "", filepath.Join(localBase, "put"))

	// upx get /path/to/dir /path/to/dir/
	localPath = filepath.Join(localBase, "getdir") + string(filepath.Separator)
	getDir(t, "../put", localPath, localPath)

	_, err = Upx("get", base, filepath.Join(localBase, "FILE"))
	NotNil(t, err)

	// upx get FILE*
	localPath = filepath.Join(localBase, "wildcard") + string(filepath.Separator)
	_, err = Upx("get", "FILE*", localPath)
	Nil(t, err)
	files, _ := Upx("ls", "FILE*")
	lfiles, _ := ioutil.ReadDir(localPath)
	NotEqual(t, len(lfiles), 0)
	Equal(t, len(lfiles)+1, len(strings.Split(string(files), "\n")))
}

func TestRm(t *testing.T) {
	SetUp()
	defer TearDown()
	base := ROOT + "/put/"
	Upx("cd", base)
	_, err := Upx("rm", "put")
	NotNil(t, err)

	_, err = Upx("rm", "put/FILE")
	Nil(t, err)
	_, err = Upx("ls", "put/FILE")
	NotNil(t, err)

	_, err = Upx("rm", "put/FILE*")
	Nil(t, err)
	_, err = Upx("ls", "put/FILE*")
	NotNil(t, err)

	_, err = Upx("rm", "-d", "put/*")
	Nil(t, err)
	_, err = Upx("ls", "-d", "put/*")
	NotNil(t, err)

	_, err = Upx("rm", "-a", "put")
	Nil(t, err)
	_, err = Upx("ls", "put")
	NotNil(t, err)
}
