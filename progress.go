package main

import (
	"github.com/gosuri/uiprogress"
	"time"
)

var progress *uiprogress.Progress

func AddBar(id, total int) (*uiprogress.Bar, int) {
	if id >= len(progress.Bars) || id < 0 {
		return progress.AddBar(total), len(progress.Bars) - 1
	} else {
		progress.Bars[id] = uiprogress.NewBar(total)
		return progress.Bars[id], id
	}
}

func initProgress() {
	progress = uiprogress.New()
	progress.RefreshInterval = time.Millisecond * 100
}
