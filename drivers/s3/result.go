package s3

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/eleven26/goss/v2/core"
)

func NewListObjectResult(entries []types.Object, hasNext bool) *core.ListObjectResult {
	return &core.ListObjectResult{
		Files:      getFiles(entries),
		IsFinished: !hasNext,
	}
}

func getFiles(objects []types.Object) []core.File {
	var files []core.File

	for _, item := range objects {
		if strings.HasSuffix(*item.Key, "/") {
			continue
		}

		files = append(files, &File{item: item})
	}

	return files
}
