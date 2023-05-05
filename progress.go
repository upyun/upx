package upx

import (
	"fmt"
	"os"
	"time"

	"github.com/arrebole/progressbar"
)

var enableBar bool = false

func EnableProgressbar() {
	enableBar = true
}

func AddBar(name string, total int) *progressbar.ProgressBar {
	if !enableBar {
		return nil
	}
	return progressbar.NewOptions64(
		int64(total),
		progressbar.OptionSetDescription(leftAlign(shortPath(name, 30), 30)),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)
}
