package s3

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/eleven26/goss/core"
)

func NewListObjectResult(entries []*s3.Object, hasNext bool) *core.ListObjectResult {
	return &core.ListObjectResult{
		Files:      getFiles(entries),
		IsFinished: !hasNext,
	}
}

func getFiles(objects []*s3.Object) []core.File {
	var files []core.File

	for _, item := range objects {
		if strings.HasSuffix(*item.Key, "/") {
			continue
		}

		files = append(files, &File{item: item})
	}

	return files
}
