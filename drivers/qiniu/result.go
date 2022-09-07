package qiniu

import (
	"strings"

	"github.com/eleven26/goss/core"
	"github.com/qiniu/go-sdk/v7/storage"
)

type ListObjectResult struct {
	entries    []storage.ListItem
	nextMarker string
	hasNext    bool
	files      []core.File
}

func NewListObjectResult(entries []storage.ListItem, nextMarker string, hasNext bool) core.ListObjectResult {
	result := ListObjectResult{
		entries:    entries,
		nextMarker: nextMarker,
		hasNext:    hasNext,
	}

	result.files = result.getFiles()

	return &result
}

func (l *ListObjectResult) getFiles() []core.File {
	var files []core.File

	for _, item := range l.entries {
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
	return !l.hasNext
}
