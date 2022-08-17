package qiniu

import (
	"errors"

	"github.com/eleven26/goss/core"
	"github.com/qiniu/go-sdk/v7/storage"
)

type nextChunk func(marker string) (entries []storage.ListItem, commonPrefixes []string, nextMarker string, hasNext bool, err error)

var errorEmpty = errors.New("empty result")

type fileIterator struct {
	dir string

	nextMarker string
	result     []storage.ListItem

	index      int
	count      int
	isFinished bool

	nextChunk nextChunk
}

func newFileIterator(dir string, chunk nextChunk) fileIterator {
	return fileIterator{
		dir:       dir,
		nextChunk: chunk,
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

	if len(f.result) == 0 {
		return nil, errorEmpty
	}

	object := f.result[f.index]
	f.index++

	file = &File{
		key:     object.Key,
		typ:     object.Type,
		fSize:   object.Fsize,
		hash:    object.Hash,
		putTime: object.PutTime,
	}

	return file, nil
}

func (f *fileIterator) All() ([]core.File, error) {
	var res []core.File

	for {
		if f.HasNext() {
			file, err := f.Next()
			if err != nil {
				if errors.Is(err, errorEmpty) {
					return res, nil
				}
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
	entries, _, nextMarker, hasNext, err := f.nextChunk(f.nextMarker)
	if err != nil {
		return err
	}

	f.index = 0
	f.count = len(entries)
	f.result = entries

	if hasNext {
		f.nextMarker = nextMarker
	} else {
		f.isFinished = true
	}

	return nil
}
