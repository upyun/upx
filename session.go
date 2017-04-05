package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"github.com/jehiah/go-strftime"
	"github.com/upyun/go-sdk/upyun"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	EXISTS = iota
	SUCC
	FAIL
	DELETE
)

type Session struct {
	Bucket   string `json:"bucket"`
	Operator string `json:"username"`
	Password string `json:"password"`
	CWD      string `json:"cwd"`

	updriver *upyun.UpYun
	color    bool
	stats    map[string]int
	smu      sync.RWMutex
}

var (
	session *Session
)

func (sess *Session) AbsPath(upPath string) (ret string) {
	if strings.HasPrefix(upPath, "/") {
		ret = path.Join(upPath)
	} else {
		ret = path.Join(sess.CWD, upPath)
	}

	if strings.HasSuffix(upPath, "/") && ret != "/" {
		ret += "/"
	}
	return
}

func (sess *Session) IsUpYunDir(upPath string) (isDir bool, exist bool) {
	upInfo, err := sess.updriver.GetInfo(sess.AbsPath(upPath))
	if err != nil {
		return false, false
	}
	return upInfo.IsDir, true
}

func (sess *Session) IsLocalDir(localPath string) (isDir bool, exist bool) {
	fInfo, err := os.Stat(localPath)
	if err != nil {
		return false, false
	}
	return fInfo.IsDir(), true
}

func (sess *Session) FormatUpInfo(upInfo *upyun.FileInfo) string {
	s := "drwxrwxrwx"
	if !upInfo.IsDir {
		s = "-rw-rw-rw-"
	}
	s += fmt.Sprintf(" 1 %s %s %12d", sess.Operator, sess.Bucket, upInfo.Size)
	if upInfo.Time.Year() != time.Now().Year() {
		s += " " + strftime.Format("%b %d  %Y", upInfo.Time)
	} else {
		s += " " + strftime.Format("%b %d %H:%M", upInfo.Time)
	}
	if upInfo.IsDir && sess.color {
		s += " " + color.BlueString(upInfo.Name)
	} else {
		s += " " + upInfo.Name
	}
	return s
}

func (sess *Session) Init() error {
	sess.updriver = upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:    sess.Bucket,
		Operator:  sess.Operator,
		Password:  sess.Password,
		UserAgent: fmt.Sprintf("upx/%s", VERSION),
	})
	_, err := sess.updriver.Usage()
	return err
}

func (sess *Session) Info() {
	n, err := sess.updriver.Usage()
	if err != nil {
		PrintErrorAndExit("usage: %v", err)
	}

	tmp := []string{
		fmt.Sprintf("ServiceName:   %s", sess.Bucket),
		fmt.Sprintf("Operator:      %s", sess.Operator),
		fmt.Sprintf("CurrentDir:    %s", sess.CWD),
		fmt.Sprintf("Usage:         %s", humanizeSize(n)),
	}

	Print(strings.Join(tmp, "\n"))
}

func (sess *Session) Pwd() {
	Print("%s", sess.CWD)
}

func (sess *Session) Mkdir(upPaths ...string) {
	for _, upPath := range upPaths {
		fpath := sess.AbsPath(upPath)
		if err := sess.updriver.Mkdir(fpath); err != nil {
			PrintErrorAndExit("mkdir %s: %v", fpath, err)
		}
	}
}

func (sess *Session) Cd(upPath string) {
	fpath := sess.AbsPath(upPath)
	if isDir, _ := sess.IsUpYunDir(fpath); isDir {
		sess.CWD = fpath
		Print(sess.CWD)
	} else {
		PrintErrorAndExit("cd: %s: Not a directory", fpath)
	}
}

func (sess *Session) Ls(upPath string, match *MatchConfig, maxItems int, isDesc bool) {
	fpath := sess.AbsPath(upPath)
	isDir, exist := sess.IsUpYunDir(fpath)
	if !exist {
		PrintErrorAndExit("ls: cannot access %s: No such file or directory", fpath)
	}

	if !isDir {
		fInfo, err := sess.updriver.GetInfo(fpath)
		if err != nil {
			PrintErrorAndExit("ls %s: %v", fpath, err)
		}
		if IsMatched(fInfo, match) {
			Print(sess.FormatUpInfo(fInfo))
		} else {
			PrintErrorAndExit("ls: cannot access %s: No such file or directory", fpath)
		}
		return
	}

	fInfoChan := make(chan *upyun.FileInfo, 50)
	go func() {
		err := sess.updriver.List(&upyun.GetObjectsConfig{
			Path:        fpath,
			ObjectsChan: fInfoChan,
			DescOrder:   isDesc,
		})
		if err != nil {
			PrintErrorAndExit("ls %s: %v", fpath, err)
		}
	}()

	objs := 0
	for fInfo := range fInfoChan {
		if IsMatched(fInfo, match) {
			Print(sess.FormatUpInfo(fInfo))
			objs++
		}
		if maxItems > 0 && objs >= maxItems {
			break
		}
	}
	if objs == 0 {
		msg := fpath
		if match.Wildcard != "" {
			msg = fpath + "/" + match.Wildcard
		}
		if match.TimeType != TIME_NOT_SET {
			msg += " timestamp@"
			if match.TimeType == TIME_AFTER || match.TimeType == TIME_INTERVAL {
				msg += "[" + match.After.Format("2006-01-02 15:04:05") + ","
			} else {
				msg += "[-oo,"
			}
			if match.TimeType == TIME_BEFORE || match.TimeType == TIME_INTERVAL {
				msg += match.Before.Format("2006-01-02 15:04:05") + "]"
			} else {
				msg += "+oo]"
			}
		}
		PrintErrorAndExit("ls: cannot access %s: No such file or directory", msg)
	}
}

func (sess *Session) getDir(upPath, localPath string, match *MatchConfig, workers int) error {
	if err := os.MkdirAll(localPath, 0755); err != nil {
		return err
	}

	var wg sync.WaitGroup

	fInfoChan := make(chan *upyun.FileInfo, workers*2)
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			id := -1
			var e error
			for fInfo := range fInfoChan {
				if IsMatched(fInfo, match) {
					fpath := path.Join(upPath, fInfo.Name)
					lpath := filepath.Join(localPath, filepath.FromSlash(fInfo.Name))
					if fInfo.IsDir {
						os.MkdirAll(lpath, 0755)
					} else {
						id, e = sess.getFileWithProgress(id, fpath, lpath, fInfo)
					}
					if e != nil {
						return
					}
				}
			}
		}()
	}

	err := sess.updriver.List(&upyun.GetObjectsConfig{
		Path:         upPath,
		ObjectsChan:  fInfoChan,
		MaxListTries: 3,
		MaxListLevel: -1,
	})
	wg.Wait()
	return err
}

func (sess *Session) getFileWithProgress(id int, upPath, localPath string, upInfo *upyun.FileInfo) (int, error) {
	var err error
	bar, idx := AddBar(id, int(upInfo.Size))
	bar = bar.AppendCompleted()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		status := "WAIT"
		if b.Current() == b.Total {
			status = "OK"
		}
		name := leftAlign(shortPath(localPath, 40), 40)
		if err != nil {
			b.Set(bar.Total)
			return fmt.Sprintf("%s ERR %s", name, err)
		}
		return fmt.Sprintf("%s %s", name, rightAlign(status, 4))
	})

	dir := filepath.Dir(localPath)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return id, err
	}

	w, err := NewFileWrappedWriter(localPath)
	if err != nil {
		return id, err
	}
	defer w.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err == nil {
			if w.Copyed == bar.Total {
				bar.Set(w.Copyed)
				break
			}
			bar.Set(w.Copyed)
		}
	}()

	_, err = sess.updriver.Get(&upyun.GetObjectConfig{
		Path:   sess.AbsPath(upPath),
		Writer: w,
	})
	wg.Wait()
	return idx, err
}

func (sess *Session) Get(upPath, localPath string, match *MatchConfig, workers int) {
	upPath = sess.AbsPath(upPath)
	upInfo, err := sess.updriver.GetInfo(upPath)
	if err != nil {
		PrintErrorAndExit("getinfo %s: %v", upPath, err)
	}

	exist, isDir := false, false
	if localInfo, _ := os.Stat(localPath); localInfo != nil {
		exist = true
		isDir = localInfo.IsDir()
	} else {
		if strings.HasSuffix(localPath, "/") {
			isDir = true
		}
	}

	if upInfo.IsDir {
		if exist {
			if !isDir {
				PrintErrorAndExit("get: %s Not a directory", localPath)
			} else {
				if match.Wildcard == "" {
					localPath = filepath.Join(localPath, path.Base(upPath))
				}
			}
		}
		sess.getDir(upPath, localPath, match, workers)
	} else {
		if isDir {
			localPath = filepath.Join(localPath, path.Base(upPath))
		}
		sess.getFileWithProgress(-1, upPath, localPath, upInfo)
	}
}

func (sess *Session) putFileWithProgress(barId int, localPath, upPath string, localInfo os.FileInfo) (int, error) {
	var err error
	bar, idx := AddBar(barId, int(localInfo.Size()))
	bar = bar.AppendCompleted()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		status := "WAIT"
		if b.Current() == b.Total {
			status = "OK"
		}
		name := leftAlign(shortPath(upPath, 40), 40)
		if err != nil {
			b.Set(bar.Total)
			return fmt.Sprintf("%s ERR %s", name, err)
		}
		return fmt.Sprintf("%s %s", name, rightAlign(status, 4))
	})

	fd, err := os.Open(localPath)
	if err != nil {
		return idx, err
	}

	var wg sync.WaitGroup
	wReader := &ProgressReader{fd: fd}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err == nil {
			if wReader.Copyed() == bar.Total {
				bar.Set(wReader.Copyed())
				return
			}
			bar.Set(wReader.Copyed())
		}
	}()

	err = sess.updriver.Put(&upyun.PutObjectConfig{
		Path: upPath,
		Headers: map[string]string{
			"Content-Length": fmt.Sprint(localInfo.Size()),
		},
		Reader: wReader,
	})
	wg.Wait()
	return idx, err
}

func (sess *Session) putDir(localPath, upPath string, workers int) {
	type FileInfo struct {
		fpath string
		fInfo os.FileInfo
	}

	localFiles := make(chan *FileInfo, workers*2)
	var wg sync.WaitGroup
	var err error
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			barId := -1
			for info := range localFiles {
				rel, _ := filepath.Rel(localPath, info.fpath)
				desPath := path.Join(upPath, filepath.ToSlash(rel))
				if fInfo, err := os.Stat(info.fpath); err == nil && fInfo.IsDir() {
					err = sess.updriver.Mkdir(desPath)
				} else {
					barId, err = sess.putFileWithProgress(barId, info.fpath, desPath, info.fInfo)
				}
				if err != nil {
					return
				}
			}
		}()
	}

	walk(localPath, func(fpath string, fInfo os.FileInfo, err error) {
		if err == nil {
			localFiles <- &FileInfo{
				fpath: fpath,
				fInfo: fInfo,
			}
		}
	})

	close(localFiles)
	wg.Wait()
}

func (sess *Session) Put(localPath, upPath string, workers int) {
	upPath = sess.AbsPath(upPath)
	localInfo, err := os.Stat(localPath)
	if err != nil {
		PrintErrorAndExit("stat %s: %v", localPath, err)
	}

	exist, isDir := false, false
	if upInfo, _ := sess.updriver.GetInfo(upPath); upInfo != nil {
		exist = true
		isDir = upInfo.IsDir
	} else {
		if strings.HasSuffix(upPath, "/") {
			isDir = true
		}
	}
	if localInfo.IsDir() {
		if exist {
			if !isDir {
				PrintErrorAndExit("put: %s: Not a directory", upPath)
			} else {
				upPath = path.Join(upPath, filepath.Base(localPath))
			}
		}
		sess.putDir(localPath, upPath, workers)
	} else {
		if isDir {
			upPath = path.Join(upPath, filepath.Base(localPath))
		}
		sess.putFileWithProgress(-1, localPath, upPath, localInfo)
	}
}

func (sess *Session) rmFile(fpath string, isAsync bool) {
	if sess.stats == nil {
		sess.stats = make(map[string]int)
	}
	err := sess.updriver.Delete(&upyun.DeleteObjectConfig{
		Path:  fpath,
		Async: isAsync,
	})
	if err == nil {
		sess.smu.Lock()
		sess.stats["ok"]++
		sess.smu.Unlock()
		PrintOnlyVerbose("DELETE %s OK", fpath)
	} else {
		sess.smu.Lock()
		sess.stats["fail"]++
		sess.smu.Unlock()
		PrintError("DELETE %s FAIL %v", fpath, err)
	}
}

func (sess *Session) rmEmptyDir(fpath string, isAsync bool) {
	sess.rmFile(fpath, isAsync)
}

func (sess *Session) rmDir(fpath string, isAsync bool) {
	fInfoChan := make(chan *upyun.FileInfo, 50)
	go func() {
		err := sess.updriver.List(&upyun.GetObjectsConfig{
			Path:        fpath,
			ObjectsChan: fInfoChan,
		})
		if err != nil {
			PrintErrorAndExit("ls %s: %v", fpath, err)
		}
	}()

	for fInfo := range fInfoChan {
		fp := path.Join(fpath, fInfo.Name)
		if fInfo.IsDir {
			sess.rmDir(fp, isAsync)
		} else {
			sess.rmFile(fp, isAsync)
		}
	}
	sess.rmEmptyDir(fpath, isAsync)
}

func (sess *Session) Rm(upPath string, match *MatchConfig, isAsync bool) {
	fpath := sess.AbsPath(upPath)
	isDir, exist := sess.IsUpYunDir(fpath)
	if !exist {
		PrintErrorAndExit("rm: cannot remove %s: No such file or directory", fpath)
	}

	if isDir && match != nil && match.Wildcard == "" {
		if match.ItemType == FILE {
			PrintErrorAndExit("rm: cannot remove %s: Is a directory, add -d/-a flag", fpath)
		}
		sess.rmDir(fpath, isAsync)
		return
	}

	if !isDir {
		fInfo, err := sess.updriver.GetInfo(fpath)
		if err != nil {
			PrintErrorAndExit("getinfo %s: %v", fpath, err)
		}
		if IsMatched(fInfo, match) {
			sess.rmFile(fpath, isAsync)
		}
		return
	}

	fInfoChan := make(chan *upyun.FileInfo, 50)
	go func() {
		err := sess.updriver.List(&upyun.GetObjectsConfig{
			Path:        fpath,
			ObjectsChan: fInfoChan,
		})
		if err != nil {
			PrintErrorAndExit("ls %s: %v", fpath, err)
		}
	}()

	for fInfo := range fInfoChan {
		fp := path.Join(fpath, fInfo.Name)
		if IsMatched(fInfo, match) {
			if fInfo.IsDir {
				sess.rmDir(fp, isAsync)
			} else {
				sess.rmFile(fp, isAsync)
			}
		}
	}
}

func (sess *Session) tree(upPath, prefix string, output chan string) (folders, files int, err error) {
	upInfos := make(chan *upyun.FileInfo, 50)
	fpath := sess.AbsPath(upPath)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		prevInfo := <-upInfos
		for fInfo := range upInfos {
			p := prefix + "|-- "
			if prevInfo.IsDir {
				if sess.color {
					output <- p + color.BlueString("%s", prevInfo.Name)
				} else {
					output <- p + prevInfo.Name
				}
				folders++
				d, f, _ := sess.tree(path.Join(fpath, prevInfo.Name), prefix+"!   ", output)
				folders += d
				files += f
			} else {
				output <- p + prevInfo.Name
				files++
			}
			prevInfo = fInfo
		}
		if prevInfo == nil {
			return
		}
		p := prefix + "`-- "
		if prevInfo.IsDir {
			if sess.color {
				output <- p + color.BlueString("%s", prevInfo.Name)
			} else {
				output <- p + prevInfo.Name
			}
			folders++
			d, f, _ := sess.tree(path.Join(fpath, prevInfo.Name), prefix+"    ", output)
			folders += d
			files += f
		} else {
			output <- p + prevInfo.Name
			files++
		}
	}()

	err = sess.updriver.List(&upyun.GetObjectsConfig{
		Path:        fpath,
		ObjectsChan: upInfos,
	})
	wg.Wait()

	return
}

func (sess *Session) Tree(upPath string) {
	fpath := sess.AbsPath(upPath)
	files, folders := 0, 0
	defer func() {
		Print("\n%d directories, %d files", folders, files)
	}()

	if isDir, _ := sess.IsUpYunDir(fpath); !isDir {
		PrintErrorAndExit("%s [error opening dir]", fpath)
	}
	Print("%s", fpath)

	output := make(chan string, 50)
	go func() {
		folders, files, _ = sess.tree(fpath, "", output)
		close(output)
	}()

	for s := range output {
		Print(s)
	}
	return
}

func (sess *Session) syncOneObject(localPath, upPath string) (status int, err error) {
	upPath = sess.AbsPath(upPath)
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return FAIL, err
	}

	diskV := &dbValue{}
	isDir, _ := sess.IsLocalDir(localPath)
	if isDir {
		err = sess.updriver.Mkdir(upPath)
	} else {
		dbV, err := getDBValue(localPath, upPath)
		if err != nil {
			return FAIL, err
		}

		fInfo, err := os.Stat(localPath)
		if err != nil {
			return FAIL, err
		}

		diskV.ModifyTime = fInfo.ModTime().UnixNano()
		diskV.Md5, _ = md5File(localPath)
		diskV.IsDir = "false"
		if dbV != nil {
			if dbV.ModifyTime == diskV.ModifyTime {
				return EXISTS, nil
			}
			if dbV.Md5 != "" && dbV.Md5 == diskV.Md5 {
				return EXISTS, nil
			}
		}
		err = sess.updriver.Put(&upyun.PutObjectConfig{
			Path:      upPath,
			LocalPath: localPath,
		})
	}

	if err == nil {
		if isDir {
			return SUCC, nil
		}
		if err = setDBValue(localPath, upPath, diskV); err == nil {
			return SUCC, nil
		}
	}
	return FAIL, err
}

func (sess *Session) Sync(localPath, upPath string, workers int, delete bool) {
	type task struct {
		src, dst, typ string
		isdir         bool
	}
	var wg sync.WaitGroup
	tasks := make(chan *task, workers*2)
	stats := make(chan int, workers*2)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	upPath = sess.AbsPath(upPath)
	localPath, _ = filepath.Abs(localPath)

	if err := initDB(); err != nil {
		PrintErrorAndExit("sync: init database: %v", err)
	}

	var dlock sync.Mutex
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				switch task.typ {
				case "sync":
					stat, err := sess.syncOneObject(task.src, task.dst)
					switch stat {
					case SUCC:
						PrintOnlyVerbose("sync %s to %s OK", task.src, task.dst)
					case EXISTS:
						PrintOnlyVerbose("sync %s to %s EXISTS", task.src, task.dst)
					case FAIL:
						PrintError("sync %s to %s FAIL %v", task.src, task.dst, err)
					}
					stats <- stat
				case "delete":
					dlock.Lock()
					if task.isdir {
						sess.rmDir(task.dst, false)
					} else {
						sess.rmFile(task.dst, false)
					}
					dlock.Unlock()
					delDBValue(task.src, task.dst)
					stats <- DELETE
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(stats)
	}()

	go func() {
		walk(localPath, func(fpath string, fInfo os.FileInfo, err error) {
			if err != nil {
				PrintError("walk dir error %s: %v", fpath, err)
				stats <- FAIL
				return
			}

			uppath, err := filepath.Rel(localPath, fpath)
			if err != nil {
				PrintError("get relative path error %s: %v", fpath, err)
				stats <- FAIL
				return
			}
			uppath = path.Join(upPath, uppath)

			if fInfo.IsDir() {
				// check node info
				dbV, err := getDBValue(fpath, uppath)
				if err != nil {
					PrintError("get db error %s %v", fpath, err)
					stats <- FAIL
					return
				}

				// get meta info
				fMetas, err := makeFileMetas(fpath)
				if err != nil {
					PrintError("make file metas error %s: %v", fpath, err)
					stats <- FAIL
					return
				}
				// update info
				curDBV, _ := makeDBValue(fpath)
				curDBV.Items = fMetas
				setDBValue(fpath, uppath, curDBV)

				// no value before
				if dbV == nil {
					tasks <- &task{src: fpath, dst: uppath, typ: "sync"}
				} else {
					if delete {
						for _, meta := range diffFileMetas(dbV.Items, fMetas) {
							tasks <- &task{
								src:   filepath.Join(fpath, meta.Name),
								dst:   path.Join(uppath, meta.Name),
								typ:   "delete",
								isdir: meta.IsDir,
							}
						}
					}

					if dbV.IsDir == "false" {
						tasks <- &task{src: fpath, dst: uppath, typ: "sync"}
					} else {
						PrintOnlyVerbose("sync %s to %s EXISTS", fpath, uppath)
						stats <- EXISTS
					}
				}
			} else {
				tasks <- &task{src: fpath, dst: uppath, typ: "sync"}
			}
		})
		close(tasks)
	}()

	counts := make([]int, 4)
	for {
		select {
		case <-sigChan:
			PrintErrorAndExit("%d succs, %d fails, %d ignores, %d/%d dels.\n",
				counts[SUCC], counts[FAIL], counts[EXISTS], sess.stats["ok"], sess.stats["fail"])
		case t, more := <-stats:
			if !more {
				if counts[FAIL] == 0 && sess.stats["fail"] == 0 {
					Print("%d succs, %d fails, %d ignores, %d/%d dels.\n",
						counts[SUCC], counts[FAIL], counts[EXISTS], sess.stats["ok"], sess.stats["fail"])
				} else {
					PrintErrorAndExit("%d succs, %d fails, %d ignores, %d/%d dels.\n",
						counts[SUCC], counts[FAIL], counts[EXISTS], sess.stats["ok"], sess.stats["fail"])
				}
				return
			}
			counts[t]++
		}
	}
}

func (sess *Session) PostTask(app, notify, taskFile string) {
	fd, err := os.Open(taskFile)
	if err != nil {
		PrintErrorAndExit("open %s: %v", taskFile, err)
	}

	body, err := ioutil.ReadAll(fd)
	fd.Close()
	if err != nil {
		PrintErrorAndExit("read %s: %v", taskFile, err)
	}

	var tasks []interface{}
	if err = json.Unmarshal(body, &tasks); err != nil {
		PrintErrorAndExit("json Unmarshal: %v", err)
	}

	if notify == "" {
		notify = "https://httpbin.org/post"
	}
	ids, err := sess.updriver.CommitTasks(&upyun.CommitTasksConfig{
		AppName:   app,
		NotifyUrl: notify,
		Tasks:     tasks,
	})
	if err != nil {
		PrintErrorAndExit("commit tasks: %v", err)
	}
	Print("%v", ids)
}
