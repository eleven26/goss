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

//func TestNewFileIterator(t *testing.T) {
//	it := NewFileIterator("", new(ChunksStub))
//
//	assert.True(t, it.HasNext())
//}

func TestNotHasNext(t *testing.T) {
	result := new(ResultStub)
	result.On("Len").Return(0)
	result.On("IsTruncated").Return(false)

	chunks := new(ChunksStub)
	chunks.On("Chunk", "first").Return(result, nil)

	fi := fileIterator{
		marker: "first",
		chunks: chunks,
	}
	assert.False(t, fi.HasNext())

	result.AssertExpectations(t)
	chunks.AssertExpectations(t)
}

func TestHasNext(t *testing.T) {
	result := new(ResultStub)
	result.On("Len").Return(1)
	result.On("IsTruncated").Return(false)

	chunks := new(ChunksStub)
	chunks.On("Chunk", "first").Return(result, nil)

	fi := fileIterator{
		marker: "first",
		chunks: chunks,
	}
	assert.True(t, fi.HasNext())

	result.AssertExpectations(t)
	chunks.AssertExpectations(t)
}

func TestNotHasNext1(t *testing.T) {
	fi := fileIterator{
		index: 100,
		count: 100,
	}
	assert.True(t, fi.HasNext())

	fi = fileIterator{
		index:      100,
		count:      100,
		isFinished: true,
	}
	assert.False(t, fi.HasNext())
}

func TestGetNextChunk(t *testing.T) {
	result := new(ResultStub)
	result.On("Len").Return(1)
	result.On("IsTruncated").Return(false)

	chunks := new(ChunksStub)
	chunks.On("Chunk", "first").Return(result, nil)

	fi := fileIterator{
		marker: "first",
		chunks: chunks,
	}
	err := fi.GetNextChunk()
	assert.Nil(t, err)

	result.AssertExpectations(t)
	chunks.AssertExpectations(t)
}

func TestAll(t *testing.T) {

}
