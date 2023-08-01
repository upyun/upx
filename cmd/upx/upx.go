package main

import (
	"os"

	"github.com/upyun/upx"
	"github.com/upyun/upx/processbar"
)

func main() {
	if upx.IsVerbose {
		processbar.ProcessBar.Enable()
		defer processbar.ProcessBar.Wait()
	}
	upx.CreateUpxApp().Run(os.Args)
}
