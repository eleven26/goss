package core

type FileIterator interface {
	HasNext() bool
	Next() (File, error)
	All() ([]File, error)
	GetNextChunk() error
}

type chunk func(marker interface{}) (ListObjectResult, error)

type fileIterator struct {
	// 获取下一批数据的 marker
	marker interface{}
	// 获取下一批数据的函数
	nextChunk chunk
	// ListObjects 的返回值
	result ListObjectResult
	// 当前获取的结果遍历到的下标
	index int
	// 当前已获取的结果的总条数
	count int
	// 是否已经获取完全部数据
	isFinished bool
}

func NewFileIterator(marker interface{}, chunk chunk) FileIterator {
	return &fileIterator{
		marker:    marker,
		nextChunk: chunk,
	}
}

// HasNext 还没有获取完全部数据、没有遍历完所有获取到的数据
func (f *fileIterator) HasNext() bool {
	return !f.isFinished || f.index < f.count
}

func (f *fileIterator) Next() (file File, err error) {
	if !f.HasNext() {
		return
	}

	if f.index >= f.count {
		err = f.GetNextChunk()
		if err != nil {
			return
		}
	}

	file = f.result.Get(f.index)

	f.index++

	return
}

func (f *fileIterator) GetNextChunk() error {
	result, err := f.nextChunk(f.marker)
	if err != nil {
		return err
	}

	f.index = 0
	f.count = result.Len()
	f.result = result

	if f.result.IsTruncated() {
		f.marker = result.NextMarker()
	} else {
		f.isFinished = true
	}

	return nil
}

func (f *fileIterator) All() ([]File, error) {
	var res []File

	for {
		if f.HasNext() {
			file, err := f.Next()
			if err != nil {
				return res, err
			}

			res = append(res, file)
		} else {
			break
		}
	}

	return res, nil
}
