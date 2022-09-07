package aliyun

import (
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/eleven26/goss/core"
)

type ListObjectResult struct {
	ossResult oss.ListObjectsResult
	files     []core.File
}

func NewListObjectResult(ossResult oss.ListObjectsResult) core.ListObjectResult {
	result := ListObjectResult{
		ossResult: ossResult,
	}

	result.files = result.getFiles()

	return &result
}

func (l *ListObjectResult) getFiles() []core.File {
	var files []core.File

	for _, properties := range l.ossResult.Objects {
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
	return !l.ossResult.IsTruncated
}
