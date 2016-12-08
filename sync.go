package main

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"os/signal"
	"path"
	"path/filepath"
)

var (
	db        *leveldb.DB
	maxWorker = 10
)

const (
	EXISTS = iota
	SUCC
	UPLOADFAIL
	LISTFAIL
)

type dbKey struct {
	SrcPath string `json:"src_path"`
	DesPath string `json:"des_path"`
}

type dbValue struct {
	Mtime int64 `json:"modify_time"`
}

type task struct {
	srcPath string
	desPath string
	err     error
	code    int
}

func (t *task) String() string {
	switch t.code {
	case SUCC:
		return fmt.Sprintf("%s to %s OK", t.srcPath, t.desPath)
	case UPLOADFAIL, LISTFAIL:
		return fmt.Sprintf("%s to %s %v fail", t.srcPath, t.desPath, t.err)
	case EXISTS:
		return fmt.Sprintf("%s to %s existed", t.srcPath, t.desPath)
	default:
		return fmt.Sprintf("%s to %s unkown", t.srcPath, t.desPath)
	}
}

func makeKey(src, des string) ([]byte, error) {
	x := dbKey{src, path.Join(user.Bucket, des)}
	return json.Marshal(&x)
}

func makeValue(filename string) (*dbValue, error) {
	info, err := os.Lstat(filename)
	if err != nil {
		return nil, err
	}
	return &dbValue{info.ModTime().UnixNano()}, nil
}

func getValue(src, des string) (*dbValue, error) {
	key, err := makeKey(src, des)
	if err != nil {
		return nil, err
	}

	raw, err := db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	var value dbValue
	if err = json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}

	return &value, nil
}

func setValue(src, des string, v *dbValue) error {
	key, err := makeKey(src, des)
	if err != nil {
		return err
	}

	if v == nil {
		v, err = makeValue(src)
		if err != nil {
			return err
		}
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return db.Put(key, b, nil)
}

func doIterDir(srcPath, desPath string, fiChannel chan *dbKey, stChannel chan *task) {
	filepath.Walk(srcPath, func(_path string, info os.FileInfo, err error) error {
		if err != nil {
			stChannel <- &task{_path, "", err, LISTFAIL}
			return filepath.SkipDir
		}

		if _path == srcPath {
			return nil
		}
		relPath, err := filepath.Rel(srcPath, _path)
		if err != nil {
			stChannel <- &task{_path, "", err, LISTFAIL}
			return filepath.SkipDir
		}
		dokey := &dbKey{
			SrcPath: _path,
			DesPath: path.Join(desPath, filepath.ToSlash(relPath)),
		}
		fiChannel <- dokey

		return nil
	})
	close(fiChannel)
}

func doUploadFile(fiChannel chan *dbKey, stChannel chan *task) {
	for {
		fiValue, more := <-fiChannel
		if !more {
			return
		}

		diskV, err := makeValue(fiValue.SrcPath)
		if err != nil {
			stChannel <- &task{fiValue.SrcPath, fiValue.DesPath, err, UPLOADFAIL}
			continue
		}

		dbV, err := getValue(fiValue.SrcPath, fiValue.DesPath)
		if err != nil {
			stChannel <- &task{fiValue.SrcPath, fiValue.DesPath, err, UPLOADFAIL}
			continue
		}

		if dbV != nil && dbV.Mtime == diskV.Mtime {
			stChannel <- &task{fiValue.SrcPath, fiValue.DesPath, nil, EXISTS}
			continue
		}

		fi, _ := os.Lstat(fiValue.SrcPath)
		if fi.IsDir() {
			err = driver.MakeDir(fiValue.DesPath)
		} else {
			err = driver.uploadFile(fiValue.SrcPath, fiValue.DesPath)
		}
		if err != nil {
			stChannel <- &task{fiValue.SrcPath, fiValue.DesPath, err, UPLOADFAIL}
			continue
		}

		err = setValue(fiValue.SrcPath, fiValue.DesPath, diskV)
		if err != nil {
			stChannel <- &task{fiValue.SrcPath, fiValue.DesPath, err, UPLOADFAIL}
		} else {
			stChannel <- &task{fiValue.SrcPath, fiValue.DesPath, nil, SUCC}
		}
	}
}

func doSync(diskSrc, upDes string) {
	fiChannel := make(chan *dbKey, 2*maxWorker)
	stChannel := make(chan *task, 2*maxWorker)
	doneChan := make(chan int, 2*maxWorker)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	if db == nil {
		var err error
		db, err = leveldb.OpenFile(dbname, nil)
		if err != nil {
			LogC("open file: %v\n", err)
		}
	}

	go doIterDir(diskSrc, upDes, fiChannel, stChannel)
	for i := 0; i < maxWorker; i++ {
		go func() {
			doUploadFile(fiChannel, stChannel)
			doneChan <- 1
		}()
	}
	succ, fails, exists, worker := 0, 0, 0, 0
	for {
		select {
		case <-sigChan:
			LogC("\n%d succ, %d fails, %d ignore.\n", succ, fails, exists)
			return
		case t, more := <-stChannel:
			if !more {
				if fails == 0 {
					LogI("%d succ, %d fails, %d ignore.\n", succ, fails, exists)
				} else {
					LogC("%d succ, %d fails, %d ignore.\n", succ, fails, exists)
				}
				return
			}
			switch t.code {
			case SUCC:
				succ++
				LogD(t.String())
			case LISTFAIL, UPLOADFAIL:
				LogE(t.String())
				fails++
			case EXISTS:
				exists++
				LogD(t.String())
			}
		case <-doneChan:
			worker++
			if worker == maxWorker {
				close(stChannel)
			}
		}
	}
}
