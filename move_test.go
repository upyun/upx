package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func mvFile(t *testing.T, oldPath string, newPath string) {
	t.Log(oldPath, newPath)
	_, err := Upx("mv", oldPath, newPath)
	Nil(t, err)

}

func TestMv(t *testing.T) {
	base := ROOT + "/mv/"
	pwd, err := ioutil.TempDir("", "test")
	Nil(t, err)
	localBase := filepath.Join(pwd, "mv")
	func() {
		SetUp()
		err := os.MkdirAll(localBase, 0755)
		err = os.MkdirAll(localBase+"/test", 0755)
		Nil(t, err)
	}()

	defer func() {
		TearDown()
	}()

	err = os.Chdir(localBase)
	//Nil(t, err)
	Upx("mkdir", base)
	Upx("mv", base)
	// upx put localBase/FILE upBase/FILE
	getwd, err := os.Getwd()
	if err != nil {
		return
	}

	t.Log("local:", getwd)
	t.Log("localbase:", localBase)
	putDir(t, localBase, base+"/putdir/", base+"/putdir/")
	CreateFile("FILE")
	oldPath := filepath.Join(base, "FILE")
	putFile(t, filepath.Join(localBase, "FILE"), "", path.Join(base, "FILE"))
	newPath := base + "putdir/"
	t.Log("dir", localBase+"test", base)
	t.Log(oldPath, newPath)
	mvFile(t, oldPath, newPath)
}
