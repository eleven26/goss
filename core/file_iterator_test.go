package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FileStub struct {
	mock.Mock
}

func (f *FileStub) Key() string {
	args := f.Called()

	return args.String(0)
}

func (f *FileStub) Type() string {
	args := f.Called()

	return args.String(0)
}

func (f *FileStub) Size() int64 {
	args := f.Called()

	return args.Get(0).(int64)
}

func (f *FileStub) ETag() string {
	args := f.Called()

	return args.String(0)
}

func (f *FileStub) LastModified() time.Time {
	args := f.Called()

	return args.Get(0).(time.Time)
}

type ResultStub struct {
	mock.Mock
}

func (r *ResultStub) Len() int {
	args := r.Called()

	return args.Int(0)
}

func (r *ResultStub) IsTruncated() bool {
	args := r.Called()

	return args.Bool(0)
}

func (r *ResultStub) NextMarker() interface{} {
	args := r.Called()

	return args.Get(0)
}

func (r *ResultStub) Get(index int) File {
	args := r.Called(index)

	return args.Get(0).(File)
}

type ChunksStub struct {
	mock.Mock
}

func (c *ChunksStub) Chunk(marker interface{}) (ListObjectResult, error) {
	args := c.Called(marker)

	return args.Get(0).(ListObjectResult), args.Error(1)
}

func TestGetNextChunk(t *testing.T) {
	file := new(FileStub)

	result1 := new(ResultStub)
	result1.On("Len").Return(1)
	result1.On("IsTruncated").Return(false)
	result1.AssertNotCalled(t, "NextMarker")
	result1.On("Get", 0).Return(file)

	chunks := new(ChunksStub)
	chunks.On("Chunk", "a").Return(result1, nil)

	fi := fileIterator{
		marker: "a",
		chunks: chunks,
	}
	assert.True(t, fi.HasNext())

	files, err := fi.All()
	assert.Nil(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, file, files[0])

	result1.AssertExpectations(t)
	chunks.AssertExpectations(t)
}
