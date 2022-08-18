package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/eleven26/goss/core"
)

type nextChunk func(option oss.Option) (oss.ListObjectsResult, error)

type fileIterator struct {
	// 前缀
	dir string
	// 获取下一批数据的函数
	nextChunk nextChunk

	// 获取下一批数据的 marker
	marker oss.Option
	// ListObjects 的返回值
	result oss.ListObjectsResult

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
		marker:     oss.Prefix(dir),
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

	object := f.result.Objects[f.index]
	f.index++

	file = &File{
		key:          object.Key,
		typ:          object.Type,
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
	f.count = len(result.Objects)
	f.result = result

	if result.IsTruncated {
		f.marker = oss.Marker(result.NextMarker)
	} else {
		f.isFinished = true
	}

	return nil
}
