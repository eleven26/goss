package aliyun

import (
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/eleven26/goss/core"
)

type ListObjectResult struct {
	files      []core.File
	isFinished bool
}

func NewListObjectResult(result oss.ListObjectsResult) core.ListObjectResult {
	return &ListObjectResult{
		files:      getFiles(result.Objects),
		isFinished: !result.IsTruncated,
	}
}

func getFiles(objects []oss.ObjectProperties) []core.File {
	var files []core.File

	for _, properties := range objects {
		if strings.HasSuffix(properties.Key, "/") {
			continue
		}

		files = append(files, &File{properties: properties})
	}

	return files
}

func (l *ListObjectResult) Files() []core.File {
	return l.files
}

func (l *ListObjectResult) Len() int {
	return len(l.files)
}

func (l *ListObjectResult) IsFinished() bool {
	return l.isFinished
}
