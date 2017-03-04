package upyun

import (
	"net/http"
	"strings"
	"time"
)

type FileInfo struct {
	Name        string
	Size        int64
	ContentType string
	IsDir       bool
	ETag        string
	Time        time.Time

	Meta map[string]string

	/* image information */
	ImgType   string
	ImgWidth  int64
	ImgHeight int64
	ImgFrames int64
}

/*
  Content-Type: image/gif
  ETag: "dc9ea7257aa6da18e74505259b04a946"
  x-upyun-file-type: GIF
  x-upyun-height: 379
  x-upyun-width: 500
  x-upyun-frames: 90
*/
func parseHeaderToFileInfo(header http.Header, getinfo bool) *FileInfo {
	fInfo := &FileInfo{}
	for k, v := range header {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "x-upyun-meta-") {
			if fInfo.Meta == nil {
				fInfo.Meta = make(map[string]string)
			}
			fInfo.Meta[lk] = v[0]
		}
	}

	if getinfo {
		// HTTP HEAD
		fInfo.Size = parseStrToInt(header.Get("x-upyun-file-size"))
		fInfo.IsDir = header.Get("x-upyun-file-type") == "folder"
		fInfo.Time = time.Unix(parseStrToInt(header.Get("x-upyun-file-date")), 0)
	} else {
		fInfo.Size = parseStrToInt(header.Get("Content-Length"))
		fInfo.ContentType = header.Get("Content-Type")
		fInfo.ETag = header.Get("ETag")
		fInfo.ImgType = header.Get("x-upyun-file-type")
		fInfo.ImgWidth = parseStrToInt(header.Get("x-upyun-width"))
		fInfo.ImgHeight = parseStrToInt(header.Get("x-upyun-height"))
		fInfo.ImgFrames = parseStrToInt(header.Get("x-upyun-frames"))
	}
	return fInfo
}
