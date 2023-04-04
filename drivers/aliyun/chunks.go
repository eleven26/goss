package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/eleven26/goss/v2/core"
)

type Chunks struct {
	prefix string
	bucket *oss.Bucket

	count      int64
	nextMarker string
}

func NewChunks(prefix string, bucket *oss.Bucket) core.Chunks {
	return &Chunks{
		prefix: prefix,
		bucket: bucket,
	}
}

func (c *Chunks) Chunk() (*core.ListObjectResult, error) {
	var result oss.ListObjectsResult
	var err error

	// 参考文档：https://help.aliyun.com/document_detail/31965.html
	// 单次最多返回 100 条，可通过 oss.MaxKeys() 设置单词最大返回条目数量
	if c.count == 0 {
		result, err = c.bucket.ListObjects(oss.Prefix(c.prefix))
	} else {
		result, err = c.bucket.ListObjects(oss.Prefix(c.prefix), oss.Marker(c.nextMarker))
	}

	c.count++
	c.nextMarker = result.NextMarker

	return NewListObjectResult(result), err
}
