package upx

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func globFiles(patterns []string) []string {
	filenames := make([]string, 0)
	for _, filename := range patterns {
		matches, err := filepath.Glob(filename)
		if err == nil {
			filenames = append(filenames, matches...)
		}
	}
	return filenames
}

func isWindowsGOOS() bool {
	return runtime.GOOS == "windows"
}
