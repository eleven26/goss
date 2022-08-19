package tencent

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/eleven26/goss/core"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type Storage struct {
	store Store
}

func (s *Storage) Put(key string, r io.Reader) error {
	response, err := s.store.Put(key, r)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return httpError(response)
	}

	return nil
}

func (s *Storage) PutFromFile(key string, localPath string) error {
	response, err := s.store.PutFromFile(key, localPath)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return httpError(response)
	}

	return nil
}

func (s *Storage) Get(key string) (rc io.ReadCloser, err error) {
	resp, err := s.store.Get(key)
	if err != nil {
		return
	}

	return resp.Body, nil
}

func (s *Storage) GetString(key string) (content string, err error) {
	bs, err := s.GetBytes(key)
	if err != nil {
		return
	}

	content = string(bs)

	return
}

func (s *Storage) GetBytes(key string) (bs []byte, err error) {
	rc, err := s.Get(key)
	if err != nil {
		return
	}

	defer func() {
		err = rc.Close()
	}()

	return ioutil.ReadAll(rc)
}

func (s *Storage) GetToFile(key string, localPath string) (err error) {
	resp, err := s.store.SaveToFile(key, localPath)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpError(resp)
	}

	return nil
}

func (s *Storage) Size(key string) (int64, error) {
	response, err := s.store.Head(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(response.Header.Get("Content-Length"), 10, 64)
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

func httpError(response *cos.Response) error {
	bytes, err := ioutil.ReadAll(response.Body)
	defer func() {
		err = response.Body.Close()
	}()
	if err != nil {
		return err
	}

	return errors.New(string(bytes))
}
