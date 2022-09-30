package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestCp(t *testing.T) {
	base := ROOT + "/cp/"
	pwd, err := ioutil.TempDir("", "test")
	Nil(t, err)
	localBase := filepath.Join(pwd, "cp")
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
	Upx("cp", base)
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
	t.Log(oldPath, newPath)
	_, err = Upx("cp", oldPath, newPath)
	Nil(t, err)
}
