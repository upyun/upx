package main

import (
	"os"

	"github.com/upyun/upx"
)

func main() {
	upx.EnableProgressbar()
	upx.CreateUpxApp().Run(os.Args)
}
