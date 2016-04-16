package main

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jehiah/go-strftime"
	"github.com/upyun/go-sdk/upyun"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type Cmd struct {
	Desc  string
	Alias string
	Func  func(args []string)
}

var (
	conf   *Config
	user   *userInfo
	driver *FsDriver

	// TODO: refine
	confname = os.Getenv("HOME") + "/.upx.cfg"
)

var CmdMap = map[string]Cmd{
	"login":    {"Log in UPYUN with service_name, username, password", "", Login},
	"logout":   {"Log out UPYUN", "", Logout},
	"cd":       {"Change working directory", "", Cd},
	"pwd":      {"Print working directory", "", Pwd},
	"mkdir":    {"Make directory", "mk", Mkdir},
	"ls":       {"List directory or file", "", Ls},
	"switch":   {"Switch service", "sw", SwitchSrv},
	"services": {"List all services", "sv", ListSrvs},
	"put":      {"Put directory or file to UPYUN", "", Put},
	"get":      {"Get directory or file from UPYUN", "", Get},
	"rm":       {"Remove one or more directories and files", "", Rm},
	"version":  {"Print version", "", nil},    // deprecated
	"help":     {"Help information", "", nil}, // deprecated
	"info":     {"Current information", "i", Info},
}

type ByName []*upyun.FileInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func Login(args []string) {
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
		b, err := gopass.GetPasswdMasked()
		if err == nil {
			user.Password = string(b)
		}
	}

	if _, err := NewFsDriver(user.Bucket, user.Username,
		user.Password, user.CurDir, 10, nil); err != nil {
		fmt.Fprintf(os.Stderr, "failed to log in. %v\n", err)
		os.Exit(-1)
	}

	// save
	conf.UpdateUserInfo(user)
	conf.Save(confname)
}

func Logout(args []string) {
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
	// save
	conf.Save(confname)
}

func SwitchSrv(args []string) {
	if len(args) > 0 {
		bucket := args[0]
		if err := conf.SwitchBucket(bucket); err != nil {
			fmt.Println("switch:", err)
		}
		// save
		conf.Save(confname)
	}
}

func ListSrvs(args []string) {
	for k, v := range conf.Users {
		if k == conf.Idx {
			fmt.Printf("* \033[33m%s\033[0m\n", v.Bucket)
		} else {
			fmt.Printf("  %s\n", v.Bucket)
		}
	}
}

func Cd(args []string) {
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

func Ls(args []string) {
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

func Pwd(args []string) {
	fmt.Println(driver.GetCurDir())
}

func Get(args []string) {
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

func Put(args []string) {
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

func Rm(args []string) {
	for _, path := range args {
		rPath, wildcard := StrSplit(path)
		if wildcard == "" {
			driver.Remove(rPath)
		} else {
			driver.RemoveMatched(rPath, &MatchConfig{wildcard: wildcard})
		}
	}
}

func Mkdir(args []string) {
	for _, path := range args {
		if err := driver.MakeDir(path); err != nil {
			fmt.Fprintf(os.Stderr, "mkdir %s: %v\n\n", path, err)
		}
	}
}

func Info(args []string) {
	output := "ServiceName: " + user.Bucket + "\n"
	output += "Operator:    " + user.Username + "\n"
	output += "CurrentDir:  " + user.CurDir + "\n"
	fmt.Println(output)
}

func init() {
	conf = &Config{}
	conf.Load(confname)

	user = conf.GetCurUser()
	logger := log.New(os.Stdout, "", 0)
	if user != nil {
		var err error
		driver, err = NewFsDriver(user.Bucket, user.Username,
			user.Password, user.CurDir, 10, logger)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to log in. %v\n", err)
			os.Exit(-1)
		}
	}
}
