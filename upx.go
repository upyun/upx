// +build linux darwin

package main

import (
	. "./fsdriver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jehiah/go-strftime"
	"github.com/polym/go-sdk/upyun"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sort"
)

const (
	version = "v0.0.1"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Bucket   string `json:"bucket"`
	CurDir   string `sjon:"curdir"`
}

var (
	conf             *Config
	driver           *FsDriver
	username, bucket string

	// TODO: refine
	confname = os.Getenv("HOME") + "/.upx.cfg"

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

func loadConfig() (*Config, error) {
	var err error
	var fd *os.File
	var b []byte
	var config Config

	if fd, err = os.Open(confname); err == nil {
		defer fd.Close()
		if b, err = ioutil.ReadAll(fd); err == nil {
			if b, err = base64.StdEncoding.DecodeString(string(b)); err == nil {
				err = json.Unmarshal(b, &config)
			}
		}
	}

	if err != nil {
		return nil, err
	}
	return &config, nil
}

func saveConfig(conf *Config) error {
	var err error
	var fd *os.File
	var b []byte
	var s string

	if fd, err = os.OpenFile(confname, os.O_RDWR|os.O_TRUNC|os.O_CREATE,
		0600); err == nil {
		defer fd.Close()
		if b, err = json.Marshal(conf); err == nil {
			s = base64.StdEncoding.EncodeToString(b)
			_, err = fd.WriteString(s)
		}
	}

	return err
}

func NewHandler() (driver *FsDriver, err error) {
	if conf, err = loadConfig(); err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "", 0)
	return NewFsDriver(conf.Bucket, conf.Username, conf.Password, conf.CurDir, 10, logger)
}

func Login(args ...string) {
	config := &Config{CurDir: "/"}

	fmt.Printf("ServiceName: ")
	fmt.Scanf("%s\n", &config.Bucket)
	fmt.Printf("Username: ")
	fmt.Scanf("%s\n", &config.Username)
	fmt.Printf("Password: ")
	config.Password = string(gopass.GetPasswdMasked())

	if err := saveConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "login: %v\n\n", err)
		os.Exit(-1)
	}
}

func Logout() {
	if err := os.Remove(confname); err != nil {
		fmt.Fprintf(os.Stderr, "logout: %v\n\n", err)
		os.Exit(-1)
	}
}

func Cd(args ...string) {
	path := "/"
	if len(args) > 0 {
		path = args[0]
	}

	var err error
	if err = driver.ChangeDir(path); err == nil {
		conf.CurDir = driver.GetCurDir()
		err = saveConfig(conf)
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
			s += fmt.Sprintf(" 1 %s %s %12d", username, bucket, v.Size)
			s += " " + strftime.Format("%b %d %H:%M", v.Time)
			s += " " + v.Name
			fmt.Println(s)
		}
	}
}

func Pwd() {
	fmt.Println(driver.GetCurDir() + "\n")
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
		fmt.Fprintf(os.Stderr, "get %s %s: %v\n\n", src, des, err)
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
		fmt.Fprintf(os.Stderr, "put %s %s: %v\n\n", src, des, err)
		os.Exit(-1)
	}
}

func Rm(args ...string) {
	for _, path := range args {
		if ok, err := driver.IsDir(path); err == nil && ok {
			fmt.Printf("< %s > is a directory. Are you sure to remove it? (y/n) ", path)
			var ans string
			if fmt.Scanf("%s", &ans); ans != "y" {
				continue
			}
		}
		if err := driver.Remove(path); err != nil {
			fmt.Fprintf(os.Stderr, "remove %s: %v\n\n", path, err)
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

func main() {
	args := os.Args
	if len(args) == 1 {
		Help(args...)
		os.Exit(-1)
	}

	switch args[1] {
	case "login":
		Login()
	case "logout":
		Logout()
		return
	case "help":
		Help(args...)
		return
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
