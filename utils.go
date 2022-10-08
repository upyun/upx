package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"github.com/upyun/go-sdk/v3/upyun"
	"encoding/json"
	"path"
	"bufio"
)

func shortPath(s string, width int) string {
	if slen(s) <= width {
		return s
	}

	dotLen := 3
	headLen := (width - dotLen) / 2
	tailLen := width - dotLen - headLen

	st := 1
	for ; st < len(s); st++ {
		if slen(s[0:st]) > headLen {
			break
		}
	}

	ed := len(s) - 1
	for ; ed >= 0; ed-- {
		if slen(s[ed:]) > tailLen {
			break
		}
	}

	return s[0:st-1] + strings.Repeat(".", dotLen) + s[ed+1:]
}

func leftAlign(s string, width int) string {
	l := slen(s)
	for i := 0; i < width-l; i++ {
		s += " "
	}
	return s
}
func rightAlign(s string, width int) string {
	l := slen(s)
	for i := 0; i < width-l; i++ {
		s = " " + s
	}
	return s
}

func slen(s string) int {
	l, rl := len(s), len([]rune(s))
	return (l-rl)/2 + rl
}

func parseMTime(value string, match *MatchConfig) error {
	if value == "" {
		return nil
	}

	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}

	if v < 0 {
		match.After = time.Now().Add(time.Duration(v) * time.Hour * 24)
		match.TimeType = TIME_AFTER
	} else {
		if strings.HasPrefix(value, "+") {
			match.Before = time.Now().Add(time.Duration(-1*(v+1)) * time.Hour * 24)
			match.TimeType = TIME_BEFORE
		} else {
			match.Before = time.Now().Add(time.Duration(-1*v) * time.Hour * 24)
			match.After = time.Now().Add(time.Duration(-1*(v+1)) * time.Hour * 24)
			match.TimeType = TIME_INTERVAL
		}
	}
	return nil
}

func humanizeSize(b int64) string {
	unit := []string{"B", "KB", "MB", "GB", "TB"}
	u, v, s := 0, float64(b), ""
	for {
		if v < 1024.0 {
			switch {
			case v < 10:
				s = fmt.Sprintf("%.3f", v)
			case v < 100:
				s = fmt.Sprintf("%.2f", v)
			case v < 1000:
				s = fmt.Sprintf("%.1f", v)
			default:
				s = fmt.Sprintf("%.0f", v)
			}
			break
		}
		v /= 1024
		u++
	}

	if strings.Contains(s, ".") {
		ed := len(s) - 1
		for ; ed > 0; ed-- {
			if s[ed] == '.' {
				ed--
				break
			}
			if s[ed] != '0' {
				break
			}
		}
		s = s[:ed+1]
	}
	return s + unit[u]
}

func md5File(fpath string) (string, error) {
	fd, err := os.Open(fpath)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, fd); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func walk(root string, f func(string, os.FileInfo, error)) {
	fi, err := os.Stat(root)
	if err == nil && fi != nil && fi.IsDir() {
		fInfos, err := ioutil.ReadDir(root)
		f(root, fi, err)
		for _, fInfo := range fInfos {
			walk(filepath.Join(root, fInfo.Name()), f)
		}
	} else {
		f(root, fi, err)
	}
}

var Temp = `resumeRecorder`

func writeInfo(point *upyun.BreakPointConfig) (string, error) {
	//保存信息
	data, err := json.Marshal(point)
	if err != nil {
		Print("marshal data failed")
		return "", err
	}
	//生成临时文件
	TempPath := path.Join(os.TempDir(), Temp, point.UploadID)

	file, err := os.OpenFile(TempPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		PrintErrorAndExit("open file %s failed", TempPath)
	}

	write := bufio.NewWriter(file)
	_, err = write.WriteString(time.Now().Format("2006-01-02-15-04"))
	if err != nil {
		PrintErrorAndExit("write string error, error= %s", err)
		return "", err
	}
	_, err = write.Write(data)
	if err != nil {
		PrintErrorAndExit("write data to file error, error= %s", err)
	}
	_, err = write.WriteString("\n")
	if err != nil {
		PrintErrorAndExit("write string error, error= %s", err)
	}
	if err = write.Flush(); err != nil {
		PrintErrorAndExit("write flush error, error= %s", err)
	}

	return TempPath, nil
}
