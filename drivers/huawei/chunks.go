package huawei

import (
	"github.com/eleven26/goss/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type Chunks struct {
	prefix string
	bucket string
	client *obs.ObsClient

	nextMarker string
}

func NewChunks(bucket string, prefix string, client *obs.ObsClient) core.Chunks {
	return &Chunks{
		prefix: prefix,
		bucket: bucket,
		client: client,
	}
}

func (c *Chunks) Chunk() (core.ListObjectResult, error) {
	input := &obs.ListObjectsInput{}
	input.Bucket = c.bucket
	input.Marker = c.nextMarker
	input.Prefix = c.prefix

	// 参考文档：https://support.huaweicloud.com/sdk-android-devg-obs/obs_26_0603.html
	// input.maxKeys 列举对象的最大数目，取值范围为1~1000，当超出范围时，按照默认的1000进行处理。
	output, err := c.client.ListObjects(input)
	if err != nil {
		return nil, err
	}

	if output.IsTruncated {
		c.nextMarker = output.NextMarker
	}

	return NewListObjectResult(output), nil
}
