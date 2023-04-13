package main

import (
	"os"

	"github.com/upyun/upx"
)

func main() {
	upx.InitProgress()
	upx.Progress.Start()
	defer upx.Progress.Stop()
	upx.CreateUpxApp().Run(os.Args)
}
