package huawei

import (
	"strings"

	"github.com/eleven26/goss/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type ListObjectResult struct {
	files      []core.File
	isFinished bool
}

func NewListObjectResult(output *obs.ListObjectsOutput) core.ListObjectResult {
	return &ListObjectResult{
		files:      getFiles(output.Contents),
		isFinished: !output.IsTruncated,
	}
}

func getFiles(contents []obs.Content) []core.File {
	var files []core.File

	for _, content := range contents {
		if strings.HasSuffix(content.Key, "/") {
			continue
		}

		files = append(files, &File{content: content})
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
