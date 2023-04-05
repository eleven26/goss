package goss

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type listObjectResult struct {
	Files      []File
	IsFinished bool
}

func newListObjectResult(entries []types.Object, hasNext bool) *listObjectResult {
	return &listObjectResult{
		Files:      getFiles(entries),
		IsFinished: !hasNext,
	}
}

func getFiles(objects []types.Object) []File {
	var files []File

	for _, item := range objects {
		if strings.HasSuffix(*item.Key, "/") {
			continue
		}

		files = append(files, &file{item: item})
	}

	return files
}
