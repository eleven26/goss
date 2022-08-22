package huawei

import (
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type File struct {
	content obs.Content
}

func (f *File) Key() string {
	return f.content.Key
}

func (f *File) Type() string {
	return ""
}

func (f *File) Size() int64 {
	return f.content.Size
}

func (f *File) ETag() string {
	return f.content.ETag
}

func (f *File) LastModified() time.Time {
	return f.content.LastModified
}
