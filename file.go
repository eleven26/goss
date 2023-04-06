package goss

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// File A slice of File is returned when all files are fetched through Store.Files.
// File hides differences in object returns from different cloud storage providers,
// but at the same time, it can only return less information.
type File interface {
	// Key is the object key.
	Key() string

	// Type is the object type.
	Type() string

	// Size is the size of current file.
	Size() int64

	// ETag is etag return by cloud storage providers.
	ETag() string

	// LastModified is the time when the object was last updated.
	LastModified() time.Time
}

type file struct {
	item types.Object
}

func (f *file) Key() string {
	return *f.item.Key
}

func (f *file) Type() string {
	return string(f.item.StorageClass)
}

func (f *file) Size() int64 {
	return f.item.Size
}

func (f *file) ETag() string {
	return *f.item.ETag
}

func (f *file) LastModified() time.Time {
	return *f.item.LastModified
}
