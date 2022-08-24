package core

import "time"

// File A slice of File is returned when all files are fetched through Storage.Files.
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
