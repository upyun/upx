package partial

import (
	"sync/atomic"
)

type Chunk struct {
	// 切片的顺序
	index int64

	// 切片内容的在源文件的开始地址
	start int64

	// 切片内容在源文件的结束地址
	end int64

	// 切片任务的下载错误
	err error

	// 下载完的切片的具体内容
	buffer []byte
}

func (p *Chunk) SetData(bytes []byte) {
	p.buffer = bytes
}

func (p *Chunk) SetError(err error) {
	p.err = err
}

func (p *Chunk) Error() error {
	return p.err
}

func (p *Chunk) Data() []byte {
	return p.buffer
}

// 切片乱序写入后，将切片顺序读取
type ChunksSorter struct {
	// 已经读取的切片数量
	readCount int64

	// 切片的所有总数
	chunkCount int64

	// 线程数，用于阻塞写入
	works int64

	// 存储切片的缓存区
	chunks []chan *Chunk
}

func NewChunksSorter(chunkCount int64, works int) *ChunksSorter {
	chunks := make([]chan *Chunk, works)
	for i := 0; i < len(chunks); i++ {
		chunks[i] = make(chan *Chunk)
	}

	return &ChunksSorter{
		chunkCount: chunkCount,
		works:      int64(works),
		chunks:     chunks,
	}
}

// 将数据写入到缓存区，如果该缓存已满，则会被阻塞
func (p *ChunksSorter) Write(chunk *Chunk) {
	p.chunks[chunk.index%p.works] <- chunk
}

// 关闭 workId 下的通道
func (p *ChunksSorter) Close(workId int) {
	if (len(p.chunks) - 1) >= workId {
		close(p.chunks[workId])
	}
}

// 顺序读取切片，如果下一个切片没有下载完，则会被阻塞
func (p *ChunksSorter) Read() *Chunk {
	if p.chunkCount == 0 {
		return nil
	}
	i := atomic.AddInt64(&p.readCount, 1)
	chunk := <-p.chunks[(i-1)%p.works]
	return chunk
}
