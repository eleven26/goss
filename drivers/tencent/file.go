package tencent

import (
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type File struct {
	object cos.Object
}

func (f *File) Key() string {
	return f.object.Key
}

func (f *File) Type() string {
	return ""
}

func (f *File) Size() int64 {
	return f.object.Size
}

func (f *File) ETag() string {
	return f.object.ETag
}

func (f *File) LastModified() time.Time {
	t, _ := time.Parse(time.RFC3339, f.object.LastModified)

	return t
}
