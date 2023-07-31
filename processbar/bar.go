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

func (p *UpxProcessBar) Start() *UpxProcessBar {
	if !enableBar {
		return p
	}

	p.bar.Set("speed", 0)
	p.bar.Start()
	wg.Add(1)

	// 避免小文件传输过快导致最后计算速度异常, 此处进行短暂睡眠
	time.Sleep(time.Millisecond * 100)
	return p
}

func (p *UpxProcessBar) Finish() {
	if !enableBar {
		return
	}
	if !p.bar.IsFinished() {
		p.bar.Finish()
	}
	wg.Done()
}

func (p *UpxProcessBar) NewProxyWriter(r io.Writer) io.WriteCloser {
	return p.bar.NewProxyWriter(r)
}

func (p *UpxProcessBar) NewProxyReader(r io.Reader) io.ReadCloser {
	return p.bar.NewProxyReader(r)
}

func NewProcessBar(filename string, limit int64) *UpxProcessBar {
	bar := pb.Full.New(int(limit))
	bar.SetTemplateString(
		`{{ with string . "filename" }}{{.}}  {{end}}{{ percent . }} {{ bar . "[" ("=" | green) (cycle . "=>" | green ) "-" "]" }} ({{counters . }}, {{speed . "%s/s" "100"}})`,
	)
	bar.SetRefreshRate(time.Millisecond * 20)
	bar.Set("filename", leftAlign(shortPath(filename, 30), 30))
	return &UpxProcessBar{bar}
}

func WaitProgressbar() {
	wg.Wait()
}
