package huawei

import (
	"github.com/eleven26/goss/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type nextChunk func(marker string) (*obs.ListObjectsOutput, error)

type fileIterator struct {
	// 前缀
	dir string
	// 获取下一批数据的函数
	nextChunk nextChunk

	// 获取下一批数据的 marker
	marker string
	// ListObjects 的返回值
	result *obs.ListObjectsOutput

	// 当前获取的结果遍历到的下标
	index int
	// 当前已获取的结果的总条数
	count int
	// 是否已经获取完全部数据
	isFinished bool
}

func newFileIterator(dir string, chunk nextChunk) fileIterator {
	return fileIterator{
		dir:        dir,
		marker:     dir,
		isFinished: false,
		nextChunk:  chunk,
	}
}

// HasNext 还没有获取完全部数据、没有遍历完所有获取到的数据
func (f *fileIterator) HasNext() bool {
	return !f.isFinished || f.index < f.count
}

func (f *fileIterator) Next() (file core.File, err error) {
	if !f.HasNext() {
		return
	}

	if f.index >= f.count {
		err = f.getNextChunk()
		if err != nil {
			return
		}
	}

	object := f.result.Contents[f.index]
	f.index++

	file = &File{
		key:          object.Key,
		typ:          "",
		size:         object.Size,
		eTag:         object.ETag,
		lastModified: object.LastModified,
	}

	return file, nil
}

func (f *fileIterator) All() ([]core.File, error) {
	var res []core.File

	for {
		if f.HasNext() {
			file, err := f.Next()
			if err != nil {
				return res, err
			}

			res = append(res, file)
		} else {
			break
		}
	}

	return res, nil
}

func (f *fileIterator) getNextChunk() error {
	result, err := f.nextChunk(f.marker)
	if err != nil {
		return err
	}

	f.index = 0
	f.count = len(result.Contents)
	f.result = result

	if result.IsTruncated {
		f.marker = result.NextMarker
	} else {
		f.isFinished = true
	}

	return nil
}
