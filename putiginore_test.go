package upx

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Ls(up string) ([]string, error) {
	b, err := Upx("ls", up)
	if err != nil {
		return nil, err
	}

	var ups = make([]string, 0)
	output := strings.TrimRight(string(b), "\n")
	for _, line := range strings.Split(output, "\n") {
		items := strings.Split(line, " ")
		ups = append(ups, items[len(items)-1])
	}
	return ups, nil
}

func TestPutIgnore(t *testing.T) {
	SetUp()
	defer TearDown()

	upRootPath := path.Join(ROOT, "iginore")
	Upx("mkdir", upRootPath)

	localRootPath, err := ioutil.TempDir("", "test")
	assert.NoError(t, err)
	localRootName := filepath.Base(localRootPath)

	CreateFile(path.Join(localRootPath, "FILE1"))
	CreateFile(path.Join(localRootPath, "FILE2"))
	CreateFile(path.Join(localRootPath, ".FILE3"))
	CreateFile(path.Join(localRootPath, ".FILES/FILE"))

	// 上传文件夹
	// 不包含隐藏的文件，所以只有FILE1和FILE2
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

	// 上传隐藏的文件夹, 无all，上传失效
	Upx(
		"put",
		path.Join(localRootPath, ".FILES"),
		path.Join(upRootPath, localRootName, ".FILES"),
	)
	files, err = Ls(path.Join(upRootPath, localRootName))
	assert.NoError(t, err)

	assert.Len(t, files, 2)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE1", "FILE2"},
	)

	time.Sleep(time.Second)

	// 上传隐藏的文件夹, 有all，上传成功
	Upx(
		"put",
		"-all",
		path.Join(localRootPath, ".FILES"),
		path.Join(upRootPath, localRootName, ".FILES"),
	)
	files, err = Ls(path.Join(upRootPath, localRootName))
	assert.NoError(t, err)
	assert.Len(t, files, 3)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE1", "FILE2", ".FILES"},
	)

	time.Sleep(time.Second)

	// 上传所有文件
	Upx("put", "-all", localRootPath, upRootPath)
	files, err = Ls(path.Join(upRootPath, localRootName))
	assert.NoError(t, err)
	assert.Len(t, files, 4)
	assert.ElementsMatch(
		t,
		files,
		[]string{"FILE1", "FILE2", ".FILE3", ".FILES"},
	)
}
