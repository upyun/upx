package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"runtime"
	"sort"
)

var cmds = []string{
	"login", "logout", "cd", "pwd", "get", "put", "sync",
	"ls", "rm", "switch", "info", "mkdir", "services",
}

var version = "v0.1.2"

func main() {
	app := cli.NewApp()
	app.Name = "upx"
	app.Usage = "a tool for managing files in UPYUN"
	app.Author = "Hongbo.Mo"
	app.Email = "zjutpolym@gmail.com"
	app.Version = fmt.Sprintf("%s %s/%s %s", version, runtime.GOOS,
		runtime.GOARCH, runtime.Version())
	app.Commands = make([]cli.Command, 0)

	sort.Strings(cmds)
	for _, cmd := range cmds {
		cm, exist := CmdMap[cmd]
		if exist {
			Cmd := cli.Command{
				Name:  cmd,
				Usage: cm.Desc,
				Action: func(c *cli.Context) error {
					if c.Command.FullName() != "login" && driver == nil {
						fmt.Println("Log in first.")
						os.Exit(-1)
					}
					opts := make(map[string]interface{})
					for k, v := range cm.Flags {
						if c.IsSet(k) {
							switch v.typ {
							case "bool":
								opts[k] = c.Bool(k)
							case "string":
								opts[k] = c.String(k)
							case "int":
								opts[k] = c.Int(k)
							}
						}
					}
					cm.Func(c.Args(), opts)
					return nil
				},
			}
			if cm.Alias != "" {
				Cmd.Aliases = []string{cm.Alias}
			}
			if cm.Flags != nil {
				Cmd.Flags = []cli.Flag{}
				for k, v := range cm.Flags {
					var flag cli.Flag
					switch v.typ {
					case "bool":
						flag = cli.BoolFlag{Name: k, Usage: v.usage}
					case "int":
						flag = cli.StringFlag{Name: k, Usage: v.usage}
					case "string":
						flag = cli.IntFlag{Name: k, Usage: v.usage}
					}
					Cmd.Flags = append(Cmd.Flags, flag)
				}
			}

			app.Commands = append(app.Commands, Cmd)
		}
	}

	app.Run(os.Args)
}
