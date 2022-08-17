package core

type Storage interface {
	// Put 将本地路径 localPath 指向的文件保存到 oss 的 key
	Put(key string, localPath string) error
	// Get 获取 key 指向的文件
	Get(key string) (string, error)
	// Delete 删除 key 指向的文件
	Delete(key string) error
	// Save 保存 key 指向的文件到本地 localPath
	Save(key string, localPath string) error
	// Exists 判断文件是否存在
	Exists(key string) (bool, error)
	// Files 列出指定目录下的所有文件
	Files(dir string) ([]File, error)
	// Size 获取文件大小
	Size(key string) (int64, error)
	// Storage 获取底层的结构体
	Storage() interface{}
}
