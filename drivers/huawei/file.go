package huawei

import "time"

type File struct {
	key          string
	typ          string
	size         int64
	eTag         string
	lastModified time.Time
}

func (f *File) Key() string {
	return f.key
}

func (f *File) Type() string {
	return f.typ
}

func (f *File) Size() int64 {
	return f.size
}

func (f *File) ETag() string {
	return f.eTag
}

func (f *File) LastModified() time.Time {
	return f.lastModified
}
