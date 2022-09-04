package core

import (
	"errors"
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

func TestNewFileIterator(t *testing.T) {
	chunks := new(ChunksStub)
	fi := NewFileIterator("foo", chunks)

	assert.Equal(t, "foo", fi.(*fileIterator).marker)
	assert.Equal(t, chunks, fi.(*fileIterator).chunks)
}

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
	emptyResult := new(ResultStub)
	emptyResult.On("Len").Return(0)
	emptyResult.On("IsTruncated").Return(false)

	chunks := new(ChunksStub)
	chunks.On("Chunk", "foo").Return(emptyResult, nil)

	fi := fileIterator{
		index:  100,
		count:  100,
		marker: "foo",
		chunks: chunks,
	}
	assert.False(t, fi.HasNext())

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

func TestHandleChunkResult(t *testing.T) {
	result := new(ResultStub)
	err := errors.New("foo")
	fi := fileIterator{}
	assert.ErrorIs(t, fi.handleChunkResult(result, err), err)

	result = new(ResultStub)
	result.On("Len").Return(10)
	result.On("IsTruncated").Return(false)
	fi = fileIterator{
		index:      10,
		count:      20,
		result:     nil,
		isFinished: false,
	}
	err = fi.handleChunkResult(result, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, fi.index)
	assert.Equal(t, 10, fi.count)
	assert.Equal(t, result, fi.result)
	assert.True(t, fi.isFinished)

	result = new(ResultStub)
	result.On("Len").Return(1)
	result.On("IsTruncated").Return(true)
	result.On("NextMarker").Return("foo")
	fi = fileIterator{
		index:      10,
		count:      20,
		result:     nil,
		isFinished: false,
	}
	err = fi.handleChunkResult(result, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, fi.index)
	assert.Equal(t, 1, fi.count)
	assert.Equal(t, result, fi.result)
	assert.False(t, fi.isFinished)
	assert.Equal(t, "foo", fi.marker)
}

func TestAll(t *testing.T) {
	file := new(FileStub)

	result := new(ResultStub)
	result.On("Get", 0).Return(file)
	result.On("Get", 1).Return(file)

	emptyResult := new(ResultStub)
	emptyResult.On("Len").Return(0)
	emptyResult.On("IsTruncated").Return(false)

	chunks := new(ChunksStub)
	chunks.On("Chunk", "foo").Return(emptyResult, nil)

	fi := fileIterator{
		index:  0,
		count:  2,
		result: result,
		marker: "foo",
		chunks: chunks,
	}
	files, err := fi.All()
	assert.Nil(t, err)
	assert.Len(t, files, 2)

	result.AssertExpectations(t)
}
