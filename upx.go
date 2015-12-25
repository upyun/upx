// +build linux darwin

package main

import (
	. "./fsdriver"
	"fmt"
	"github.com/jehiah/go-strftime"
	"github.com/polym/go-sdk/upyun"
	"log"
	"os"
	"runtime"
	"sort"
	"syscall"
)

const version = "v0.0.1"

var (
	driver           *FsDriver
	username, bucket string

	cmdDesc = map[string]string{
		"cd":      "Change working directory",
		"pwd":     "Print working directory",
		"mkdir":   "Make directory",
		"ls":      "List directory or file",
		"login":   "Log in UPYUN with username, password, bucket",
		"logout":  "Log out UPYUN",
		"put":     "Put directory or file to UPYUN",
		"get":     "Get directory or file from UPYUN",
		"rm":      "Remove one or more directories and files",
		"version": "Print version",
	}
)

type ByName []*upyun.FileInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func NewHandler() (*FsDriver, error) {
	username = os.Getenv("UPX_USERNAME")
	password := os.Getenv("UPX_PASSWORD")
	bucket = os.Getenv("UPX_BUCKET")
	curDir := os.Getenv("UPX_CURDIR")
	if curDir == "" {
		curDir = "/"
	}

	logger := log.New(os.Stdout, "", 0)
	return NewFsDriver(bucket, username, password, curDir, 10, logger)
}

func SetEnvALL(envs map[string]string) error {
	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, syscall.Environ())
}

func Login(args ...string) {
	if len(args) != 3 {
		fmt.Fprintf(os.Stderr, "login operator password bucket\n")
		os.Exit(-1)
	}
	smap := map[string]string{
		"UPX_USERNAME": args[0],
		"UPX_PASSWORD": args[1],
		"UPX_BUCKET":   args[2],
		"UPX_CURDIR":   "/",
	}

	SetEnvALL(smap)
}

func Logout() {
	smap := map[string]string{
		"UPX_USERNAME": "",
		"UPX_PASSWORD": "",
		"UPX_BUCKET":   "",
		"UPX_CURDIR":   "/",
	}

	SetEnvALL(smap)
}

func Cd(args ...string) {
	path := "/"
	if len(args) > 0 {
		path = args[0]
	}

	if err := driver.ChangeDir(path); err != nil {
		fmt.Fprintf(os.Stderr, "cd %s: %v\n", path, err)
		os.Exit(-1)
	}

	smap := map[string]string{
		"UPX_CURDIR": driver.GetCurDir(),
	}

	SetEnvALL(smap)
}

func Ls(args ...string) {
	path := driver.GetCurDir()
	if len(args) > 0 {
		path = args[0]
	}
	if infos, err := driver.ListDir(path); err != nil {
		fmt.Fprintf(os.Stderr, "ls %s: %v\n", path, err)
		os.Exit(-1)
	} else {
		sort.Sort(ByName(infos))
		for _, v := range infos {
			s := "drwxrwxrwx"
			if v.Type != "folder" {
				s = "-rw-rw-rw-"
			}
			s += fmt.Sprintf(" 1 %s %s %12d", username, bucket, v.Size)
			s += " " + strftime.Format("%b %d %H:%M", v.Time)
			s += " " + v.Name
			fmt.Println(s)
		}
	}
}

func Pwd() {
	fmt.Println(driver.GetCurDir())
}

func Get(args ...string) {
	var src, des string
	switch len(args) {
	case 0:
		// TODO
	case 1:
		src = args[0]
		des = "./"
	case 2:
		src = args[0]
		des = args[1]
	}

	if err := driver.GetItems(src, des); err != nil {
		fmt.Fprintf(os.Stderr, "get %s %s: %v\n", src, des, err)
		os.Exit(-1)
	}
}

func Put(args ...string) {
	var src, des string
	switch len(args) {
	case 0:
		// TODO
	case 1:
		src = args[0]
		des = ""
	case 2:
		src = args[0]
		des = args[1]
	}

	if err := driver.PutItems(src, des); err != nil {
		fmt.Fprintf(os.Stderr, "put %s %s: %v\n", src, des, err)
		os.Exit(-1)
	}
}

func Rm(args ...string) {
	for _, path := range args {
		if err := driver.Remove(path); err != nil {
			fmt.Fprintf(os.Stderr, "remove %s: %v\n", path, err)
		}
	}
}

func Mkdir(args ...string) {
	for _, path := range args {
		fmt.Println(path)
		if err := driver.MakeDir(path); err != nil {
			fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", path, err)
		}
	}
}

func Help(args ...string) {
	cmd := args[0]
	s := cmd + " is a tool for managing files in UPYUN\n\n"
	s = s + "Usage:\n\n"
	s = s + "\t" + cmd + " command [arguments]\n\n"
	s = s + "The commands are:\n\n"

	var desc []string
	for k, _ := range cmdDesc {
		desc = append(desc, k)
	}
	sort.Strings(desc)
	for _, k := range desc {
		s += fmt.Sprintf("\t%-8s %s\n", k, cmdDesc[k])
	}
	s += "\n"
	fmt.Println(s)
}

func main() {
	args := os.Args
	if len(args) == 1 {
		Help(args...)
		os.Exit(-1)
	}

	switch args[1] {
	case "login":
		Login(args[2:]...)
	case "logout":
		Logout()
		return
	case "help":
		Help(args...)
		return
	case "version":
		fmt.Printf("%s version %s %s/%s\n", args[0], version, runtime.GOOS, runtime.GOARCH)
		return
	}

	var err error
	driver, err = NewHandler()
	if err != nil {
		fmt.Fprintf(os.Stderr, "re login %v\n", err)
		os.Exit(-1)
	}

	switch args[1] {
	case "cd":
		Cd(args[2:]...)
	case "ls":
		Ls(args[2:]...)
	case "get":
		Get(args[2:]...)
	case "put":
		Put(args[2:]...)
	case "rm":
		Rm(args[2:]...)
	case "pwd":
		Pwd()
	case "mkdir":
		Mkdir(args[2:]...)
	}
}
