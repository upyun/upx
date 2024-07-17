package upx

import (
	"fmt"
	"runtime"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/urfave/cli"
)

const VERSION = "v0.4.6"

func CreateUpxApp() *cli.App {
	app := cli.NewApp()
	app.Name = "upx"
	app.Usage = "a tool for driving UpYun Storage"
	app.Author = "Hongbo.Mo"
	app.Email = "zjutpolym@gmail.com"
	app.Version = fmt.Sprintf("%s %s/%s %s", VERSION,
		runtime.GOOS, runtime.GOARCH, runtime.Version())
	app.EnableBashCompletion = true
	app.Compiled = time.Now()
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "quiet, q", Usage: "not verbose"},
		cli.StringFlag{Name: "auth", Usage: "auth string"},
		cli.StringFlag{Name: "pprof-addr", Usage: "enable pprof"},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("q") {
			IsVerbose = false
		}
		if c.String("auth") != "" {
			err := authStrToConfig(c.String("auth"))
			if err != nil {
				PrintErrorAndExit("%s: invalid auth string", c.Command.FullName())
			}
		}
		if c.String("pprof-addr") != "" {
			go func() {
				http.ListenAndServe(c.String("pprof-addr"), nil)
			}()
		}
		return nil
	}
	app.Commands = []cli.Command{
		NewLoginCommand(),
		NewLogoutCommand(),
		NewListSessionsCommand(),
		NewSwitchSessionCommand(),
		NewInfoCommand(),
		NewCdCommand(),
		NewPwdCommand(),
		NewMkdirCommand(),
		NewLsCommand(),
		NewTreeCommand(),
		NewGetCommand(),
		NewPutCommand(),
		NewUploadCommand(),
		NewRmCommand(),
		NewSyncCommand(),
		NewAuthCommand(),
		NewPostCommand(),
		NewPurgeCommand(),
		NewGetDBCommand(),
		NewCleanDBCommand(),
		NewUpgradeCommand(),
		NewCopyCommand(),
		NewMoveCommand(),
	}
	return app
}
