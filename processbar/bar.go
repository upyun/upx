package processbar

import (
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var enableBar bool = false
var barPool *pb.Pool
var wg sync.WaitGroup

func EnableProgressbar() {
	enableBar = true
}

func NewProcessBar(filename string, current, limit int64) *pb.ProgressBar {
	bar := pb.Full.New(int(limit))
	bar.SetTemplateString(
		`{{ with string . "filename" }}{{.}}  {{end}}{{ percent . }} {{ bar . "[" ("=" | green) (cycle . "=>" | green ) "-" "]" }} ({{counters .}}, {{speed . "%s/s"}})`,
	)
	bar.SetRefreshRate(time.Millisecond * 125)
	bar.Set("filename", leftAlign(shortPath(filename, 30), 30))
	return bar
}

func StartBar(bar *pb.ProgressBar) {
	if enableBar {
		if barPool == nil {
			barPool = pb.NewPool()
			barPool.Start()
		}
		wg.Add(1)
		barPool.Add(bar)
	}
}

func FinishBar(bar *pb.ProgressBar) {
	if enableBar && bar != nil {
		wg.Done()
	}
}

func WaitProgressbar() {
	if barPool != nil {
		wg.Wait()
		barPool.Stop()
	}
}
