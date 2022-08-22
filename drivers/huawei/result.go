package huawei

import (
	"github.com/eleven26/goss/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type ListObjectResult struct {
	output *obs.ListObjectsOutput
}

func (l *ListObjectResult) Len() int {
	return len(l.output.Contents)
}

func (l *ListObjectResult) IsTruncated() bool {
	return l.output.IsTruncated
}

func (l *ListObjectResult) NextMarker() interface{} {
	return l.output.NextMarker
}

func (l *ListObjectResult) Get(index int) core.File {
	return &File{
		content: l.output.Contents[index],
	}
}
