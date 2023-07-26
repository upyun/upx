package cache

import (
	"encoding/json"
	"fmt"
	"time"
)

// 分片上传任务
type MutUpload struct {
	UploadID string

	// 文件总计大小
	Size int64

	// 分块大小
	PartSize int64

	// 本都文件路径
	Path string

	// 云端文件路径
	UpPath string

	// 上传时间
	CreateAt time.Time
}

func (p *MutUpload) Key() string {
	return fmt.Sprintf("mutupload-%s", p.UpPath)
}

// 查询分片上传任务
func FindMutUpload(fn func(key string, entity *MutUpload) bool) ([]*MutUpload, error) {
	var result []*MutUpload
	err := Range("mutupload-", func(key []byte, value []byte) {
		var item = &MutUpload{}
		if err := json.Unmarshal(value, item); err != nil {
			db.Delete(key, nil)
			return
		}

		// 删除过期的分片上传记录
		if time.Since(item.CreateAt).Hours() > 24 {
			FindMutUploadPart(func(key string, part *MutUploadPart) bool {
				if part.UploadID == item.UploadID {
					db.Delete([]byte(key), nil)
				}
				return false
			})
			db.Delete(key, nil)
		}

		if fn(string(key), item) {
			result = append(result, item)
		}
	})
	return result, err
}

// 添加分片上传
func AddMutUpload(entity *MutUpload) error {
	db, err := GetClient()
	if err != nil {
		return err
	}

	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	return db.Put([]byte(entity.Key()), data, nil)
}

// 分片上传任务下的具体分片信息
type MutUploadPart struct {
	UploadID string
	PartId   int64
	Len      int64
}

func (p *MutUploadPart) Key() string {
	return fmt.Sprintf("part-%s-%d", p.UploadID, p.PartId)
}

// 获取已经上传的分片
func FindMutUploadPart(fn func(key string, entity *MutUploadPart) bool) ([]*MutUploadPart, error) {
	var result []*MutUploadPart
	err := Range("part-", func(key []byte, value []byte) {
		var item = &MutUploadPart{}
		if err := json.Unmarshal(value, item); err != nil {
			db.Delete(key, nil)
			return
		}

		if fn(string(key), item) {
			result = append(result, item)
		}
	})
	return result, err
}

// 记录已经上传的分片
func AddMutUploadPart(entity *MutUploadPart) error {
	db, err := GetClient()
	if err != nil {
		return err
	}

	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return db.Put([]byte(entity.Key()), data, nil)
}

func DeleteUpload(upPath, uploadID string) error {
	Range("mutupload-"+upPath, func(key, data []byte) {
		Delete(string(key))
	})
	Range("part-"+uploadID, func(key, data []byte) {
		Delete(string(key))
	})
	return nil
}
