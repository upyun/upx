package main

import (
	"os"

	"github.com/upyun/upx"
	"github.com/upyun/upx/processbar"
)

func main() {
	if upx.IsVerbose {
		processbar.EnableProgressbar()
		defer processbar.WaitProgressbar()
	}
	upx.CreateUpxApp().Run(os.Args)
}
