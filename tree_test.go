package upx

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	base := ROOT + "/ls"
	dirs, files := []string{}, []string{}

	func() {
		SetUp()
		Upx("mkdir", base)
		Upx("cd", base)

		for i := 0; i < 11; i++ {
			Upx("mkdir", fmt.Sprintf("dir%d", i))
			dirs = append(dirs, fmt.Sprintf("dir%d", i))
		}

		CreateFile("FILE")
		for i := 0; i < 5; i++ {
			Upx("put", "FILE", fmt.Sprintf("FILE%d", i))
			files = append(files, fmt.Sprintf("FILE%d", i))
		}
	}()

	defer func() {
		TearDown()
	}()

	tree1, err := Upx("tree")
	assert.NoError(t, err)
	tree1s := string(tree1)
	arr := strings.Split(tree1s, "\n")
	assert.Equal(t, len(arr), len(dirs)+len(files)+4)
	pwd, _ := Upx("pwd")
	assert.Equal(t, arr[0]+"\n", string(pwd))
	assert.Equal(t, arr[len(arr)-3], "")
	assert.Equal(t, arr[len(arr)-2], fmt.Sprintf("%d directories, %d files", len(dirs), len(files)))

	tree2, err := Upx("tree", base)
	assert.NoError(t, err)
	assert.Equal(t, string(tree2), string(tree1))
}
