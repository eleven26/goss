package goss

// FileIterator is an iterator used to iterate over all objects in the cloud.
type FileIterator interface {
	// All Get all objects.
	All() ([]File, error)
}

// Chunks is used to get the next "page" of objects.
type Chunks interface {
	Chunk() (*listObjectResult, error)
}

// fileIterator is the iterator used to iterate over all matching objects.
type fileIterator struct {
	chunks      Chunks // provides next "page" for current iterator.
	index       int    // the index of result
	count       int    // the length of result
	isFinished  bool   // Whether all data has been obtained.
	chunksCount int    // counter for Chunks.Chunk

	files []File
}

// newFileIterator creates an instance of FileIterator.
func newFileIterator(chunks Chunks) FileIterator {
	return &fileIterator{
		chunks: chunks,
	}
}

// hasNext check if there is a next object.
func (f *fileIterator) hasNext() bool {
	if f.shouldGetNextChunk() {
		err := f.getNextChunk()
		if err != nil {
			return false
		}
	}

	return f.index < f.count
}

func (f *fileIterator) shouldGetNextChunk() bool {
	if f.chunksCount == 0 {
		return true
	}

	if f.isFinished {
		return false
	}

	return f.index == f.count
}

// next get the next object.
func (f *fileIterator) next() File {
	file := f.files[f.index]

	f.index++

	return file
}

// getNextChunk get next "page" of objects.
func (f *fileIterator) getNextChunk() error {
	f.chunksCount++

	result, err := f.chunks.Chunk()
	if err != nil {
		return err
	}

	f.index = 0
	f.files = result.Files
	f.count = len(f.files)
	f.isFinished = result.IsFinished

	return nil
}

// All get all objects.
func (f *fileIterator) All() ([]File, error) {
	var res []File

	for {
		if !f.hasNext() {
			break
		}

		res = append(res, f.next())
	}

	return res, nil
}
