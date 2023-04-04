package qiniu

import (
	"github.com/eleven26/goss/v2/core"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Chunks struct {
	bucket        string
	prefix        string
	bucketManager *storage.BucketManager

	nextMarker string
}

func NewChunks(bucket string, prefix string, bucketManager *storage.BucketManager) core.Chunks {
	return &Chunks{
		bucket:        bucket,
		prefix:        prefix,
		bucketManager: bucketManager,
	}
}

func (c *Chunks) Chunk() (*core.ListObjectResult, error) {
	// 参考文档：https://developer.qiniu.com/kodo/1284/list
	// ListFiles 最后一个参数 limit 为单次列举的条目数，范围为1-1000。 默认值为1000。
	entries, _, nextMarker, hasNext, err := c.bucketManager.ListFiles(c.bucket, c.prefix, "", c.nextMarker, 100)
	if err != nil {
		return nil, err
	}

	c.nextMarker = nextMarker

	return NewListObjectResult(entries, hasNext), nil
}
