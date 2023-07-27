package processbar

import (
	"io"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	enableBar bool = false
	wg        sync.WaitGroup
)

func EnableProgressbar() {
	enableBar = true
}

type UpxProcessBar struct {
	bar *pb.ProgressBar
}

func (p *UpxProcessBar) SetCurrent(value int64) {
	p.bar.SetCurrent(value)
}

func (p *UpxProcessBar) StartBar() {
	if !enableBar {
		return
	}
	wg.Add(1)
	p.bar.Start()
}

func (p *UpxProcessBar) Finish() {
	if !enableBar {
		return
	}
	p.bar.Finish()
	wg.Done()
}

func (p *UpxProcessBar) NewProxyWriter(r io.Writer) io.WriteCloser {
	return p.bar.NewProxyWriter(r)
}

func (p *UpxProcessBar) NewProxyReader(r io.Reader) io.ReadCloser {
	return p.bar.NewProxyReader(r)
}

func NewProcessBar(filename string, current, limit int64) *UpxProcessBar {
	bar := pb.Full.New(int(limit))
	bar.SetTemplateString(
		`{{ with string . "filename" }}{{.}}  {{end}}{{ percent . }} {{ bar . "[" ("=" | green) (cycle . "=>" | green ) "-" "]" }} ({{counters .}}, {{speed . "%s/s"}})`,
	)
	bar.SetRefreshRate(time.Millisecond * 125)
	bar.Set("filename", leftAlign(shortPath(filename, 30), 30))
	return &UpxProcessBar{bar}
}

func WaitProgressbar() {
	wg.Wait()
}
