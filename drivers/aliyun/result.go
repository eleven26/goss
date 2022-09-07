package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/eleven26/goss/core"
)

type ListObjectResult struct {
	ossResult oss.ListObjectsResult
}

func (l *ListObjectResult) Len() int {
	return len(l.ossResult.Objects)
}

func (l *ListObjectResult) IsFinished() bool {
	return !l.ossResult.IsTruncated
}

func (l *ListObjectResult) Get(index int) core.File {
	return &File{
		properties: l.ossResult.Objects[index],
	}
}
