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
	// TODO implement me
	panic("implement me")
}

func (f *FileStub) Type() string {
	// TODO implement me
	panic("implement me")
}

func (f *FileStub) Size() int64 {
	// TODO implement me
	panic("implement me")
}

func (f *FileStub) ETag() string {
	// TODO implement me
	panic("implement me")
}

func (f *FileStub) LastModified() time.Time {
	// TODO implement me
	panic("implement me")
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

//func chunks(marker interface{}) (ListObjectResult, error) {
//	return nil, nil
//}

func TestGetNextChunk(t *testing.T) {
	fi := fileIterator{
		marker:     nil,
		result:     nil,
		index:      0,
		count:      0,
		isFinished: false,
	}

	assert.True(t, fi.HasNext())
}
