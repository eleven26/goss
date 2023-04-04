package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/eleven26/goss/v2/core"
)

type Chunks struct {
	bucket string
	prefix string

	s3    *s3.Client
	token *string
}

func NewChunks(bucket string, prefix string, s3 *s3.Client) core.Chunks {
	return &Chunks{
		bucket: bucket,
		prefix: prefix,
		s3:     s3,
	}
}

func (c *Chunks) Chunk() (*core.ListObjectResult, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:            aws.String(c.bucket),
		ContinuationToken: c.token,
		Prefix:            aws.String(c.prefix),
	}

	output, err := c.s3.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	c.token = output.ContinuationToken

	return NewListObjectResult(output.Contents, output.IsTruncated), nil
}
