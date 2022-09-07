package tencent

import (
	"strings"

	"github.com/eleven26/goss/core"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type ListObjectResult struct {
	files      []core.File
	isFinished bool
}

func NewListObjectResult(r *cos.BucketGetResult) core.ListObjectResult {
	return &ListObjectResult{
		files:      getFiles(r.Contents),
		isFinished: !r.IsTruncated,
	}
}

func getFiles(contents []cos.Object) []core.File {
	var files []core.File

	for _, content := range contents {
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
	return l.isFinished
}
