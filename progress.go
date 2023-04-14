package upx

import (
	"time"

	"github.com/gosuri/uiprogress"
)

var Progress *uiprogress.Progress

func AddBar(id, total int) (*uiprogress.Bar, int) {
	if id >= len(Progress.Bars) || id < 0 {
		return Progress.AddBar(total), len(Progress.Bars) - 1
	} else {
		Progress.Bars[id] = uiprogress.NewBar(total)
		return Progress.Bars[id], id
	}
}

func InitProgress() {
	Progress = uiprogress.New()
	Progress.RefreshInterval = time.Millisecond * 100
}
