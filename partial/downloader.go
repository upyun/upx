package partial

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"
)

type Downloader interface {
	Download() error
}

type ChunkDownFunc func(start, end int64) ([]byte, error)

type MultiPartialDownloader struct {

	// 文件路径
	filePath string

	// 最终文件大小
	finalSize int64

	// 本地文件大小
	localSize int64

	//分片大小
	chunkSize int64

	writer   io.Writer
	works    int
	downFunc ChunkDownFunc
}

func NewMultiPartialDownloader(filePath string, finalSize, chunkSize int64, writer io.Writer, works int, fn ChunkDownFunc) Downloader {
	return &MultiPartialDownloader{
		filePath:  filePath,
		finalSize: finalSize,
		works:     works,
		writer:    writer,
		chunkSize: chunkSize,
		downFunc:  fn,
	}
}

func (p *MultiPartialDownloader) Download() error {
	fileinfo, err := os.Stat(p.filePath)

	// 如果异常
	// - 文件不存在异常: localSize 默认值 0
	// - 不是文件不存在异常: 报错
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		p.localSize = fileinfo.Size()
	}

	// 计算需要下载的块数
	needDownSize := p.finalSize - p.localSize
	chunkCount := needDownSize / p.chunkSize
	if needDownSize%p.chunkSize != 0 {
		chunkCount++
	}

	chunksSorter := NewChunksSorter(
		chunkCount,
		p.works,
	)

	// 下载切片任务
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		// 取消切片下载任务，并等待
		cancel()
		wg.Wait()
	}()

	for i := 0; i < p.works; i++ {
		wg.Add(1)
		go func(ctx context.Context, workId int) {
			defer func() {
				// 关闭 workId 下的接收通道
				chunksSorter.Close(workId)
				wg.Done()
			}()

			// 每个 work 取自己倍数的 chunk
			for j := workId; j < int(chunkCount); j += p.works {
				select {
				case <-ctx.Done():
					return
				default:
					var (
						err    error
						buffer []byte
					)
					start := p.localSize + int64(j)*p.chunkSize
					end := p.localSize + int64(j+1)*p.chunkSize
					if end > p.finalSize {
						end = p.finalSize
					}

					chunk := &Chunk{
						index: int64(j),
						start: start,
						end:   end,
					}

					// 重试三次
					for t := 0; t < 3; t++ {
						// ? 由于长度是从1开始，而数据是从0地址开始
						// ? 计算字节时容量会多出开头的一位，所以末尾需要减少一位
						buffer, err = p.downFunc(chunk.start, chunk.end-1)
						if err == nil {
							break
						}
					}
					chunk.SetData(buffer)
					chunk.SetError(err)
					chunksSorter.Write(chunk)

					if err != nil {
						return
					}
				}
			}
		}(ctx, i)
	}

	// 将分片顺序写入到文件
	for {
		chunk := chunksSorter.Read()
		if chunk == nil {
			break
		}
		if chunk.Error() != nil {
			return chunk.Error()
		}
		if len(chunk.Data()) == 0 {
			return errors.New("chunk buffer download but size is 0")
		}
		p.writer.Write(chunk.Data())
	}
	return nil
}
