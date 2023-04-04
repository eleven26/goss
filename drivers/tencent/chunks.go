package tencent

import (
	"context"

	"github.com/eleven26/goss/v2/core"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type Chunks struct {
	prefix string
	bucket *cos.BucketService

	count      int
	nextMarker string
}

func NewChunks(prefix string, bucket *cos.BucketService) core.Chunks {
	return &Chunks{
		prefix: prefix,
		bucket: bucket,
	}
}

func (c *Chunks) Chunk() (*core.ListObjectResult, error) {
	var opt cos.BucketGetOptions
	var result *cos.BucketGetResult
	var err error

	// 参考文档：https://cloud.tencent.com/document/product/436/7734
	// 单次返回最大的条目数量，默认值为1000，最大为1000
	// BucketGetOptions.MaxKeys 可以设置单次获取的条目数量
	if c.count == 0 {
		opt = cos.BucketGetOptions{Prefix: c.prefix}
	} else {
		opt = cos.BucketGetOptions{Prefix: c.prefix, Marker: c.nextMarker}
	}

	result, _, err = c.bucket.Get(context.Background(), &opt)
	if err != nil {
		return nil, err
	}

	c.count++
	c.nextMarker = result.NextMarker

	return NewListObjectResult(result), nil
}
