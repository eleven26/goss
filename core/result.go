package core

// ListObjectResult is the result of Chunks.Chunk.
type ListObjectResult interface {
	// Len is the length of current chunk.
	Len() int

	// IsFinished is used to check if the result is finished.
	IsFinished() bool

	// Files return all files of the result.
	Files() []File
}
