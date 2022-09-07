package qiniu

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/eleven26/goss/core"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Store struct {
	config
	bucketManager *storage.BucketManager
}

func (s *Store) Put(key string, r io.Reader) error {
	upToken := s.upToken()
	formUploader := s.uploader()
	ret := storage.PutRet{}

	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, r)
	if err != nil {
		return err
	}

	err = formUploader.Put(context.Background(), &ret, upToken, key, buf, size, nil)
	if err != nil {
		return err
	}

	return err
}

func (s *Store) upToken() string {
	putPolicy := storage.PutPolicy{
		Scope: s.config.Bucket,
	}

	mac := qbox.NewMac(s.config.AccessKey, s.config.SecretKey)

	return putPolicy.UploadToken(mac)
}

func (s *Store) uploader() *storage.FormUploader {
	cfg := storage.Config{}
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	return storage.NewFormUploader(&cfg)
}

func (s *Store) PutFromFile(key string, localPath string) error {
	upToken := s.upToken()
	formUploader := s.uploader()
	ret := storage.PutRet{}

	return formUploader.PutFile(context.Background(), &ret, upToken, key, localPath, nil)
}

func (s *Store) Get(key string) (io.ReadCloser, error) {
	url := s.getDownloadUrl(key)

	return s.getUrlContent(url)
}

func (s *Store) getDownloadUrl(key string) string {
	var url string

	if s.config.Private {
		mac := qbox.NewMac(s.config.AccessKey, s.config.SecretKey)
		deadline := time.Now().Add(time.Second * 3600).Unix() // 1小时有效期
		url = storage.MakePrivateURL(mac, s.config.Domain, key, deadline)
	} else {
		url = storage.MakePublicURL(s.config.Domain, key)
	}

	return url
}

func (s *Store) getUrlContent(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (s *Store) meta(key string) (http.Header, error) {
	fi, err := s.bucketManager.Stat(s.config.Bucket, key)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Set("Content-Length", strconv.FormatInt(fi.Fsize, 10))

	return header, nil
}

func (s *Store) Size(key string) (int64, error) {
	fi, err := s.bucketManager.Stat(s.config.Bucket, key)
	if err != nil {
		return 0, err
	}

	return fi.Fsize, nil
}

func (s *Store) Exists(key string) (bool, error) {
	_, err := s.meta(key)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Store) Delete(key string) error {
	return s.bucketManager.Delete(s.config.Bucket, key)
}

func (s *Store) Iterator(dir string) core.FileIterator {
	return core.NewFileIterator(&Chunks{
		bucket:        s.config.Bucket,
		bucketManager: s.bucketManager,
		prefix:        dir,
	})
}

type Chunks struct {
	prefix        string
	bucket        string
	bucketManager *storage.BucketManager
	nextMarker    string
}

func (c *Chunks) Chunk() (core.ListObjectResult, error) {
	// 参考文档：https://developer.qiniu.com/kodo/1284/list
	// ListFiles 最后一个参数 limit 为单次列举的条目数，范围为1-1000。 默认值为1000。
	entries, _, nextMarker, hasNext, err := c.bucketManager.ListFiles(c.bucket, c.prefix, "", c.nextMarker, 100)
	if err != nil {
		return nil, err
	}

	c.nextMarker = nextMarker

	return NewListObjectResult(entries, hasNext), nil
}
