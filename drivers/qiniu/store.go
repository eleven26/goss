package qiniu

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/eleven26/goss/v2/core"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Store struct {
	config
	bucketManager *storage.BucketManager
}

func (s *Store) put(key string, r io.Reader) (*storage.PutRet, error) {
	upToken := s.upToken()
	formUploader := s.uploader()
	ret := storage.PutRet{}

	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	err = formUploader.Put(context.Background(), &ret, upToken, key, buf, size, nil)
	if err != nil {
		return nil, err
	}

	return &ret, nil
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

func (s *Store) putFromFile(key string, localPath string) (*storage.PutRet, error) {
	upToken := s.upToken()
	formUploader := s.uploader()
	ret := storage.PutRet{}

	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localPath, nil)
	if err != nil {
		return nil, err
	}

	return &ret, nil
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

func (s *Store) Stat(key string) (fileInfo storage.FileInfo, err error) {
	return s.bucketManager.Stat(s.config.Bucket, key)
}

func (s *Store) SaveToFile(key string, localPath string) error {
	url := s.getDownloadUrl(key)
	rc, err := s.getUrlContent(url)
	if err != nil {
		return err
	}

	defer func(rc io.ReadCloser) {
		err = rc.Close()
	}(rc)

	// 保存到文件 localPath
	f, _ := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY, 0o644)
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	_, err = io.Copy(f, rc)

	return err
}

func (s *Store) Delete(key string) error {
	return s.bucketManager.Delete(s.config.Bucket, key)
}

func (s *Store) Iterator(dir string) core.FileIterator {
	chunk := func(marker string) (entries []storage.ListItem, commonPrefixes []string, nextMarker string, hasNext bool, err error) {
		return s.bucketManager.ListFiles(s.config.Bucket, dir, "", marker, 100)
	}

	it := newFileIterator(dir, chunk)

	return &it
}
