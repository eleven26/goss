package tencent

import (
	"strings"

	"github.com/eleven26/goss/core"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type ListObjectResult struct {
	result *cos.BucketGetResult
	files  []core.File
}

func NewListObjectResult(r *cos.BucketGetResult) core.ListObjectResult {
	result := ListObjectResult{
		result: r,
	}

	result.files = result.getFiles()

	return &result
}

func (l *ListObjectResult) getFiles() []core.File {
	var files []core.File

	for _, content := range l.result.Contents {
		if strings.HasSuffix(content.Key, "/") {
			continue
		}

		files = append(files, &File{object: content})
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
	return !l.result.IsTruncated
}
