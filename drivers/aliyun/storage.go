package aliyun

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/eleven26/goss/core"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
)

type Storage struct {
	store Store
}

func (s *Storage) Put(key string, localPath string) error {
	return s.store.Put(key, localPath)
}

func (s *Storage) Get(key string) (content string, err error) {
	rc, err := s.store.Get(key)
	if err != nil {
		return
	}

	defer func() {
		err = rc.Close()
	}()

	bs, err := ioutil.ReadAll(rc)
	if err != nil {
		return
	}

	content = string(bs)

	return
}

func (s *Storage) Save(key string, localPath string) (err error) {
	if !viper.GetBool("show_progress_bar") {
		return s.store.SaveToFile(key, localPath)
	}

	return s.saveWithProgress(key, localPath)
}

func (s Storage) saveWithProgress(key string, localPath string) (err error) {
	rc, err := s.store.Get(key)
	if err != nil {
		return
	}

	defer func() {
		err = rc.Close()
	}()

	// 获取文件长度
	length, err := s.Size(key)
	if err != nil {
		return
	}

	// 保存到文件 localPath
	f, _ := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY, 0o644)
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	// 初始化进度条
	bar := progressbar.DefaultBytes(length, fmt.Sprintf("\"%s\" -> \"%s\"", key, localPath))

	// io.MultiWriter 同时输出到文件和进度条
	_, err = io.Copy(io.MultiWriter(f, bar), rc)

	return
}

func (s *Storage) Size(key string) (int64, error) {
	header, err := s.store.Meta(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(header.Get("Content-Length"), 10, 64)
}

func (s *Storage) Delete(key string) error {
	return s.store.Delete(key)
}

func (s *Storage) Exists(key string) (bool, error) {
	return s.store.Exists(key)
}

func (s *Storage) Files(dir string) ([]core.File, error) {
	return s.store.Iterator(dir).All()
}

func (s *Storage) Storage() interface{} {
	return s
}
