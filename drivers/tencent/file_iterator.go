package tencent

import (
	"time"

	"github.com/eleven26/goss/core"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type nextChunk func(opt *cos.BucketGetOptions) (*cos.BucketGetResult, *cos.Response, error)

type fileIterator struct {
	dir string

	opt    *cos.BucketGetOptions
	result *cos.BucketGetResult

	index      int
	count      int
	isFinished bool

	nextChunk nextChunk
}

func newFileIterator(dir string) fileIterator {
	return fileIterator{
		dir: dir,
		opt: &cos.BucketGetOptions{
			Prefix: dir,
		},
	}
}

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

	t, _ := time.Parse(time.RFC3339, object.LastModified)

	file = &File{
		key:          object.Key,
		typ:          "",
		size:         object.Size,
		eTag:         object.ETag,
		lastModified: t,
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
	result, _, err := f.nextChunk(f.opt)
	if err != nil {
		return err
	}

	f.index = 0
	f.count = len(result.Contents)
	f.result = result

	if result.IsTruncated {
		f.opt.Marker = result.NextMarker
		f.opt.Prefix = ""
	} else {
		f.isFinished = true
	}

	return nil
}
