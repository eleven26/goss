# goss

✨ `goss` 是一个简洁的云存储 golang 库，支持**阿里云**、**腾讯云**、**七牛云**。

[![Go Reference](https://pkg.go.dev/badge/github.com/eleven26/go-filesystem.svg)](https://pkg.go.dev/github.com/eleven26/goss)
[![Go Report Card](https://goreportcard.com/badge/github.com/eleven26/go-filesystem)](https://goreportcard.com/report/github.com/eleven26/goss)
[![Go](https://github.com/eleven26/goss/actions/workflows/go.yml/badge.svg)](https://github.com/eleven26/goss/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/eleven26/goss/branch/main/graph/badge.svg?token=UU4lLD2n4k)](https://codecov.io/gh/eleven26/goss)
[![GitHub license](https://img.shields.io/github/license/eleven26/goss)](https://github.com/eleven26/goss/blob/main/LICENSE)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/eleven26/goss)


## 🚀 安装

```
go get -u github.com/eleven26/goss
```


## ⚙️ 配置

所有支持的配置项：

```yaml
# 云存储类型
# 可选值为： aliyun、tencent，不区分大小写
driver: aliyun

# 阿里云 oss 配置
aliyun:
  # oss 的链接，不同区域不同
  endpoint:
  # bucket
  bucket:
  access_key_id:
  access_key_secret:

# 腾讯云 cos 配置 
tencent:
  # 腾讯云 bucket 对应的的 url
  url:
  secret_id:
  secret_key:

# 七牛云 kodo 配置
qiniu:
  # bucket 名称
  bucket:
  access_key:
  secret_key:
  # bucket 外链域名
  domain:
  # 是否是私有空间
  private:
```

样例配置：

> 比如，如果只是使用阿里云 oss，则只需添加以下配置项就可以了：

```yaml
driver: aliyun

aliyun:
  endpoint: oss-cn-shenzhen.aliyuncs.com
  bucket: images
  access_key_id: LT2I316210b3JlXj
  access_key_secret: 4IZq10e233Ya1ZS18JDG0ZfvBBnYva
```


## 💡 基本用法

1. 你可以通过下面的代码来导入 `goss`:

```go
import "github.com/eleven26/goss/goss"
```

2. 使用之前需要创建实例：

```go
// path 是配置文件的路径
path := "./goss.yml"
// storage 是云存储对象
storage := goss.Storage(path)
```

3. 使用

```go
// storage.Get 会获取路径指定的文件，返回字符串
fmt.Println(storage.Get("test/foo.txt"))
```


## 📚 接口

`goss` 支持以下操作：

- [Put](#Put)
- [PutFromFile](#PutFromFile)
- [Get](#Get)
- [GetString](#GetString)
- [GetBytes](#GetBytes)
- [GetToFile](#GetToFile)
- [Delete](#Delete)
- [Exists](#Exists)
- [Files](#Files)
- [Size](#Size)

### Put

上传文件到云存储。第一个参数是 key，第二个参数是 `io.Reader`。

```go
data := []byte("this is some data stored as a byte slice in Go Lang!")
r := bytes.NewReader(data)
err := storage.Put("test/test.txt", r)
```

### PutFromFile

上传文件到云存储。第一个参数是 key，第二个参数是本地文件路径。

```go
err := storage.PutFromFile("test/test.txt", "/path/to/test.txt")
```

### Get

从云存储获取文件，返回字符串。参数是 key。返回值是 `io.ReadCloser` 和 `error`。

```go
// rc 是 `io.ReadCloser`
rc, err := storage.Get("test/test.txt")
defer rc.Close()

bs, err := ioutil.ReadAll(rc)
fmt.Println(string(bs))
```

### GetString

从云存储获取文件，返回字符串。参数是 key。返回值是 `string` 和 `error`

```go
content, err := storage.GetString("test/test.txt")
fmt.Println(content)
```

### GetBytes

从云存储获取文件，返回字符串。参数是 key。返回值是 `[]byte` 和 `error`

```go
bs, err := storage.Get("test/test.txt")
fmt.Println(string(bs))
```

### GetToFile

下载云存储文件到本地。第一个参数是 key，第二个参数是本地路径。

```go
// 第一个参数是云端路径，第二个参数是本地路径
err := storage.GetToFile("test/test.txt", "/path/to/local")
```

### Delete

删除云存储文件。

```go
err := storage.Delete("test/test.txt")
```

### Exists

判断云存储文件是否存在。

```go
exists, err := storage.Exists("test/test.txt")
```

### Files

根据前缀获取文件列表。

```go
exists, err := storage.Files("test/")
```

### Size

获取云存储文件大小。

```go
size, err := storage.Size("test/test.txt")
```
