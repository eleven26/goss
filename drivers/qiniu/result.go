package qiniu

import (
	"strings"

	"github.com/eleven26/goss/core"
	"github.com/qiniu/go-sdk/v7/storage"
)

type ListObjectResult struct {
	files      []core.File
	isFinished bool
}

func NewListObjectResult(entries []storage.ListItem, hasNext bool) core.ListObjectResult {
	return &ListObjectResult{
		isFinished: !hasNext,
		files:      getFiles(entries),
	}
}

func getFiles(entries []storage.ListItem) []core.File {
	var files []core.File

	for _, item := range entries {
		if strings.HasSuffix(item.Key, "/") {
			continue
		}

		files = append(files, &File{item: item})
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
