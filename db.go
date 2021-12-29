package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

type dbKey struct {
	SrcPath string `json:"src_path"`
	DstPath string `json:"dst_path"`
}

type fileMeta struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isdir"`
}

type dbValue struct {
	ModifyTime int64       `json:"modify_time"`
	Md5        string      `json:"md5"`
	IsDir      string      `json:"isdir"`
	Items      []*fileMeta `json:"items"`
}

func getDBName() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("USERPROFILE"), ".upx.db")
	}
	return filepath.Join(os.Getenv("HOME"), ".upx.db")
}

func makeDBKey(src, dst string) ([]byte, error) {
	return json.Marshal(&dbKey{
		SrcPath: src,
		DstPath: path.Join(session.Bucket, dst),
	})
}

func makeDBValue(filename string, md5 bool) (*dbValue, error) {
	finfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	dbV := &dbValue{
		ModifyTime: finfo.ModTime().UnixNano(),
	}

	if !finfo.IsDir() {
		if md5 {
			md5Str, _ := md5File(filename)
			dbV.Md5 = md5Str
		}
		dbV.IsDir = "false"
	} else {
		dbV.IsDir = "true"
	}
	return dbV, nil
}

func getDBValue(src, dst string) (*dbValue, error) {
	key, err := makeDBKey(src, dst)
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

func setDBValue(src, dst string, v *dbValue) error {
	key, err := makeDBKey(src, dst)
	if err != nil {
		return err
	}

	if v == nil {
		v, err = makeDBValue(src, true)
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

func delDBValue(src, dst string) error {
	key, err := makeDBKey(src, dst)
	if err != nil {
		return err
	}

	return db.Delete(key, nil)
}

func delDBValues(srcPrefix, dstPrefix string) {
	dstPrefix = path.Join(session.Bucket, dstPrefix)
	iter := db.NewIterator(nil, nil)
	if ok := iter.First(); !ok {
		return
	}
	for {
		k := new(dbKey)
		key := iter.Key()
		err := json.Unmarshal(key, k)
		if err != nil {
			PrintError("decode %s: %v", string(key), err)
		}
		if strings.HasPrefix(k.SrcPath, srcPrefix) && strings.HasPrefix(k.DstPath, dstPrefix) {
			PrintOnlyVerbose("found %s => %s to delete", k.SrcPath, k.DstPath)
			db.Delete(iter.Key(), nil)
		}
		if ok := iter.Next(); !ok {
			break
		}
	}
}

func makeFileMetas(dirname string) ([]*fileMeta, error) {
	var res []*fileMeta
	fInfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return res, err
	}
	for _, fInfo := range fInfos {
		fpath := filepath.Join(dirname, fInfo.Name())
		fi, _ := os.Stat(fpath)
		if fi != nil && fi.IsDir() {
			res = append(res, &fileMeta{fInfo.Name(), true})
		} else {
			res = append(res, &fileMeta{fInfo.Name(), false})
		}
	}
	return res, nil
}

func diffFileMetas(src []*fileMeta, dst []*fileMeta) []*fileMeta {
	i, j := 0, 0
	var res []*fileMeta
	for i < len(src) && j < len(dst) {
		if src[i].Name < dst[j].Name {
			res = append(res, src[i])
			i++
		} else if src[i].Name == dst[j].Name {
			if src[i].IsDir != dst[j].IsDir {
				res = append(res, src[i])
			}
			i++
			j++
		} else {
			j++
		}
	}

	res = append(res, src[i:]...)
	return res
}

func initDB() (err error) {
	db, err = leveldb.OpenFile(getDBName(), nil)
	if err != nil {
		Print("db %v %s", err, getDBName())
	}
	return err
}
