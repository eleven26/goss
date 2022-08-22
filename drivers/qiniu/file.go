package qiniu

import (
	"strconv"
	"time"

	"github.com/qiniu/go-sdk/v7/storage"
)

type File struct {
	item storage.ListItem
}

func (f *File) Key() string {
	return f.item.Key
}

func (f *File) Type() string {
	return strconv.Itoa(f.item.Type)
}

func (f *File) Size() int64 {
	return f.item.Fsize
}

func (f *File) ETag() string {
	return f.item.Hash
}

func (f *File) LastModified() time.Time {
	return time.UnixMicro(f.item.PutTime / 10)
}
