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

func (r *ResultStub) IsFinished() bool {
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
	marker string
}

func (c *ChunksStub) Chunk() (ListObjectResult, error) {
	args := c.Called()

	return args.Get(0).(ListObjectResult), args.Error(1)
}

func TestNewFileIterator(t *testing.T) {
	chunks := new(ChunksStub)
	chunks.marker = "foo"
	fi := NewFileIterator(chunks)

	assert.Equal(t, chunks, fi.(*fileIterator).chunks)
}

func TestNotHasNext(t *testing.T) {
	result := new(ResultStub)
	result.On("Len").Return(0)
	result.On("IsFinished").Return(false)

	chunks := new(ChunksStub)
	chunks.On("Chunk").Return(result, nil)

	fi := fileIterator{
		chunks: chunks,
	}
	assert.False(t, fi.HasNext())

	result.AssertExpectations(t)
	chunks.AssertExpectations(t)
}

func TestHasNext(t *testing.T) {
	result := new(ResultStub)
	result.On("Len").Return(1)
	result.On("IsFinished").Return(false)

	chunks := new(ChunksStub)
	chunks.On("Chunk").Return(result, nil)

	fi := fileIterator{
		chunks: chunks,
	}
	assert.True(t, fi.HasNext())

	result.AssertExpectations(t)
	chunks.AssertExpectations(t)
}

func TestNotHasNext1(t *testing.T) {
	emptyResult := new(ResultStub)
	emptyResult.On("Len").Return(0)
	emptyResult.On("IsFinished").Return(false)

	chunks := new(ChunksStub)
	chunks.On("Chunk").Return(emptyResult, nil)

	fi := fileIterator{
		index:  100,
		count:  100,
		chunks: chunks,
	}
	assert.False(t, fi.HasNext())

	fi = fileIterator{
		index:       100,
		count:       100,
		isFinished:  true,
		chunksCount: 1,
	}
	assert.False(t, fi.HasNext())
}

func TestGetNextChunk(t *testing.T) {
	result := new(ResultStub)
	result.On("Len").Return(1)
	result.On("IsFinished").Return(false)

	chunks := new(ChunksStub)
	chunks.On("Chunk").Return(result, nil)

	fi := fileIterator{
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
	result.On("IsFinished").Return(false)
	fi = fileIterator{
		index:  10,
		count:  20,
		result: nil,
	}
	err = fi.handleChunkResult(result, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, fi.index)
	assert.Equal(t, 10, fi.count)
	assert.Equal(t, result, fi.result)
	assert.False(t, fi.isFinished)

	result = new(ResultStub)
	result.On("Len").Return(1)
	result.On("IsFinished").Return(true)
	result.On("NextMarker").Return("foo")
	fi = fileIterator{
		index:  10,
		count:  20,
		result: nil,
	}
	err = fi.handleChunkResult(result, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, fi.index)
	assert.Equal(t, 1, fi.count)
	assert.Equal(t, result, fi.result)
	assert.True(t, fi.isFinished)
}

func TestAll(t *testing.T) {
	file := new(FileStub)
	file.On("Key").Return("foo")

	result := new(ResultStub)
	result.On("Get", 0).Return(file)
	result.On("Get", 1).Return(file)
	result.On("Len").Return(2)
	result.On("IsFinished").Return(true)

	chunks := new(ChunksStub)
	chunks.On("Chunk").Return(result, nil)

	fi := fileIterator{
		index:       0,
		count:       2,
		result:      result,
		chunks:      chunks,
		chunksCount: 0,
	}
	files, err := fi.All()
	assert.Nil(t, err)
	assert.Len(t, files, 2)

	result.AssertExpectations(t)
}
