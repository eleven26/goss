package huawei

import (
	"strings"

	"github.com/eleven26/goss/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type ListObjectResult struct {
	output *obs.ListObjectsOutput
	files  []core.File
}

func NewListObjectResult(output *obs.ListObjectsOutput) core.ListObjectResult {
	result := ListObjectResult{
		output: output,
	}

	result.files = result.getFiles()

	return &result
}

func (l *ListObjectResult) getFiles() []core.File {
	var files []core.File

	for _, content := range l.output.Contents {
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
	return !l.output.IsTruncated
}
