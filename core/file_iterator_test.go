package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ChunksStub struct {
	mock.Mock
}

func (c *ChunksStub) Chunk() (*ListObjectResult, error) {
	args := c.Called()

	return args.Get(0).(*ListObjectResult), args.Error(1)
}

func TestNewFileIterator(t *testing.T) {
	chunks := new(ChunksStub)
	fi := NewFileIterator(chunks)

	assert.Equal(t, chunks, fi.(*fileIterator).chunks)
}

func TestNotHasNext(t *testing.T) {
	result := &ListObjectResult{}

	chunks := new(ChunksStub)
	chunks.On("Chunk").Return(result, nil)

	fi := fileIterator{
		chunks: chunks,
	}
	assert.False(t, fi.hasNext())

	chunks.AssertExpectations(t)
}

func TestHasNext(t *testing.T) {
	result := &ListObjectResult{
		Files:      make([]File, 1),
		IsFinished: false,
	}

	chunks := new(ChunksStub)
	chunks.On("Chunk").Return(result, nil)

	fi := fileIterator{
		chunks: chunks,
	}
	assert.True(t, fi.hasNext())

	chunks.AssertExpectations(t)
}

func TestNotHasNext1(t *testing.T) {
	emptyResult := &ListObjectResult{}

	chunks := new(ChunksStub)
	chunks.On("Chunk").Return(emptyResult, nil)

	fi := fileIterator{
		index:  100,
		count:  100,
		chunks: chunks,
	}
	assert.False(t, fi.hasNext())

	fi = fileIterator{
		index:       100,
		count:       100,
		isFinished:  true,
		chunksCount: 1,
	}
	assert.False(t, fi.hasNext())
}

func TestGetNextChunk(t *testing.T) {
	result := &ListObjectResult{
		Files:      make([]File, 1),
		IsFinished: false,
	}

	chunks := new(ChunksStub)
	chunks.On("Chunk").Return(result, nil)

	fi := fileIterator{
		chunks: chunks,
	}
	err := fi.getNextChunk()
	assert.Nil(t, err)

	chunks.AssertExpectations(t)
}

func TestAll(t *testing.T) {
	result := &ListObjectResult{
		Files:      make([]File, 2),
		IsFinished: true,
	}
	chunks := new(ChunksStub)
	chunks.On("Chunk").Return(result, nil)

	fi := fileIterator{
		index:       0,
		count:       2,
		chunks:      chunks,
		chunksCount: 0,
	}
	files, err := fi.All()
	assert.Nil(t, err)
	assert.Len(t, files, 2)
}
