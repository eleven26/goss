package core

// ListObjectResult is the result of Chunks.Chunk.
type ListObjectResult interface {
	// Len is the length of current chunk.
	Len() int

	// IsTruncated is used to check if the result is truncated.
	IsTruncated() bool

	// NextMarker return the next marker of next "page".
	NextMarker() interface{}

	// Get file in result by index.
	Get(index int) File
}
