// +build linux darwin

package main

import (
	"errors"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jehiah/go-strftime"
	"github.com/upyun/go-sdk/upyun"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

var (
	conf    *Config
	user    *userInfo
	driver  *FsDriver
	version string

	// TODO: refine
	confname = os.Getenv("HOME") + "/.upx.cfg"

	cmdDesc = map[string]string{
		"cd":       "Change working directory",
		"pwd":      "Print working directory",
		"mkdir":    "Make directory",
		"ls":       "List directory or file",
		"login":    "Log in UPYUN with username, password, bucket",
		"logout":   "Log out UPYUN",
		"switch":   "Switch service",
		"services": "List all services",
		"put":      "Put directory or file to UPYUN",
		"get":      "Get directory or file from UPYUN",
		"rm":       "Remove one or more directories and files",
		"version":  "Print version",
		"help":     "Help information",
		"info":     "Current information",
	}
)

type ByName []*upyun.FileInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func NewHandler() (driver *FsDriver, err error) {
	user = conf.GetCurUser()
	logger := log.New(os.Stdout, "", 0)
	if user == nil {
		return nil, errors.New("no user")
	}
	return NewFsDriver(user.Bucket, user.Username, user.Password, user.CurDir, 10, logger)
}

func SaveConfig() {
	if err := conf.Save(confname); err != nil {
		fmt.Fprintf(os.Stderr, "save config file:%v\n\n", err)
		os.Exit(-1)
	}
}

func Login(args ...string) {
	user := &userInfo{CurDir: "/"}
	if len(args) == 3 {
		user.Bucket = args[0]
		user.Username = args[1]
		user.Password = args[2]
	} else {
		fmt.Printf("ServiceName: ")
		fmt.Scanf("%s\n", &user.Bucket)
		fmt.Printf("Operator: ")
		fmt.Scanf("%s\n", &user.Username)
		fmt.Printf("Password: ")
		user.Password = string(gopass.GetPasswdMasked())
	}

	conf.UpdateUserInfo(user)
}

func Logout() {
	var err error
	if err = conf.RemoveBucket(); err == nil {
		if len(conf.Users) == 0 {
			err = os.Remove(confname)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "logout: %v\n\n", err)
		os.Exit(-1)
	}
}

func SwitchSrv(args ...string) {
	if len(args) > 0 {
		bucket := args[0]
		if err := conf.SwitchBucket(bucket); err != nil {
			fmt.Println("switch:", err)
		}
	}
}

func ListSrvs() {
	for k, v := range conf.Users {
		if k == conf.Idx {
			fmt.Printf("* \033[33m%s\033[0m\n", v.Bucket)
		} else {
			fmt.Printf("  %s\n", v.Bucket)
		}
	}
}

func Cd(args ...string) {
	path := "/"
	if len(args) > 0 {
		path = args[0]
	}

	var err error
	if err = driver.ChangeDir(path); err == nil {
		user.CurDir = driver.GetCurDir()
		err = conf.Save(confname)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "cd %s: %v\n\n", path, err)
		os.Exit(-1)
	}
}

func Ls(args ...string) {
	path := driver.GetCurDir()
	if len(args) > 0 {
		path = args[0]
	}
	if infos, err := driver.ListDir(path); err != nil {
		fmt.Fprintf(os.Stderr, "ls %s: %v\n\n", path, err)
		os.Exit(-1)
	} else {
		sort.Sort(ByName(infos))
		for _, v := range infos {
			s := "drwxrwxrwx"
			if v.Type != "folder" {
				s = "-rw-rw-rw-"
			}
			s += fmt.Sprintf(" 1 %s %s %12d", user.Username, user.Bucket, v.Size)
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

	if err := driver.Downloads(src, des); err != nil {
		fmt.Fprintf(os.Stderr, "get %s %s: %v\n\n", src, des, err)
		os.Exit(-1)
	}

	time.Sleep(time.Second)
}

func Put(args ...string) {
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

	if err := driver.Uploads(src, des); err != nil {
		fmt.Fprintf(os.Stderr, "put %s %s: %v\n\n", src, des, err)
		os.Exit(-1)
	}

	time.Sleep(time.Second)
}

func StrSplit(s string) (path, wildcard string) {
	idx := strings.Index(s, "*")
	if idx == -1 {
		return s, ""
	}
	idx = strings.LastIndex(s[:idx], "/")
	if idx == -1 {
		return "./", s
	}
	return s[:idx+1], s[idx+1:]
}

func Rm(args ...string) {
	for _, path := range args {
		rPath, wildcard := StrSplit(path)
		if wildcard == "" {
			driver.Remove(rPath)
		} else {
			driver.RemoveMatched(rPath, &MatchConfig{wildcard: wildcard})
		}
	}
}

func Mkdir(args ...string) {
	for _, path := range args {
		if err := driver.MakeDir(path); err != nil {
			fmt.Fprintf(os.Stderr, "mkdir %s: %v\n\n", path, err)
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

func Info() {
	output := "ServiceName: " + user.Bucket + "\n"
	output += "Operator:    " + user.Username + "\n"
	output += "CurrentDir:  " + user.CurDir + "\n"
	fmt.Println(output)
}

func main() {
	args := os.Args
	if len(args) == 1 {
		Help(args...)
		os.Exit(-1)
	}

	conf = &Config{}
	conf.Load(confname)
	defer SaveConfig()

	switch args[1] {
	case "login":
		Login(args[2:]...)
	case "logout":
		Logout()
		return
	case "help":
		Help(args...)
		return
	case "services":
		ListSrvs()
		return
	case "switch":
		SwitchSrv(args[2:]...)
	case "version":
		fmt.Printf("%s version %s %s/%s\n\n", args[0], version, runtime.GOOS, runtime.GOARCH)
		return
	}

	var err error
	driver, err = NewHandler()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nfailed to log in. %v\n\n", err)
		os.Exit(-1)
	}

	defer driver.progress.Stop()

	switch args[1] {
	case "login", "switch":
		return
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
	case "info":
		Info()
	default:
		Help(args...)
	}
}
