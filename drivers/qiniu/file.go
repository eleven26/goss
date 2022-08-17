package qiniu

import (
	"strconv"
	"time"
)

type File struct {
	key     string
	typ     int // storage.fileInfo.type
	fSize   int64
	hash    string
	putTime int64
}

func (f *File) Key() string {
	return f.key
}

func (f *File) Type() string {
	return strconv.Itoa(f.typ)
}

func (f *File) Size() int64 {
	return f.fSize
}

func (f *File) ETag() string {
	return f.hash
}

func (f *File) LastModified() time.Time {
	return time.UnixMicro(f.putTime / 10)
}
