package processbar

import (
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type UpxProcessBar struct {
	process *mpb.Progress
	enable  bool
}

var ProcessBar = &UpxProcessBar{
	process: mpb.New(
		mpb.WithWidth(100),
		mpb.WithRefreshRate(180*time.Millisecond),
		mpb.WithWaitGroup(&sync.WaitGroup{}),
	),
	enable: false,
}

func (p *UpxProcessBar) Enable() {
	p.enable = true
}

func (p *UpxProcessBar) AddBar(name string, total int64) *mpb.Bar {
	if !p.enable {
		return nil
	}

	bar := p.process.AddBar(0,
		mpb.PrependDecorators(
			decor.Name(leftAlign(shortPath(name, 30), 30), decor.WCSyncWidth),
			decor.Counters(decor.SizeB1024(0), "%.2f / %.2f", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.NewPercentage("%d", decor.WCSyncWidth),
			decor.OnComplete(
				decor.Name("...", decor.WCSyncWidth), " done",
			),
			decor.AverageSpeed(decor.SizeB1024(0), " %.1f", decor.WCSyncWidth),
		))

	bar.SetTotal(total, false)
	bar.DecoratorAverageAdjust(time.Now())
	return bar
}

func (p *UpxProcessBar) Wait() {
	if p.enable {
		p.process.Wait()
	}
}
