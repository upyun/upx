package partial

import (
	"context"
	"io"

	"golang.org/x/sync/errgroup"
)

// 下载器接口
type Uploader interface {
	Upload() error
}

// 具体的分片下载函数
type ChunkUploadFunc func(partId int64, body []byte) error

// 多线程分片上传器
type MultiPartialUploader struct {
	//分片大小
	chunkSize int64

	// 本地文件
	reader io.Reader

	// 线程数
	works int

	// 上传函数
	handleFunc ChunkUploadFunc
}

func NewMultiPartialUploader(chunkSize int64, reader io.Reader, works int, fn ChunkUploadFunc) Uploader {
	if works <= 0 {
		panic("multiPartialUploader works must > 0")
	}
	if chunkSize <= 0 {
		panic("multiPartialUploader chunkSize must > 0")
	}

	return &MultiPartialUploader{
		works:      works,
		chunkSize:  chunkSize,
		handleFunc: fn,
		reader:     reader,
	}
}

func (p *MultiPartialUploader) Upload() error {
	chunkUploadTask := make(chan *Chunk, p.works)

	// 任务发布者
	// 从reader中读取分片大小的数据，提交到上传任务队列
	go func() {
		var chunkIndex int64 = 0
		for {
			// 已经上传完成则跳过
			buffer := make([]byte, p.chunkSize)
			nRead, err := p.reader.Read(buffer)

			chunkUploadTask <- &Chunk{
				index:  chunkIndex,
				buffer: buffer[0:nRead],
				start:  p.chunkSize * chunkIndex,
				end:    p.chunkSize*chunkIndex + int64(nRead),
				err:    err,
			}
			if err != nil {
				break
			}
			chunkIndex++
		}
		close(chunkUploadTask)
	}()

	// 上传任务到云端
	group, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < p.works; i++ {
		group.Go(func() error {
			if err := p.uploadChunks(ctx, chunkUploadTask); err != nil {
				return err
			}
			return nil
		})
	}
	return group.Wait()
}

func (p *MultiPartialUploader) uploadChunks(ctx context.Context, channel <-chan *Chunk) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case chunk, ok := <-channel:
			if !ok {
				return nil
			}
			if chunk.err != nil {
				if chunk.err == io.EOF {
					return nil
				}
				return chunk.err
			}
			if err := p.handleFunc(chunk.index, chunk.buffer); err != nil {
				return err
			}
		}
	}
}
