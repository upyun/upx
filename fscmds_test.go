package main

import (
	"fmt"
	"path"
	"strings"
	"testing"
)

/*
	mkdir /path/to/mkdir/case1
	cd /path/to/mkdir
	mkdir case2
	cd case2
	mkdir ../case3
	cd ../case3
	ls /path/to/mkdir
*/
func TestMkdirAndCdAndPwd(t *testing.T) {
	SetUp()
	defer TearDown()

	base := path.Join(ROOT, "mkdir")

	case1 := path.Join(base, "case1")
	b, err := Upx("mkdir", case1)
	Nil(t, err)
	Equal(t, string(b), "")

	Upx("cd", base)
	b, _ = Upx("pwd")
	Equal(t, string(b), base+"\n")

	case2 := path.Join(base, "case2")
	b, err = Upx("mkdir", "case2")
	Nil(t, err)
	Equal(t, string(b), "")

	Upx("cd", "case2")
	b, _ = Upx("pwd")
	Equal(t, string(b), case2+"\n")

	case3 := path.Join(base, "case3")
	b, err = Upx("mkdir", "../case3")
	Nil(t, err)
	Equal(t, string(b), "")

	Upx("cd", "../case3")
	b, _ = Upx("pwd")
	Equal(t, string(b), case3+"\n")

	// check
	b, err = Upx("ls", base)
	Nil(t, err)
	output := string(b)
	lines := strings.Split(output, "\n")
	Equal(t, len(lines), 4)
	Equal(t, strings.Contains(output, " case1\n"), true)
	Equal(t, strings.Contains(output, " case2\n"), true)
	Equal(t, strings.Contains(output, " case3\n"), true)
}

/*
	ls /path/to/file
	ls -r /path/to/dir
	ls -c 10 /path/to/dir
	ls -d /path/to/dir
	ls -r -d -c 10 /path/to/dir
*/
func TestLs(t *testing.T) {
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
		for _, file := range files {
			Upx("rm", file)
		}
		for _, dir := range dirs {
			Upx("rm", dir)
		}
		TearDown()
	}()

	b, err := Upx("ls")
	Nil(t, err)
	Equal(t, len(strings.Split(string(b), "\n")), len(dirs)+len(files)+1)

	normal, err := Upx("ls", base)
	Nil(t, err)
	Equal(t, len(strings.Split(string(normal), "\n")), len(dirs)+len(files)+1)

	c := (len(dirs) + len(files)) - 1
	limited, err := Upx("ls", "-c", fmt.Sprint(c))
	Nil(t, err)
	Equal(t, len(strings.Split(string(limited), "\n")), c+1)

	folders, err := Upx("ls", "-d")
	Nil(t, err)
	Equal(t, len(strings.Split(string(folders), "\n")), len(dirs)+1)

	c = len(dirs) - 1
	lfolders, err := Upx("ls", "-d", "-c", fmt.Sprint(c))
	Nil(t, err)
	Equal(t, len(strings.Split(string(lfolders), "\n")), c+1)
	for _, line := range strings.Split(string(lfolders), "\n")[0:c] {
		Equal(t, strings.HasPrefix(line, "drwxrwxrwx "), true)
	}

	lfiles, err := Upx("ls", "FILE*")
	Nil(t, err)
	Equal(t, len(strings.Split(string(lfiles), "\n")), 6)

	reversed, err := Upx("ls", "-r", base)
	Nil(t, err)
	NotEqual(t, string(reversed), string(normal))
}
