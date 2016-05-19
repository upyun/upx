package main

import (
	"errors"
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/upyun/go-sdk/upyun"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FsDriver struct {
	// base infomation
	curDir   string
	operator string
	bucket   string

	// config
	maxConc int

	up       *upyun.UpYun
	logger   *log.Logger
	progress *uiprogress.Progress
}

func NewFsDriver(bucket, username, password, curDir string, conc int,
	logger *log.Logger) (*FsDriver, error) {
	driver := &FsDriver{
		curDir:   curDir,
		operator: username,
		bucket:   bucket,
		up:       upyun.NewUpYun(bucket, username, password),
		maxConc:  conc,
		logger:   logger,
	}

	var err error
	_, err = driver.up.Usage()
	if err != nil {
		return nil, err
	}

	driver.progress = uiprogress.New()
	driver.progress.RefreshInterval = time.Millisecond * 100
	driver.progress.Start()

	return driver, nil
}

// Make directory on UPYUN
func (driver *FsDriver) MakeDir(path string) error {
	path = driver.AbsPath(path)
	return driver.up.Mkdir(path)
}

func (dr *FsDriver) ListDirWithCount(path string,
	maxCount int) (infos []*upyun.FileInfo, err error) {
	path = dr.AbsPath(path)
	if info, err := dr.up.GetInfo(path); err != nil {
		return nil, err
	} else {
		if info.Type != "folder" {
			return []*upyun.FileInfo{info}, nil
		}
	}

	ch, errChannel := dr.up.GetLargeList(path, false, false)
	for {
		select {
		case info, more := <-ch:
			if !more {
				return infos, nil
			}
			infos = append(infos, info)
		case err := <-errChannel:
			if err != nil {
				return nil, err
			}
		}
	}

	return infos, nil
}

// Download <src> from UPYUN to <des>. <src> <des> must be file-path
func (driver *FsDriver) dlFile(src, des string) error {
	// Make dir
	if err := os.MkdirAll(filepath.Dir(des), os.ModePerm); err != nil {
		return err
	}
	fd, err := os.Create(des)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = driver.up.Get(src, fd)
	return err
}

func (driver *FsDriver) NewProgressBar(barSize int, skip bool,
	f func(src, des string) error, srcPath, desPath string) *uiprogress.Bar {
	var err error
	bar := driver.progress.AddBar(barSize).AppendCompleted()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		status := "WAIT"
		if skip {
			status = "SKIP"
		} else {
			if b.Current() == int(barSize) {
				status = "OK"
			}
		}

		if err != nil {
			return fmt.Sprintf("%-40s  ERR %s", driver.short(desPath), err)
		}
		return fmt.Sprintf("%-40s %+4s", driver.short(desPath), status)
	})

	go func() {
		err = f(srcPath, desPath)
		bar.Set(bar.Total)
	}()

	return bar
}

func (driver *FsDriver) dlFileWithProgress(src, des string) {
	src = driver.AbsPath(src)
	des, _ = filepath.Abs(des)

	barSize := 1
	upInfo, err := driver.up.GetInfo(src)
	if err == nil && upInfo.Size != 0 {
		barSize = int(upInfo.Size)
	}

	skip := false
	bar := driver.NewProgressBar(barSize, skip, driver.dlFile, src, des)

	for upInfo != nil && bar.Current() != bar.Total {
		time.Sleep(time.Millisecond * 40)
		if dkInfo, e := os.Lstat(des); e == nil {
			v := int(dkInfo.Size())
			if v == int(upInfo.Size) {
				bar.Set(bar.Total)
				break
			}
			bar.Set(v)
		}
	}
}

func (driver *FsDriver) dlDir(src, des string) {
	var wg sync.WaitGroup
	ups, _ := driver.up.GetLargeList(src, false, true)
	desDir := filepath.Join(des, path.Base(src))

	wg.Add(driver.maxConc)
	for w := 0; w < driver.maxConc; w++ {
		go func() {
			defer wg.Done()
			for {
				upInfo, more := <-ups
				if !more {
					break
				}
				if upInfo.Type == "file" {
					driver.dlFileWithProgress(path.Join(src, upInfo.Name),
						filepath.Join(desDir, upInfo.Name))
				}
			}
		}()
	}

	wg.Wait()
}

// Download <src> on UPYUN to <des> in local disk
// <src>, <des> are files or <src>, <des> are folders.
func (driver *FsDriver) Downloads(src, des string) error {
	srcPath := driver.AbsPath(src)
	if desPath, ok := driver.parseDiskDes(srcPath, des); ok {
		if driver.IsUPDir(srcPath) {
			driver.dlDir(srcPath, desPath)
		} else {
			driver.dlFileWithProgress(srcPath, desPath)
		}
		return nil
	} else {
		return errors.New("no support download folder to file.")
	}
}

func (driver *FsDriver) uploadFile(src, des string) error {
	if fd, err := os.Open(src); err == nil {
		_, err = driver.up.Put(des, fd, false, nil)
		return err
	} else {
		return err
	}
}

func (driver *FsDriver) uploadFileWithProgress(src, des string) {
	var dkInfo os.FileInfo
	var err error
	des = driver.AbsPath(des)

	barSize := 1
	if dkInfo, err = os.Lstat(src); err == nil && dkInfo.Size() != 0 {
		barSize = int(dkInfo.Size())
	}

	skip := false
	bar := driver.NewProgressBar(barSize, skip, driver.uploadFile, src, des)
	for {
		time.Sleep(time.Millisecond * 20)
		v := bar.Current()
		if v == int(barSize) {
			return
		}
		add := 102400
		if add+v < int(barSize)*98/100 {
			bar.Set(add + v)
		}
	}
}

func (driver *FsDriver) uploadDir(src, des string) {
	var wg sync.WaitGroup
	fnames := make(chan string, driver.maxConc)
	desDir := path.Join(des, filepath.Base(src))
	wg.Add(driver.maxConc)
	for w := 0; w < driver.maxConc; w++ {
		go func() {
			defer wg.Done()
			for {
				fname, more := <-fnames
				if !more {
					return
				}

				rel, _ := filepath.Rel(src, fname)
				rel = filepath.ToSlash(rel)
				desPath := path.Join(desDir, rel)
				if driver.IsDiskDir(fname) {
					driver.MakeDir(desPath)
				} else {
					driver.uploadFileWithProgress(fname, desPath)
				}
			}
		}()
	}

	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		fnames <- path
		return nil
	})

	close(fnames)
	wg.Wait()
}

// Upload <src> in local disk to <des> on UPYUN.
func (driver *FsDriver) Uploads(src, des string) error {
	desPath := driver.AbsPath(des)
	if desPath, ok := driver.parseUPYUNDes(src, desPath); ok {
		if driver.IsDiskDir(src) {
			driver.uploadDir(src, desPath)
		} else {
			driver.uploadFileWithProgress(src, desPath)
		}
		return nil
	} else {
		return errors.New("no support upload folder to file.")
	}
}

func (driver *FsDriver) rmFile(path string, async bool) {
	path = driver.AbsPath(path)
	remove := driver.up.Delete
	if async {
		remove = driver.up.AsyncDelete
	}

	err := remove(path)
	if err != nil {
		driver.logger.Printf("DELETE %s FAIL %v", path, err)
	} else {
		driver.logger.Printf("DELETE %s OK", path)
	}
}

func (driver *FsDriver) rmDir(path string, async bool) {
	// more friendly
	path = driver.AbsPath(path)
	infoChannel, errChannel := driver.up.GetLargeList(path, false, true)
	var wg sync.WaitGroup
	maxWorker := 1
	if async {
		maxWorker = 200
	}
	wg.Add(maxWorker)
	for w := 0; w < maxWorker; w++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case upInfo, more := <-infoChannel:
					if !more {
						return
					}
					driver.rmFile(path+"/"+upInfo.Name, async)
				case err := <-errChannel:
					if err != nil {
						driver.logger.Printf("rmDir GetLargeList error %v", err)
						return
					}
				}
			}
		}()
	}

	wg.Wait()
	driver.rmFile(path, async)
}

func (driver *FsDriver) Remove(path string, async bool) {
	path = driver.AbsPath(path)
	if driver.IsUPDir(path) {
		driver.rmDir(path, async)
	} else {
		driver.rmFile(path, async)
	}
}

// path MUST be a folder
func (driver *FsDriver) RemoveMatched(path string, mc *MatchConfig, async bool) {
	path = driver.AbsPath(path)
	if mc.wildcard != "" {
		upInfos, errChannel := driver.up.GetLargeList(path, false, false)
		for {
			select {
			case upInfo, more := <-upInfos:
				if !more {
					return
				}
				if mc.IsMatched(upInfo) {
					driver.Remove(path+"/"+upInfo.Name, async)
				}
			case err := <-errChannel:
				if err != nil {
					driver.logger.Printf("RemoveMatched GetLargeList error %v", err)
				}
			}
		}
	} else {
		upInfo, err := driver.up.GetInfo(path)
		if err == nil && mc.IsMatched(upInfo) {
			driver.Remove(path, async)
		} else {
			driver.logger.Printf("DELETE %s: Not matched", path)
		}
	}
}

// Get current working diretory
func (driver *FsDriver) GetCurDir() string {
	return driver.curDir
}

// Change working directory
func (driver *FsDriver) ChangeDir(path string) error {
	rPath := driver.AbsPath(path)
	if !driver.IsUPDir(rPath) {
		return errors.New(fmt.Sprintf("%s: Not a directory", rPath))
	}
	driver.curDir = rPath
	return nil
}

func (driver *FsDriver) MaybeUPDir(path string) bool {
	if driver.IsUPDir(path) || strings.HasSuffix(path, "/") {
		return true
	}
	return false
}

func (driver *FsDriver) IsUPDir(path string) bool {
	upInfo, err := driver.up.GetInfo(path)
	if err == nil {
		if upInfo.Type == "folder" {
			return true
		}
	}
	return false
}

func (driver *FsDriver) MaybeDiskDir(path string) bool {
	if driver.IsDiskDir(path) || strings.HasSuffix(path, string(filepath.Separator)) {
		return true
	}
	return false
}

func (driver *FsDriver) IsDiskDir(path string) bool {
	dkInfo, err := os.Lstat(path)
	if err == nil {
		if dkInfo.IsDir() {
			return true
		}
	}
	return false
}

func (driver *FsDriver) parseDiskDes(src, des string) (string, bool) {
	if driver.IsUPDir(src) {
		if driver.MaybeDiskDir(des) {
			return des, true
		}
		return "", false
	}
	if driver.MaybeDiskDir(des) {
		des = filepath.Join(des, path.Base(src))
	}
	return des, true
}

func (driver *FsDriver) parseUPYUNDes(src, des string) (string, bool) {
	if driver.IsDiskDir(src) {
		if driver.MaybeUPDir(des + "/") {
			return des, true
		}
		return "", false
	}
	if driver.MaybeUPDir(des) {
		des = path.Join(des, filepath.Base(src))
	}
	return des, true
}

func (dr *FsDriver) short(s string) string {
	l := len(s)
	if l <= 40 {
		return s
	}

	return s[0:17] + "..." + s[l-20:l]
}

func (driver *FsDriver) AbsPath(_path string) string {
	suffix := ""
	if strings.HasSuffix(_path, "/") {
		suffix = "/"
	}
	_path = filepath.ToSlash(_path)
	if _path[0] != '/' {
		_path = path.Join(driver.curDir, _path)
	}

	return path.Join(_path) + suffix
}
