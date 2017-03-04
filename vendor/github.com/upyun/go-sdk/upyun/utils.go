package upyun

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func makeRFC1123Date(d time.Time) string {
	utc := d.UTC().Format(time.RFC1123)
	return strings.Replace(utc, "UTC", "GMT", -1)
}

func makeUserAgent(version string) string {
	return fmt.Sprintf("UPYUN Go SDK V2/%s", version)
}

func md5Str(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func base64ToStr(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func hmacSha1(key string, data []byte) []byte {
	hm := hmac.New(sha1.New, []byte(key))
	hm.Write(data)
	return hm.Sum(nil)
}

func escapeUri(uri string) (string, error) {
	uri = path.Join("/", uri)
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

var readHTTPBody = ioutil.ReadAll

func readHTTPBodyToStr(resp *http.Response) (string, error) {
	b, err := readHTTPBody(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("read http body: %v", err)
	}
	return string(b), nil
}

func addQueryToUri(rawurl string, kwargs map[string]string) string {
	u, _ := url.ParseRequestURI(rawurl)
	q := u.Query()
	for k, v := range kwargs {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func encodeQueryToPayload(kwargs map[string]string) string {
	payload := url.Values{}
	for k, v := range kwargs {
		payload.Set(k, v)
	}
	return payload.Encode()
}

func readHTTPBodyToInt(resp *http.Response) (int64, error) {
	b, err := readHTTPBody(resp.Body)
	resp.Body.Close()
	if err != nil {
		return 0, fmt.Errorf("read http body: %v", err)
	}

	n, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse int: %v", err)
	}
	return n, nil
}

func parseStrToInt(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func md5File(f io.ReadSeeker) (string, error) {
	offset, _ := f.Seek(0, 0)
	defer f.Seek(offset, 0)
	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func parseBodyToFileInfos(b []byte) (fInfos []*FileInfo) {
	line := strings.Split(string(b), "\n")
	for _, l := range line {
		if len(l) == 0 {
			continue
		}
		items := strings.Split(l, "\t")
		if len(items) != 4 {
			continue
		}

		fInfos = append(fInfos, &FileInfo{
			Name:  items[0],
			IsDir: items[1] == "F",
			Size:  int64(parseStrToInt(items[2])),
			Time:  time.Unix(parseStrToInt(items[3]), 0),
		})
	}
	return
}
