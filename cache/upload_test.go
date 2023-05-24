package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMutUpload(t *testing.T) {
	mutUpload := &MutUpload{
		UploadID: "1",
		Size:     100 * 12,
		PartSize: 100,
		Path:     "a.jpg",
		UpPath:   "b.jpg",
		CreateAt: time.Now(),
	}
	assert.NoError(t, AddMutUpload(mutUpload))
	assert.NoError(t, AddMutUpload(&MutUpload{
		UploadID: "2",
		Size:     100 * 12,
		PartSize: 100,
		Path:     "/c/a.jpg",
		UpPath:   "b.jpg",
		CreateAt: time.Now(),
	}))
	results, err := FindMutUpload(func(key string, entity *MutUpload) bool {
		return key == mutUpload.Key()
	})

	assert.NoError(t, err)
	assert.Equal(t, len(results), 1)
	assert.Equal(
		t,
		results[0].Key(),
		mutUpload.Key(),
	)
}

func TestMutUploadPart(t *testing.T) {
	part1s := []int64{}
	for i := 0; i < 100; i++ {
		part1s = append(part1s, int64(i))
	}

	for _, v := range part1s {
		err := AddMutUploadPart(&MutUploadPart{
			UploadID: "1",
			PartId:   v,
			Len:      100,
		})
		assert.NoError(t, err)
	}

	part2s := []int64{}
	records, err := FindMutUploadPart(func(key string, entity *MutUploadPart) bool {
		return entity.UploadID == "1"
	})
	assert.NoError(t, err)
	for _, v := range records {
		part2s = append(part2s, v.PartId)
	}

	assert.ElementsMatch(t, part1s, part2s)
}
