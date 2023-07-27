package main

import (
	"os"

	"github.com/upyun/upx"
	"github.com/upyun/upx/processbar"
)

func main() {
	processbar.EnableProgressbar()
	upx.CreateUpxApp().Run(os.Args)
	processbar.WaitProgressbar()
}
