package s3

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

type File struct {
	item *s3.Object
}

func (f *File) Key() string {
	return *f.item.Key
}

func (f *File) Type() string {
	return *f.item.StorageClass
}

func (f *File) Size() int64 {
	return *f.item.Size
}

func (f *File) ETag() string {
	return *f.item.ETag
}

func (f *File) LastModified() time.Time {
	return *f.item.LastModified
}
