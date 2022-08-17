package qiniu

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/eleven26/goss/core"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Store struct {
	config
	bucketManager *storage.BucketManager
}

func (s *Store) put(key string, localPath string) (*storage.PutRet, error) {
	putPolicy := storage.PutPolicy{
		Scope: s.config.Bucket,
	}

	mac := qbox.NewMac(s.config.AccessKey, s.config.SecretKey)

	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
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

func (s *Store) getUrlContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(resp.Body)

	return ioutil.ReadAll(resp.Body)
}

func (s *Store) Stat(key string) (fileInfo storage.FileInfo, err error) {
	return s.bucketManager.Stat(s.config.Bucket, key)
}

func (s *Store) SaveToFile(key string, localPath string) error {
	url := s.getDownloadUrl(key)
	bs, err := s.getUrlContent(url)
	if err != nil {
		return err
	}

	// 保存到文件 localPath
	f, _ := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY, 0o644)
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	_, err = f.Write(bs)

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
