// +build linux darwin

package updriver

import (
	"errors"
	"fmt"
	"github.com/polym/go-sdk/upyun"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type FsDriver struct {
	curDir   string
	operator string
	bucket   string
	maxConc  int
	up       *upyun.UpYun
	logger   *log.Logger
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
	_, err := driver.up.Usage()
	if err != nil {
		return nil, errors.New("username or password is wrong")
	}

	return driver, nil
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

func (dr *FsDriver) GetItems(src, des string) error {
	var ch chan *upyun.FileInfo

	src = dr.abs(src)
	if ok, err := dr.isDir(src); err != nil {
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

					var err error
					if err := os.MkdirAll(path.Dir(desPath), os.ModePerm); err == nil {
						// open a new file
						if info, err := os.Stat(desPath); err == nil {
							if upi, err := dr.up.GetInfo(srcPath); err == nil && upi.Size == info.Size() {
								dr.logger.Printf("GET %s %s EXIST", srcPath, desPath)
								continue
							}
						}
						var fd *os.File
						if fd, err = os.OpenFile(desPath,
							os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600); err == nil {
							err = dr.up.Get(srcPath, fd)
							fd.Close()
						}
					}
					if err != nil {
						// TODO: error
						dr.logger.Printf("GET %s FAIL %v", srcPath, err)
					} else {
						dr.logger.Printf("GET %s %s OK", srcPath, desPath)
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
	for k := 0; k < 200; k++ {
		info, more := <-ch
		if !more {
			return infos[0:k], nil
		}
		infos = append(infos, info)
	}

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

func (dr *FsDriver) PutItems(src, des string) error {
	var wg sync.WaitGroup
	ch := make(chan string, dr.maxConc+10)
	des = dr.abs(des)

	isUpDir, err := dr.isDir(des)
	if err != nil && err.Error() != "X-Error-Code=40400001" {
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
				fname, more := <-ch
				if !more {
					return
				}

				var err error
				var desPath string

				desPath = des
				if isUpDir {
					desPath += "/" + strings.Replace(fname, src, "", 1)
				}

				if strings.HasSuffix(desPath, "/") {
					panic(fmt.Sprintf("desPath should not HasSuffix with / %s", desPath))
				}

				var fd *os.File
				if fd, err = os.OpenFile(fname, os.O_RDWR, 0600); err == nil {
					// check whether is already existed
					fdInfo, _ := fd.Stat()
					if info, err := dr.up.GetInfo(desPath); err == nil && info.Size == fdInfo.Size() {
						dr.logger.Printf("PUT %s -> %s EXIST", fname, desPath)
						fd.Close()
						continue
					} else {
						_, err = dr.up.Put(desPath, fd, false, "", "", nil)
						fd.Close()
					}
				}

				if err != nil {
					// TODO error
					dr.logger.Printf("PUT %s %v FAIL", fname, err)
				} else {
					dr.logger.Printf("PUT %s -> %s OK", fname, desPath)
				}
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
	ok, err := dr.isDir(path)
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
					// TODO: retry
					if err := dr.up.Delete(path + "/" + info.Name); err != nil {
						//TODO: error
						dr.logger.Printf("DELETE %s FAIL %v", dr.abs(path+"/"+info.Name), err)
					} else {
						dr.logger.Printf("DELETE %s OK", dr.abs(path+"/"+info.Name))
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

func (dr *FsDriver) isDir(path string) (bool, error) {
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
