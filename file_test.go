package goss

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	now := aws.Time(time.Now())

	f := file{item: types.Object{
		Key:          aws.String("test"),
		StorageClass: types.ObjectStorageClass("test"),
		ETag:         aws.String("test"),
		Size:         1,
		LastModified: now,
	}}

	assert.Equal(t, "test", f.Key())
	assert.Equal(t, "test", f.Type())
	assert.Equal(t, "test", f.ETag())
	assert.Equal(t, int64(1), f.Size())
	assert.Equal(t, now, f.LastModified())
}
