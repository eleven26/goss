package core

type ListObjectResult interface {
	Len() int
	IsTruncated() bool
	NextMarker() interface{}
	Get(index int) File
}
