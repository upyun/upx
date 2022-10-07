package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestResumeput(t *testing.T) {
	base := ROOT + "/put/"
	pwd, err := ioutil.TempDir("", "test")
	Nil(t, err)
	localBase := filepath.Join(pwd, "put")
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

	_, err = Upx("resume-put", localBase, path.Join(base, "FILE"))
	NotNil(t, err)

}
