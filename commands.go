package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jehiah/go-strftime"
	"github.com/upyun/go-sdk/upyun"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Cmd struct {
	Desc  string
	Alias string
	Func  func(args []string, opts map[string]interface{})
	Flags map[string]CmdFlag
}

type CmdFlag struct {
	usage string
	typ   string
}

var (
	conf   *Config
	user   *userInfo
	driver *FsDriver

	// TODO: refine
	confname = filepath.Join(os.Getenv("HOME"), ".upx.cfg")
	dbname   = filepath.Join(os.Getenv("HOME"), ".upx.db")
)

var (
	GlobalFlags = map[string]CmdFlag{
		"auth": CmdFlag{"auth information", "string"},
		"v":    CmdFlag{"verbose", "bool"},
	}
	RmFlags = map[string]CmdFlag{
		"d":     CmdFlag{"only remove directories", "bool"},
		"a":     CmdFlag{"remove files, directories and their contents recursively, never prompt", "bool"},
		"async": CmdFlag{"remove files async", "bool"},
	}
	LsFlags = map[string]CmdFlag{
		"r": CmdFlag{"reverse order", "bool"},
		"c": CmdFlag{"max items to list", "int"},
		"d": CmdFlag{"only show directory", "bool"},
	}
	SyncFlags = map[string]CmdFlag{
		"w": CmdFlag{"worker number", "int"},
	}
)

var CmdMap = map[string]Cmd{
	"login":    {"Log in UPYUN with service_name, username, password", "", Login, nil},
	"logout":   {"Log out UPYUN", "", Logout, nil},
	"cd":       {"Change working directory", "", Cd, nil},
	"pwd":      {"Print working directory", "", Pwd, nil},
	"mkdir":    {"Make directory", "mk", Mkdir, nil},
	"ls":       {"List directory or file", "", Ls, LsFlags},
	"switch":   {"Switch service, alias sw", "sw", SwitchSrv, nil},
	"services": {"List all services, alias sv", "sv", ListSrvs, nil},
	"sync":     {"sync folder to UPYUN", "", Sync, SyncFlags},
	"put":      {"Put directory or file to UPYUN", "", Put, nil},
	"get":      {"Get directory or file from UPYUN", "", Get, nil},
	"rm":       {"Remove one or more directories and files", "", Rm, RmFlags},
	"version":  {"Print version", "", nil, nil},    // deprecated
	"help":     {"Help information", "", nil, nil}, // deprecated
	"info":     {"Current information", "i", Info, nil},
	"auth":     {"generate auth string", "", GenAuth, nil},
}

type ByName []*upyun.FileInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func Login(args []string, opts map[string]interface{}) {
	user := &userInfo{CurDir: "/"}
	if len(args) == 3 {
		user.Bucket = args[0]
		user.Username = args[1]
		user.Password = args[2]
	} else {
		fmt.Printf("BucketName: ")
		fmt.Scanf("%s\n", &user.Bucket)
		fmt.Printf("Operator: ")
		fmt.Scanf("%s\n", &user.Username)
		fmt.Printf("Password: ")
		b, err := gopass.GetPasswdMasked()
		if err == nil {
			user.Password = string(b)
		}
	}

	LogI("\n")

	if _, err := NewFsDriver(user.Bucket, user.Username,
		user.Password, user.CurDir, 10, nil); err != nil {
		LogC("login: %v", err)
		os.Exit(-1)
	}

	// save
	conf.UpdateUserInfo(user)
	conf.Save(confname)

	LogI("Welcome to %s, %s!", user.Bucket, user.Username)
}

func Logout(args []string, opts map[string]interface{}) {
	var err error
	if err = conf.RemoveBucket(); err == nil {
		if len(conf.Users) == 0 {
			err = os.Remove(confname)
		}
	}
	if err != nil {
		LogC("logout: %v", err)
	}
	// save
	conf.Save(confname)
}

func SwitchSrv(args []string, opts map[string]interface{}) {
	if len(args) > 0 {
		bucket := args[0]
		if err := conf.SwitchBucket(bucket); err != nil {
			LogE("switch: %v", err)
		}
		// save
		conf.Save(confname)
	}
}

func ListSrvs(args []string, opts map[string]interface{}) {
	for k, v := range conf.Users {
		if k == conf.Idx {
			LogI("* \033[33m%s\033[0m\n", v.Bucket)
		} else {
			LogI("  %s\n", v.Bucket)
		}
	}
}

func Cd(args []string, opts map[string]interface{}) {
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
		LogC("cd %s: %v", path, err)
	}
}

func Ls(args []string, opts map[string]interface{}) {
	fpath := driver.GetCurDir()
	if len(args) > 0 {
		fpath = args[0]
	}

	if !driver.IsUPDir(fpath) {
		info, err := driver.up.GetInfo(fpath)
		if err != nil {
			LogE("getinfo: %v", err)
			return
		}
		info.Name = path.Base(fpath)
		LogI(parseInfo(info))
		return
	}
	maxCount, cnt, asc, onlyDir := 0, 0, true, false

	if v, ok := opts["c"]; ok {
		maxCount = v.(int)
	}
	if v, ok := opts["r"]; ok {
		if v.(bool) {
			asc = false
		}
	}
	if v, ok := opts["d"]; ok {
		if v.(bool) {
			onlyDir = true
		}
	}

	infos, errChannel := driver.up.GetLargeList(fpath, asc, false)
	for {
		select {
		case info, more := <-infos:
			if !more {
				return
			}
			if maxCount > 0 && cnt >= maxCount {
				return
			}
			if onlyDir && info.Type != "folder" {
				continue
			}
			LogI("%s\n", parseInfo(info))
			cnt++
		case err := <-errChannel:
			if err != nil {
				LogE("ls %s: %v", fpath, err)
				return
			}
		}
	}
}

func Pwd(args []string, opts map[string]interface{}) {
	LogI(driver.GetCurDir())
}

func Get(args []string, opts map[string]interface{}) {
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
		LogC("get %s %s: %v\n\n", src, des, err)
	}

	time.Sleep(time.Second)
}

func Sync(args []string, opts map[string]interface{}) {
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

	if v, ok := opts["w"]; ok {
		maxWorker = v.(int)
	}

	doSync(src, des)
}

func Put(args []string, opts map[string]interface{}) {
	var src, des string
	switch len(args) {
	case 1:
		src, des = args[0], "./"
	case 2:
		src, des = args[0], args[1]
	default:
	}

	if err := driver.Uploads(src, des); err != nil {
		LogC("put %s %s: %v\n\n", src, des, err)
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

func Rm(args []string, opts map[string]interface{}) {
	for _, path := range args {
		rPath, wildcard := StrSplit(path)
		match := &MatchConfig{
			wildcard: wildcard,
			itemType: "file",
		}
		if v, exists := opts["d"]; exists && v.(bool) {
			match.itemType = "folder"
		}
		if v, exists := opts["a"]; exists && v.(bool) {
			match.itemType = ""
		}
		async := false
		if v, exists := opts["async"]; exists {
			async = v.(bool)
		}
		driver.RemoveMatched(rPath, match, async)
	}
}

func Mkdir(args []string, opts map[string]interface{}) {
	for _, path := range args {
		if err := driver.MakeDir(path); err != nil {
			LogE("mkdir %s: %v\n\n", path, err)
		}
	}
}

func Info(args []string, opts map[string]interface{}) {
	usage, _ := driver.up.Usage()
	output := fmt.Sprintf("BucketName: %s\n", user.Bucket)
	output += fmt.Sprintf("Operator:    %s\n", user.Username)
	output += fmt.Sprintf("CurrentDir:  %s\n", user.CurDir)
	output += fmt.Sprintf("Usage:       %.3fMB\n", float64(usage)/1024/1024)
	LogI(output)
}

func GenAuth(args []string, opts map[string]interface{}) {
	if len(args) != 3 {
		LogC("not enough arguments. bucket username password")
	}
	bucket, username, passwd := args[0], args[1], args[2]
	LogI(genAuth(bucket, username, passwd))
}

func parseInfo(info *upyun.FileInfo) string {
	s := "drwxrwxrwx"
	if info.Type != "folder" {
		s = "-rw-rw-rw-"
	}
	s += fmt.Sprintf(" 1 %s %s %12d", user.Username, user.Bucket, info.Size)
	if info.Time.Year() != time.Now().Year() {
		s += " " + strftime.Format("%b %d  %Y", info.Time)
	} else {
		s += " " + strftime.Format("%b %d %H:%M", info.Time)
	}
	s += " " + info.Name
	return s
}

func initDriver(auth string) {
	if runtime.GOOS == "windows" {
		confname = filepath.Join(os.Getenv("USERPROFILE"), ".upx.cfg")
		dbname = filepath.Join(os.Getenv("USERPROFILE"), ".upx.db")
	}

	logger := log.New(os.Stdout, "", 0)

	if auth == "" {
		conf = &Config{}
		conf.Load(confname)

		user = conf.GetCurUser()
		if user != nil {
			var err error
			driver, err = NewFsDriver(user.Bucket, user.Username,
				user.Password, user.CurDir, 10, logger)
			if err != nil {
				conf.RemoveBucket()
				conf.Save(confname)
				LogC("failed to log in. %v\n", err)
			}
		}
	} else {
		var err error
		var b []byte
		if b, err = base64.StdEncoding.DecodeString(auth); err == nil {
			var u userInfo
			if err = json.Unmarshal(b, &u); err == nil {
				user = &u
				driver, err = NewFsDriver(user.Bucket, user.Username,
					user.Password, "/", 10, logger)
			}
		}
		if err != nil {
			LogC("failed to log in. %v\n", err)
		}
	}
}
