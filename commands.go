package upx

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/upyun/upx/xerrors"
	"github.com/urfave/cli"
	"golang.org/x/term"
)

const (
	NO_CHECK = false
	CHECK    = true
)

func InitAndCheck(login, check bool, c *cli.Context) (err error) {
	if login == LOGIN && session == nil {
		err = readConfigFromFile(LOGIN)
	}
	if login == NO_LOGIN {
		err = readConfigFromFile(NO_LOGIN)
	}

	if check && c.NArg() == 0 && c.NumFlags() == 0 {
		err = xerrors.ErrInvalidCommand
	}
	return
}

func CreateInitCheckFunc(login, check bool) cli.BeforeFunc {
	return func(ctx *cli.Context) error {
		if err := InitAndCheck(login, check, ctx); err != nil {
			if errors.Is(err, xerrors.ErrInvalidCommand) {
				cli.ShowCommandHelp(ctx, ctx.Command.Name)
				return &cli.ExitError{}
			}
			return cli.NewExitError(err.Error(), -1)
		}
		return nil
	}
}

func NewLoginCommand() cli.Command {
	return cli.Command{
		Name:   "login",
		Usage:  "Log in to UpYun",
		Before: CreateInitCheckFunc(NO_LOGIN, NO_CHECK),
		Action: func(c *cli.Context) error {
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
				b, err := term.ReadPassword(int(os.Stdin.Fd()))
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
		Name:   "logout",
		Usage:  "Log out of your UpYun account",
		Before: CreateInitCheckFunc(NO_LOGIN, NO_CHECK),
		Action: func(c *cli.Context) error {
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
		Name:   "sessions",
		Usage:  "List all sessions",
		Before: CreateInitCheckFunc(NO_LOGIN, NO_CHECK),
		Action: func(c *cli.Context) error {
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
		Name:   "switch",
		Usage:  "Switch to specific session",
		Before: CreateInitCheckFunc(NO_LOGIN, CHECK),
		Action: func(c *cli.Context) error {
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
		Name:   "info",
		Usage:  "Current session information",
		Before: CreateInitCheckFunc(LOGIN, NO_CHECK),
		Action: func(c *cli.Context) error {
			session.Info()
			return nil
		},
	}
}

func NewMkdirCommand() cli.Command {
	return cli.Command{
		Name:      "mkdir",
		Usage:     "Make directory",
		ArgsUsage: "<remote-dir>",
		Before:    CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
			session.Mkdir(c.Args()...)
			return nil
		},
	}
}

func NewCdCommand() cli.Command {
	return cli.Command{
		Name:      "cd",
		Usage:     "Change directory",
		ArgsUsage: "<remote-path>",
		Before:    CreateInitCheckFunc(LOGIN, NO_CHECK),
		Action: func(c *cli.Context) error {
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
		Name:   "pwd",
		Usage:  "Print working directory",
		Before: CreateInitCheckFunc(LOGIN, NO_CHECK),
		Action: func(c *cli.Context) error {
			session.Pwd()
			return nil
		},
	}
}

func NewLsCommand() cli.Command {
	return cli.Command{
		Name:      "ls",
		Usage:     "List directory or file",
		ArgsUsage: "<remote-path>",
		Before:    CreateInitCheckFunc(LOGIN, NO_CHECK),
		Action: func(c *cli.Context) error {
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
		Name:      "get",
		Usage:     "Get directory or file",
		ArgsUsage: "[-c] <remote-path> [save-path]",
		Before:    CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
			upPath := c.Args().First()
			localPath := "." + string(filepath.Separator)

			if c.NArg() > 2 {
				PrintErrorAndExit("upx get args limit 2")
			}
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
			if c.Int("w") > 10 || c.Int("w") < 1 {
				PrintErrorAndExit("max concurrent threads must between (1 - 10)")
			}
			if mc.Start != "" || mc.End != "" {
				session.GetStartBetweenEndFiles(upPath, localPath, mc, c.Int("w"))
			} else {
				session.Get(upPath, localPath, mc, c.Int("w"), c.Bool("c"))
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.IntFlag{Name: "w", Usage: "max concurrent threads (1-10)", Value: 5},
			cli.BoolFlag{Name: "c", Usage: "continue download, Resume Broken Download"},
			cli.StringFlag{Name: "mtime", Usage: "file's data was last modified n*24 hours ago, same as linux find command."},
			cli.StringFlag{Name: "start", Usage: "file download range starting location"},
			cli.StringFlag{Name: "end", Usage: "file download range ending location"},
		},
	}
}

func NewPutCommand() cli.Command {
	return cli.Command{
		Name:      "put",
		Usage:     "Put directory or file",
		ArgsUsage: "<local-path> [remote-path]",
		Before:    CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
			localPath := c.Args().First()
			upPath := "./"

			if c.NArg() > 2 {
				fmt.Println("Use the upload command instead of the put command for multiple file uploads")
				os.Exit(0)
			}

			if c.NArg() > 1 {
				upPath = c.Args().Get(1)
			}
			if c.Int("w") > 10 || c.Int("w") < 1 {
				PrintErrorAndExit("max concurrent threads must between (1 - 10)")
			}
			session.Put(
				localPath,
				upPath,
				c.Int("w"),
				c.Bool("all"),
			)
			return nil
		},
		Flags: []cli.Flag{
			cli.IntFlag{Name: "w", Usage: "max concurrent threads", Value: 5},
			cli.BoolFlag{Name: "all", Usage: "upload all files including hidden files"},
		},
	}
}

func NewUploadCommand() cli.Command {
	return cli.Command{
		Name:      "upload",
		Usage:     "upload multiple directory or file or http url",
		ArgsUsage: "[local-path...] [url...] [--remote remote-path]",
		Before:    CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
			if c.Int("w") > 10 || c.Int("w") < 1 {
				PrintErrorAndExit("max concurrent threads must between (1 - 10)")
			}
			filenames := c.Args()
			if isWindowsGOOS() {
				filenames = globFiles(filenames)
			}
			session.Upload(
				filenames,
				c.String("remote"),
				c.Int("w"),
				c.Bool("all"),
			)
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all", Usage: "upload all files including hidden files"},
			cli.IntFlag{Name: "w", Usage: "max concurrent threads", Value: 5},
			cli.StringFlag{Name: "remote", Usage: "remote path", Value: "./"},
		},
	}
}

func NewRmCommand() cli.Command {
	return cli.Command{
		Name:      "rm",
		Usage:     "Remove directory or file",
		ArgsUsage: "<remote-path>",
		Before:    CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
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
		Name:      "tree",
		Usage:     "List contents of directories in a tree-like format",
		ArgsUsage: "<remote-path>",
		Before:    CreateInitCheckFunc(LOGIN, NO_CHECK),
		Action: func(c *cli.Context) error {
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
		Name:      "sync",
		Usage:     "Sync local directory to UpYun",
		ArgsUsage: "<local-path> [remote-path]",
		Before:    CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
			localPath := c.Args().First()
			upPath := session.CWD
			if c.NArg() > 1 {
				upPath = c.Args().Get(1)
			}
			if c.Int("w") > 10 || c.Int("w") < 1 {
				PrintErrorAndExit("max concurrent threads must between (1 - 10)")
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
		Name:   "post",
		Usage:  "Post async process task",
		Before: CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
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
		Name:   "purge",
		Usage:  "refresh CDN cache",
		Before: CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
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
		Name:   "get-db",
		Usage:  "get db value",
		Before: CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
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
		Name:   "clean-db",
		Usage:  "clean db by local_prefx and remote_prefix",
		Before: CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
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

func NewCopyCommand() cli.Command {
	return cli.Command{
		Name:      "cp",
		Usage:     "copy files inside cloud storage",
		ArgsUsage: "[remote-source-path] [remote-target-path]",
		Before:    CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				PrintErrorAndExit("invalid command args")
			}
			if err := session.Copy(c.Args()[0], c.Args()[1], c.Bool("f")); err != nil {
				PrintErrorAndExit(err.Error())
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "f", Usage: "Force overwrite existing files"},
		},
	}
}

func NewMoveCommand() cli.Command {
	return cli.Command{
		Name:      "mv",
		Usage:     "move files inside cloud storage",
		ArgsUsage: "[remote-source-path] [remote-target-path]",
		Before:    CreateInitCheckFunc(LOGIN, CHECK),
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				PrintErrorAndExit("invalid command args")
			}
			if err := session.Move(c.Args()[0], c.Args()[1], c.Bool("f")); err != nil {
				PrintErrorAndExit(err.Error())
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "f", Usage: "Force overwrite existing files"},
		},
	}
}
