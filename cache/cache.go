package cache

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// 存储本地数据库连接
var db *leveldb.DB

func GetDBName() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("USERPROFILE"), ".upx.db")
	}
	return filepath.Join(os.Getenv("HOME"), ".upx.db")
}

func GetClient() (*leveldb.DB, error) {
	var err error
	if db == nil {
		db, err = leveldb.OpenFile(GetDBName(), nil)
	}
	return db, err
}

func Delete(key string) error {
	db, err := GetClient()
	if err != nil {
		return err
	}
	return db.Delete([]byte(key), nil)
}

func Range(scoop string, fn func(key []byte, data []byte)) error {
	db, err := GetClient()
	if err != nil {
		return err
	}

	iter := db.NewIterator(
		util.BytesPrefix([]byte(scoop)),
		nil,
	)

	for iter.Next() {
		fn(iter.Key(), iter.Value())
	}

	iter.Release()
	return iter.Error()
}
