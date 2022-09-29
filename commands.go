package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

const (
	NO_CHECK = false
	CHECK    = true
)

func Init(login bool) {
	InitAndCheck(login, false, nil)
}

func InitAndCheck(login, check bool, c *cli.Context) {
	if login == LOGIN && session == nil {
		readConfigFromFile(LOGIN)
	}
	if login == NO_LOGIN {
		readConfigFromFile(NO_LOGIN)
	}
	if check && c.NArg() == 0 && c.NumFlags() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(-1)
	}
}

func NewLoginCommand() cli.Command {
	return cli.Command{
		Name:  "login",
		Usage: "Log in to UpYun",
		Action: func(c *cli.Context) error {
			Init(NO_LOGIN)
			session = &Session{CWD: "/"}
			args := c.Args()
			if len(args) == 3 {
				session.Bucket = args.Get(0)
				session.Operator = args.Get(1)
				session.Password = args.Get(2)
			} else {
				fmt.Printf("ServiceName: ")
				fmt.Scanf("%s\n", &session.Bucket)
				fmt.Printf("Operator: ")
				fmt.Scanf("%s\n", &session.Operator)
				fmt.Printf("Password: ")
				b, err := gopass.GetPasswdMasked()
				if err == nil {
					session.Password = string(b)
				}
				// TODO
				Print("")
			}

			if err := session.Init(); err != nil {
				PrintErrorAndExit("login failed: %v", err)
			}
			Print("Welcome to %s, %s!", session.Bucket, session.Operator)

			if config == nil {
				config = &Config{
					SessionId: 0,
					Sessions:  []*Session{session},
				}
			} else {
				config.Insert(session)
			}
			saveConfigToFile()

			return nil
		},
	}
}

func NewLogoutCommand() cli.Command {
	return cli.Command{
		Name:  "logout",
		Usage: "Log out of your UpYun account",
		Action: func(c *cli.Context) error {
			Init(NO_LOGIN)
			if session != nil {
				op, bucket := session.Operator, session.Bucket
				config.PopCurrent()
				saveConfigToFile()
				Print("Goodbye %s/%s ~~", op, bucket)
			} else {
				PrintErrorAndExit("nothing to do")
			}
			return nil
		},
	}
}

func NewAuthCommand() cli.Command {
	return cli.Command{
		Name:  "auth",
		Usage: "Generate auth string",
		Action: func(c *cli.Context) error {
			if c.NArg() == 3 {
				s, err := makeAuthStr(c.Args()[0], c.Args()[1], c.Args()[2])
				if err != nil {
					PrintErrorAndExit("auth: %v", err)
				}
				Print(s)
			} else {
				PrintErrorAndExit("auth: invalid parameters")
			}
			return nil
		},
	}
}

func NewListSessionsCommand() cli.Command {
	return cli.Command{
		Name:  "sessions",
		Usage: "List all sessions",
		Action: func(c *cli.Context) error {
			Init(NO_LOGIN)
			for k, v := range config.Sessions {
				if k == config.SessionId {
					Print("> %s", color.YellowString(v.Bucket))
				} else {
					Print("  %s", v.Bucket)
				}
			}
			return nil
		},
	}
}

func NewSwitchSessionCommand() cli.Command {
	return cli.Command{
		Name:  "switch",
		Usage: "Switch to specific session",
		Action: func(c *cli.Context) error {
			InitAndCheck(NO_LOGIN, CHECK, c)
			bucket := c.Args().First()
			for k, v := range config.Sessions {
				if bucket == v.Bucket {
					session = v
					config.SessionId = k
					saveConfigToFile()
					Print("Welcome to %s, %s!", session.Bucket, session.Operator)
					return nil
				}
			}
			PrintErrorAndExit("switch %s: No such session", bucket)
			return nil
		},
	}
}

func NewInfoCommand() cli.Command {
	return cli.Command{
		Name:  "info",
		Usage: "Current session information",
		Action: func(c *cli.Context) error {
			Init(LOGIN)
			session.Info()
			return nil
		},
	}
}

func NewMkdirCommand() cli.Command {
	return cli.Command{
		Name:  "mkdir",
		Usage: "Make directory",
		Action: func(c *cli.Context) error {
			InitAndCheck(LOGIN, CHECK, c)
			session.Mkdir(c.Args()...)
			return nil
		},
	}
}

func NewCdCommand() cli.Command {
	return cli.Command{
		Name:  "cd",
		Usage: "Change directory",
		Action: func(c *cli.Context) error {
			Init(LOGIN)
			fpath := "/"
			if c.NArg() > 0 {
				fpath = c.Args().First()
			}
			session.Cd(fpath)
			saveConfigToFile()
			return nil
		},
	}
}

func NewPwdCommand() cli.Command {
	return cli.Command{
		Name:  "pwd",
		Usage: "Print working directory",
		Action: func(c *cli.Context) error {
			Init(LOGIN)
			session.Pwd()
			return nil
		},
	}
}

func NewLsCommand() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "List directory or file",
		Action: func(c *cli.Context) error {
			Init(LOGIN)
			fpath := session.CWD
			if c.NArg() > 0 {
				fpath = c.Args().First()
			}
			mc := &MatchConfig{}
			if c.Bool("d") {
				mc.ItemType = DIR
			}
			base := path.Base(fpath)
			dir := path.Dir(fpath)
			if strings.Contains(base, "*") {
				mc.Wildcard = base
				fpath = dir
			}
			if c.String("mtime") != "" {
				err := parseMTime(c.String("mtime"), mc)
				if err != nil {
					PrintErrorAndExit("ls %s: parse mtime: %v", fpath, err)
				}
			}
			session.color = c.Bool("color")
			session.Ls(fpath, mc, c.Int("c"), c.Bool("r"))
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "r", Usage: "reverse order"},
			cli.BoolFlag{Name: "d", Usage: "only show directory"},
			cli.BoolFlag{Name: "color", Usage: "colorful output"},
			cli.IntFlag{Name: "c", Usage: "max items to list"},
			cli.StringFlag{Name: "mtime", Usage: "file's data was last modified n*24 hours ago, same as linux find command."},
		},
	}
}

func NewGetCommand() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "Get directory or file",
		Action: func(c *cli.Context) error {
			InitAndCheck(LOGIN, CHECK, c)
			upPath := c.Args().First()
			localPath := "." + string(filepath.Separator)
			if c.NArg() > 1 {
				localPath = c.Args().Get(1)
			}

			mc := &MatchConfig{}
			base := path.Base(upPath)
			dir := path.Dir(upPath)
			if strings.Contains(base, "*") {
				mc.Wildcard, upPath = base, dir
			}
			if c.String("start") != "" {
				mc.Start = c.String("start")
			}
			if c.String("end") != "" {
				mc.End = c.String("end")
			}
			if c.String("mtime") != "" {
				err := parseMTime(c.String("mtime"), mc)
				if err != nil {
					PrintErrorAndExit("get %s: parse mtime: %v", upPath, err)
				}
			}
			if mc.Start != "" || mc.End != "" {
				session.GetStartBetweenEndFiles(upPath, localPath, mc, c.Int("w"))
			} else {
				session.Get(upPath, localPath, mc, c.Int("w"))
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.IntFlag{Name: "w", Usage: "max concurrent threads", Value: 5},
			cli.StringFlag{Name: "mtime", Usage: "file's data was last modified n*24 hours ago, same as linux find command."},
			cli.StringFlag{Name: "start", Usage: "file download range starting location"},
			cli.StringFlag{Name: "end", Usage: "file download range ending location"},
		},
	}
}

func NewPutCommand() cli.Command {
	return cli.Command{
		Name:  "put",
		Usage: "Put directory or file",
		Action: func(c *cli.Context) error {
			InitAndCheck(LOGIN, CHECK, c)
			localPath := c.Args().First()
			upPath := "./"
			if c.NArg() > 1 {
				upPath = c.Args().Get(1)
			}

			session.Put(localPath, upPath, c.Int("w"))

			return nil
		},
		Flags: []cli.Flag{
			cli.IntFlag{Name: "w", Usage: "max concurrent threads", Value: 5},
		},
	}
}

func NewRmCommand() cli.Command {
	return cli.Command{
		Name:  "rm",
		Usage: "Remove directory or file",
		Action: func(c *cli.Context) error {
			InitAndCheck(LOGIN, CHECK, c)
			fpath := c.Args().First()
			base := path.Base(fpath)
			dir := path.Dir(fpath)
			mc := &MatchConfig{
				ItemType: FILE,
			}
			if strings.Contains(base, "*") {
				mc.Wildcard, fpath = base, dir
			}

			if c.Bool("d") {
				mc.ItemType = DIR
			}
			if c.Bool("a") {
				mc.ItemType = ITEM_NOT_SET
			}

			if c.String("mtime") != "" {
				err := parseMTime(c.String("mtime"), mc)
				if err != nil {
					PrintErrorAndExit("rm %s: parse mtime: %v", fpath, err)
				}
			}

			session.Rm(fpath, mc, c.Bool("async"))
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "d", Usage: "only remove directories"},
			cli.BoolFlag{Name: "a", Usage: "remove files, directories and their contents recursively, never prompt"},
			cli.BoolFlag{Name: "async", Usage: "remove asynchronously"},
			cli.StringFlag{Name: "mtime", Usage: "file's data was last modified n*24 hours ago, same as linux find command."},
		},
	}
}

func NewTreeCommand() cli.Command {
	return cli.Command{
		Name:  "tree",
		Usage: "List contents of directories in a tree-like format",
		Action: func(c *cli.Context) error {
			Init(LOGIN)
			fpath := session.CWD
			if c.NArg() > 0 {
				fpath = c.Args().First()
			}
			session.color = c.Bool("color")
			session.Tree(fpath)
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "color", Usage: "colorful output"},
		},
	}
}

func NewSyncCommand() cli.Command {
	return cli.Command{
		Name:  "sync",
		Usage: "Sync local directory to UpYun",
		Action: func(c *cli.Context) error {
			InitAndCheck(LOGIN, CHECK, c)
			localPath := c.Args().First()
			upPath := session.CWD
			if c.NArg() > 1 {
				upPath = c.Args().Get(1)
			}
			session.Sync(localPath, upPath, c.Int("w"), c.Bool("delete"), c.Bool("strong"))
			return nil
		},
		Flags: []cli.Flag{
			cli.IntFlag{Name: "w", Usage: "max concurrent threads", Value: 5},
			cli.BoolFlag{Name: "delete", Usage: "delete extraneous files from last sync"},
			cli.BoolFlag{Name: "strong", Usage: "strong consistency"},
		},
	}
}

func NewPostCommand() cli.Command {
	return cli.Command{
		Name:  "post",
		Usage: "Post async process task",
		Action: func(c *cli.Context) error {
			InitAndCheck(LOGIN, CHECK, c)
			app := c.String("app")
			notify := c.String("notify")
			task := c.String("task")
			session.PostTask(app, notify, task)
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "app", Usage: "app name"},
			cli.StringFlag{Name: "notify", Usage: "notify url"},
			cli.StringFlag{Name: "task", Usage: "task file"},
		},
	}
}

func NewPurgeCommand() cli.Command {
	return cli.Command{
		Name:  "purge",
		Usage: "refresh CDN cache",
		Action: func(c *cli.Context) error {
			InitAndCheck(LOGIN, CHECK, c)
			list := c.String("list")
			session.Purge(c.Args(), list)
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "list", Usage: "file which contains urls"},
		},
	}
}

func NewGetDBCommand() cli.Command {
	return cli.Command{
		Name:  "get-db",
		Usage: "get db value",
		Action: func(c *cli.Context) error {
			InitAndCheck(LOGIN, CHECK, c)
			if c.NArg() != 2 {
				PrintErrorAndExit("get-db local remote")
			}
			if err := initDB(); err != nil {
				PrintErrorAndExit("get-db: init database: %v", err)
			}
			value, err := getDBValue(c.Args()[0], c.Args()[1])
			if err != nil {
				PrintErrorAndExit("get-db: %v", err)
			}
			b, _ := json.MarshalIndent(value, "", "  ")
			Print("%s", string(b))
			return nil
		},
	}
}

func NewCleanDBCommand() cli.Command {
	return cli.Command{
		Name:  "clean-db",
		Usage: "clean db by local_prefx and remote_prefix",
		Action: func(c *cli.Context) error {
			InitAndCheck(LOGIN, CHECK, c)
			if c.NArg() != 2 {
				PrintErrorAndExit("clean-db local remote")
			}
			if err := initDB(); err != nil {
				PrintErrorAndExit("clean-db: init database: %v", err)
			}
			delDBValues(c.Args()[0], c.Args()[1])
			return nil
		},
	}
}

func NewUpgradeCommand() cli.Command {
	return cli.Command{
		Name:  "upgrade",
		Usage: "upgrade upx to latest version",
		Action: func(c *cli.Context) error {
			Upgrade()
			return nil
		},
	}
}
