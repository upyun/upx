package main

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var db *leveldb.DB

type dbKey struct {
	SrcPath string `json:"src_path"`
	DstPath string `json:"dst_path"`
}

type dbValue struct {
	ModifyTime int64 `json:"modify_time"`
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

func makeDBValue(filename string) (*dbValue, error) {
	finfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	return &dbValue{finfo.ModTime().UnixNano()}, nil
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
		v, err = makeDBValue(src)
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

func initDB() (err error) {
	db, err = leveldb.OpenFile(getDBName(), nil)
	if err != nil {
		Print("db %v %s", err, getDBName())
	}
	return err
}
