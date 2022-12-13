package minio

import (
	"time"

	"github.com/minio/minio-go/v7"
)

type File struct {
	obj minio.ObjectInfo
}

func (f *File) Key() string {
	return f.obj.Key
}

func (f *File) Type() string {
	return f.obj.ContentType
}

func (f *File) Size() int64 {
	return f.obj.Size
}

func (f *File) ETag() string {
	return f.obj.ETag
}

func (f *File) LastModified() time.Time {
	return f.obj.LastModified
}
