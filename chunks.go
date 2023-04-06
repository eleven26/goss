package goss

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var maxKeys int32 = 1000

type chunks struct {
	bucket string
	prefix string

	s3    *s3.Client
	token *string
}

func newChunks(bucket string, prefix string, s3 *s3.Client) Chunks {
	return &chunks{
		bucket: bucket,
		prefix: prefix,
		s3:     s3,
	}
}

func (c *chunks) Chunk() (*listObjectResult, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:            aws.String(c.bucket),
		ContinuationToken: c.token,
		Prefix:            aws.String(c.prefix),
		MaxKeys:           maxKeys,
	}

	output, err := c.s3.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	c.token = output.NextContinuationToken

	return newListObjectResult(output.Contents, output.IsTruncated), nil
}
