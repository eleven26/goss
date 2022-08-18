package core

import "io"

type Storage interface {
	// Put 将从 r 读取的内容保存到 oss 的 key
	Put(key string, r io.Reader) error
	// PutFromFile 将本地路径 localPath 指向的文件保存到 oss 的 key
	PutFromFile(key string, localPath string) error

	// Get 获取 key 指向的文件
	Get(key string) (io.ReadCloser, error)
	// GetString 获取 key 指向的文件，返回字符串
	GetString(key string) (string, error)
	// GetBytes 获取 key 指向的文件，返回字节数组
	GetBytes(key string) ([]byte, error)
	// GetToFile 保存 key 指向的文件到本地 localPath
	GetToFile(key string, localPath string) error

	// Delete 删除 key 指向的文件
	Delete(key string) error
	// Exists 判断文件是否存在
	Exists(key string) (bool, error)
	// Files 列出指定目录下的所有文件
	Files(dir string) ([]File, error)
	// Size 获取文件大小
	Size(key string) (int64, error)

	// Storage 获取底层的结构体
	Storage() interface{}
}
