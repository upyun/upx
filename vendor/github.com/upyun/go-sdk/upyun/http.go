package upyun

import (
	//	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (up *UpYun) doHTTPRequest(method, url string, headers map[string]string,
	body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		if strings.ToLower(k) == "host" {
			req.Host = v
		} else {
			req.Header.Set(k, v)
		}
	}

	req.Header.Set("User-Agent", up.UserAgent)
	if method == "PUT" || method == "POST" {
		length := req.Header.Get("Content-Length")
		if length != "" {
			req.ContentLength, _ = strconv.ParseInt(length, 10, 64)
		} else {
			switch v := body.(type) {
			case *os.File:
				if fInfo, err := v.Stat(); err == nil {
					req.ContentLength = fInfo.Size()
				}
			case UpYunPutReader:
				req.ContentLength = int64(v.Len())
			}
		}
	}

	//	fmt.Printf("%+v\n", req)

	return up.httpc.Do(req)
}

func (up *UpYun) doGetEndpoint(host string) string {
	s := up.Hosts[host]
	if s != "" {
		return s
	}
	return host
}
