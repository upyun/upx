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

func TestMove(t *testing.T) {
	SetUp()
	defer TearDown()

	upRootPath := path.Join(ROOT, "move")
	Upx("mkdir", upRootPath)

	localRootPath, err := ioutil.TempDir("", "test")
	assert.NoError(t, err)
	localRootName := filepath.Base(localRootPath)

	CreateFile(path.Join(localRootPath, "FILE1"))
	CreateFile(path.Join(localRootPath, "FILE2"))

	// 上传文件
	Upx("put", localRootPath, upRootPath)
	files, err := Ls(path.Join(upRootPath, localRootName))

	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE1", "FILE2"},
	)

	time.Sleep(time.Second)

	// 正常移动文件
	_, err = Upx(
		"mv",
		path.Join(upRootPath, localRootName, "FILE1"),
		path.Join(upRootPath, localRootName, "FILE3"),
	)
	assert.NoError(t, err)

	files, err = Ls(path.Join(upRootPath, localRootName))
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE2", "FILE3"},
	)

	time.Sleep(time.Second)

	// 目标文件已存在
	_, err = Upx(
		"mv",
		path.Join(upRootPath, localRootName, "FILE2"),
		path.Join(upRootPath, localRootName, "FILE3"),
	)
	assert.Equal(
		t,
		err.Error(),
		fmt.Sprintf(
			"target path %s already exists use -f to force overwrite\n",
			path.Join(upRootPath, localRootName, "FILE3"),
		),
	)

	files, err = Ls(path.Join(upRootPath, localRootName))
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE2", "FILE3"},
	)

	time.Sleep(time.Second)

	// 目标文件已存在, 强制覆盖
	_, err = Upx(
		"mv",
		"-f",
		path.Join(upRootPath, localRootName, "FILE2"),
		path.Join(upRootPath, localRootName, "FILE3"),
	)
	assert.NoError(t, err)

	files, err = Ls(path.Join(upRootPath, localRootName))
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE3"},
	)
}
