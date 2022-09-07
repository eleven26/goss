package core

import "strings"

// FileIterator is an iterator used to iterate over all objects in the cloud.
type FileIterator interface {
	// HasNext Determine if there is a next object.
	HasNext() bool

	// Next Get the next object.
	Next() (File, error)

	// All Get all objects.
	All() ([]File, error)

	// GetNextChunk Get the next batch of objects.
	GetNextChunk() error
}

// Chunks is used to get the next "page" of objects.
type Chunks interface {
	Chunk(marker interface{}) (ListObjectResult, error)
}

// fileIterator is the iterator used to iterate over all matching objects.
type fileIterator struct {
	marker      interface{}      // used to get next "page".
	chunks      Chunks           // provides next "page" for current iterator.
	result      ListObjectResult // the return value of chunks
	index       int              // the index of result
	count       int              // the length of result
	isFinished  bool             // Whether all data has been obtained.
	chunksCount int              // counter for Chunks.Chunk
}

// NewFileIterator creates an instance of FileIterator.
func NewFileIterator(marker interface{}, chunks Chunks) FileIterator {
	return &fileIterator{
		marker: marker,
		chunks: chunks,
	}
}

// HasNext check if there is a next object.
func (f *fileIterator) HasNext() bool {
	// No chunk have been got.
	if f.chunksCount == 0 {
		err := f.GetNextChunk()
		if err != nil {
			return false
		}
	}

	// No more chunks.
	if f.isFinished {
		return f.index < f.count
	}

	// Get next "page".
	if f.index == f.count {
		err := f.GetNextChunk()
		if err != nil {
			return false
		}
	}

	return f.index < f.count
}

// Next get the next object.
func (f *fileIterator) Next() (file File, err error) {
	if !f.HasNext() {
		return
	}

	file = f.result.Get(f.index)

	f.index++

	return
}

// GetNextChunk get next "page" of objects.
func (f *fileIterator) GetNextChunk() error {
	f.chunksCount++
	return f.handleChunkResult(f.chunks.Chunk(f.marker))
}

func (f *fileIterator) handleChunkResult(result ListObjectResult, err error) error {
	if err != nil {
		return err
	}

	f.index = 0
	f.count = result.Len()
	f.result = result

	if result.IsTruncated() {
		f.marker = result.NextMarker()
	} else {
		f.isFinished = true
	}

	return nil
}

// All get all objects.
func (f *fileIterator) All() ([]File, error) {
	var res []File

	for {
		if f.HasNext() {
			file, err := f.Next()
			if err != nil {
				return res, err
			}

			if strings.HasSuffix(file.Key(), "/") {
				continue
			}

			res = append(res, file)
		} else {
			break
		}
	}

	return res, nil
}
