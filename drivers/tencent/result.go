package tencent

import (
	"github.com/eleven26/goss/core"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type ListObjectResult struct {
	result *cos.BucketGetResult
}

func (l *ListObjectResult) Len() int {
	return len(l.result.Contents)
}

func (l *ListObjectResult) IsTruncated() bool {
	return l.result.IsTruncated
}

func (l *ListObjectResult) NextMarker() interface{} {
	return l.result.NextMarker
}

func (l *ListObjectResult) Get(index int) core.File {
	return &File{
		object: l.result.Contents[index],
	}
}
