package core

import "time"

type File interface {
	Key() string
	Type() string
	Size() int64
	ETag() string
	LastModified() time.Time
}
