package qiniu

import (
	"github.com/eleven26/goss/core"
	"github.com/qiniu/go-sdk/v7/storage"
)

type ListObjectResult struct {
	entries    []storage.ListItem
	nextMarker string
	hasNext    bool
}

func (l *ListObjectResult) Len() int {
	return len(l.entries)
}

func (l *ListObjectResult) IsFinished() bool {
	return !l.hasNext
}

func (l *ListObjectResult) Get(index int) core.File {
	return &File{
		item: l.entries[index],
	}
}
