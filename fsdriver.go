// +build linux darwin

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
	curDir   string
	operator string
	bucket   string
	maxConc  int
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

func short(s string) string {
	l := len(s)
	if l <= 40 {
		return s
	}

	return s[0:17] + "..." + s[l-20:l]
}

func (dr *FsDriver) ChangeDir(path string) error {
	path = dr.abs(path)
	if info, err := dr.up.GetInfo(path); err != nil {
		return err
	} else {
		if info.Type == "folder" {
			dr.curDir = dr.abs(path + "/")
			fmt.Println(dr.curDir)
			return nil
		}
		return errors.New(fmt.Sprintf("%s: Not a directory", path))
	}
}
func (dr *FsDriver) GetCurDir() string {
	return dr.curDir
}

func (dr *FsDriver) getItem(src, des string) error {
	var dkInfo os.FileInfo
	var fd *os.File
	var upInfo *upyun.FileInfo
	var err error
	var wg sync.WaitGroup
	var skip bool = false

	if upInfo, err = dr.up.GetInfo(src); err != nil {
		return err
	}

	if dkInfo, err = os.Lstat(des); err == nil && dkInfo.Size() == upInfo.Size {
		skip = true
	}

	barSize := upInfo.Size
	// hack for empty file
	if barSize == 0 {
		barSize = 1
	}

	bar := dr.progress.AddBar(int(barSize)).AppendCompleted()
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
			return fmt.Sprintf("%-40s  ERR %s", short(src), err)
		}
		return fmt.Sprintf("%-40s %+4s", short(src), status)
	})

	wg.Add(1)
	go func() {
		v := 0
		defer wg.Done()
		for {
			var verr error
			time.Sleep(time.Millisecond * 40)
			if dkInfo, verr = os.Lstat(des); verr == nil {
				v = int(dkInfo.Size())
			}
			bar.Set(v)
			if v == int(upInfo.Size) {
				if v == 0 {
					bar.Set(1)
				}
				return
			}
		}
	}()

	if !skip {
		if fd, err = os.OpenFile(des, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600); err == nil {
			defer fd.Close()
			if err = dr.up.Get(src, fd); err != nil {
				return err
			}
		}
	}

	wg.Wait()

	return err
}

func (dr *FsDriver) GetItems(src, des string) error {
	var ch chan *upyun.FileInfo

	src = dr.abs(src)
	if ok, err := dr.IsDir(src); err != nil {
		return err
	} else {
		if ok {
			ch = dr.up.GetLargeList(src, true)
		} else {
			ch = make(chan *upyun.FileInfo, 10)
			ch <- &upyun.FileInfo{
				Type: "file",
				Name: "",
			}
			close(ch)
		}
	}

	var wg sync.WaitGroup
	for w := 0; w < dr.maxConc; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				upInfo, more := <-ch
				if !more {
					return
				}

				if upInfo.Type != "folder" {
					// if des is a directory, should add basename to des filename
					srcPath := src
					desPath := des
					if upInfo.Name != "" {
						srcPath += "/" + upInfo.Name
						desPath += "/" + upInfo.Name
					} else {
						if info, err := os.Lstat(des); err == nil && info.IsDir() {
							desPath += "/" + path.Base(srcPath)
						}
					}

					srcPath = dr.abs(srcPath)
					desPath = path.Clean(desPath)
					var err error
					if err = os.MkdirAll(path.Dir(desPath), os.ModePerm); err == nil {
						err = dr.getItem(srcPath, desPath)
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}

func (dr *FsDriver) ListDir(path string) (infos []*upyun.FileInfo, err error) {
	path = dr.abs(path)
	if info, err := dr.up.GetInfo(path); err != nil {
		return nil, err
	} else {
		if info.Type != "folder" {
			return []*upyun.FileInfo{info}, nil
		}
	}

	ch := dr.up.GetLargeList(path, false)
	for k := 0; k < 1000; k++ {
		info, more := <-ch
		if !more {
			return infos[0:k], nil
		}
		infos = append(infos, info)
	}

	close(ch)

	return infos, nil
}

func (dr *FsDriver) MakeDir(path string) error {
	path = dr.abs(path)
	fmt.Println(path)
	if err := dr.up.Mkdir(path); err != nil {
		return err
	}
	return nil
}

func (dr *FsDriver) putItem(src, des string) error {
	var dkInfo os.FileInfo
	var upInfo *upyun.FileInfo
	var err error
	var skip bool = false
	var wg sync.WaitGroup

	if dkInfo, err = os.Lstat(src); err != nil {
		return err
	}

	if upInfo, err = dr.up.GetInfo(des); err == nil && dkInfo.Size() == upInfo.Size {
		skip = true
	}

	err = nil
	barSize := dkInfo.Size()
	if barSize == 0 {
		barSize = 1
	}

	bar := dr.progress.AddBar(int(barSize)).AppendCompleted()
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
			return fmt.Sprintf("%-40s  ERR %s", short(des), err)
		}
		return fmt.Sprintf("%-40s %+4s", short(des), status)
	})

	wg.Add(1)
	lock := new(sync.Mutex)
	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Millisecond * 20)
			lock.Lock()
			v := bar.Current()
			if v == int(barSize) || err != nil {
				lock.Unlock()
				return
			}
			add := 102400
			if add+v < int(barSize)*98/100 {
				bar.Set(add + v)
			}
			lock.Unlock()
		}
	}()

	if !skip {
		var fd *os.File
		if fd, err = os.OpenFile(src, os.O_RDWR, 0600); err == nil {
			_, err = dr.up.Put(des, fd, false, "", "", nil)
			fd.Close()
		}
	}

	// hack
	if err == nil {
		lock.Lock()
		bar.Set(int(barSize))
		lock.Unlock()
	}

	wg.Wait()

	return err
}

func (dr *FsDriver) PutItems(src, des string) error {
	var wg sync.WaitGroup
	ch := make(chan string, dr.maxConc+10)
	des = dr.abs(des)

	isUpDir, err := dr.IsDir(des)
	if err != nil && err.Error() != "404" {
		return err
	}

	srcInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// upload items to the directory which doesn't exist
	if strings.HasSuffix(des, "/") {
		isUpDir = true
	}

	for w := 0; w < dr.maxConc; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				srcPath, more := <-ch
				if !more {
					return
				}

				var desPath string

				desPath = des
				if isUpDir {
					if srcInfo.IsDir() {
						desPath += "/" + strings.Replace(srcPath, src, "", 1)
					} else {
						desPath += "/" + filepath.Base(srcPath)
					}
				}

				if strings.HasSuffix(desPath, "/") {
					panic(fmt.Sprintf("desPath should not HasSuffix with / %s", desPath))
				}

				desPath = dr.abs(desPath)

				dr.putItem(srcPath, desPath)
			}
		}()
	}

	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			ch <- path
		}
		return nil
	})

	close(ch)
	wg.Wait()

	return nil
}

func (dr *FsDriver) Remove(path string) error {
	path = dr.abs(path)
	ok, err := dr.IsDir(path)
	if err != nil {
		return err
	}

	if ok {
		// remove all items in directory
		var wg sync.WaitGroup
		ch := dr.up.GetLargeList(path, true)
		// UPYUN Delete Limit Rate
		maxWorker := 1

		for w := 0; w < maxWorker; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					info, more := <-ch
					if !more {
						return
					}
					desPath := dr.abs(path + "/" + info.Name)
					// TODO: retry
					if err := dr.up.Delete(desPath); err != nil {
						//TODO: error
						dr.logger.Printf("DELETE %s FAIL %v", desPath, err)
					} else {
						dr.logger.Printf("DELETE %s OK", desPath)
					}
				}
			}()
		}

		wg.Wait()
	}

	if err = dr.up.Delete(path); err != nil {
		dr.logger.Printf("DELETE %s FAIL %v", path, err)
	} else {
		dr.logger.Printf("DELETE %s OK", path)
	}

	// TODO: more information
	return nil
}

func (dr *FsDriver) IsDir(path string) (bool, error) {
	path = dr.abs(path)
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	if info, err := dr.up.GetInfo(path); err != nil || info.Type != "folder" {
		return false, err
	}
	return true, nil
}

func (dr *FsDriver) abs(path string) string {
	if path[0] != '/' {
		path = dr.curDir + "/" + path
	}

	if strings.HasSuffix(path, "/.") || strings.HasSuffix(path, "/..") {
		path += "/"
	}

	size := 0
	parts := strings.Split(path, "/")
	for _, p := range parts {
		switch p {
		case "", ".":
			continue
		case "..":
			size--
			if size < 0 {
				return "/"
			}
		default:
			parts[size] = p
			size++
		}
	}

	if size == 0 {
		return "/"
	}

	if strings.HasSuffix(path, "/") {
		return "/" + strings.Join(parts[0:size], "/") + "/"
	}
	return "/" + strings.Join(parts[0:size], "/")
}

//func (dr *FsDriver) limitConcRun(f func(ch chan interface{}, args ...string) error, args ...string) error {
//	var wg sync.WaitGroup
//	for w := 0; w < dr.maxConc; w++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			sargs := make([]string, len(args))
//			for k, arg := range args {
//				sargs[k] = arg
//			}
//		}()
//	}
//	wg.Wait()
//}
