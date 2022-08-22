package aliyun

import (
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type File struct {
	properties oss.ObjectProperties
}

func (f *File) Key() string {
	return f.properties.Key
}

func (f *File) Type() string {
	return f.properties.Type
}

func (f *File) Size() int64 {
	return f.properties.Size
}

func (f *File) ETag() string {
	return f.properties.ETag
}

func (f *File) LastModified() time.Time {
	return f.properties.LastModified
}
