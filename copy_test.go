package upx

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	SetUp()
	defer TearDown()

	upRootPath := path.Join(ROOT, "copy")
	Upx("mkdir", upRootPath)

	localRootPath, err := ioutil.TempDir("", "test")
	assert.NoError(t, err)
	localRootName := filepath.Base(localRootPath)

	CreateFile(path.Join(localRootPath, "FILE1"))
	CreateFile(path.Join(localRootPath, "FILE2"))

	// 上传文件
	_, err = Upx("put", localRootPath, upRootPath)
	assert.NoError(t, err)

	files, err := Ls(path.Join(upRootPath, localRootName))
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE1", "FILE2"},
	)

	time.Sleep(time.Second)

	// 正常复制文件
	_, err = Upx(
		"cp",
		path.Join(upRootPath, localRootName, "FILE1"),
		path.Join(upRootPath, localRootName, "FILE3"),
	)
	assert.NoError(t, err)

	files, err = Ls(path.Join(upRootPath, localRootName))
	assert.NoError(t, err)
	assert.Len(t, files, 3)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE1", "FILE2", "FILE3"},
	)

	time.Sleep(time.Second)

	// 目标文件已存在
	_, err = Upx(
		"cp",
		path.Join(upRootPath, localRootName, "FILE1"),
		path.Join(upRootPath, localRootName, "FILE2"),
	)
	assert.Error(t, err)
	assert.Equal(
		t,
		err.Error(),
		fmt.Sprintf(
			"target path %s already exists use -f to force overwrite\n",
			path.Join(upRootPath, localRootName, "FILE2"),
		),
	)

	files, err = Ls(path.Join(upRootPath, localRootName))
	assert.NoError(t, err)
	assert.Len(t, files, 3)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE1", "FILE2", "FILE3"},
	)

	time.Sleep(time.Second)

	// 目标文件已存在, 强制覆盖
	_, err = Upx(
		"cp",
		"-f",
		path.Join(upRootPath, localRootName, "FILE1"),
		path.Join(upRootPath, localRootName, "FILE2"),
	)
	assert.NoError(t, err)

	files, err = Ls(path.Join(upRootPath, localRootName))
	assert.NoError(t, err)
	assert.Len(t, files, 3)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE1", "FILE2", "FILE3"},
	)
}
