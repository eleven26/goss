//go:build integration

package tencent

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	config2 "github.com/eleven26/goss/config"
	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/utils"

	fs "github.com/eleven26/go-filesystem"
	"github.com/stretchr/testify/assert"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	storage core.Storage

	object *cos.ObjectService

	key          = "test/foo.txt"
	testdata     string
	fooPath      string
	localFooPath string
)

func init() {
	vip, err := config2.ReadInUserHomeConfig()
	if err != nil {
		log.Fatal(err)
	}

	d := NewDriver(core.WithViper(vip))
	storage, err = d.Storage()
	if err != nil {
		log.Fatal(err)
	}

	object = storage.Store().(*Store).client.Object

	testdata = filepath.Join(utils.RootDir(), "testdata")
	fooPath = filepath.Join(testdata, "foo.txt")
	localFooPath = filepath.Join(testdata, "foo1.txt")
}

func setUp(t *testing.T) {
	_, err := object.PutFromFile(context.Background(), key, fooPath, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func tearDown(t *testing.T) {
	deleteLocal(t)
	deleteRemote(t)
}

func deleteRemote(t *testing.T) {
	_, err := object.Delete(context.Background(), key)
	if err != nil {
		t.Fatal(err)
	}
}

func deleteLocal(t *testing.T) {
	exists, _ := fs.Exists(localFooPath)
	if exists {
		err := fs.Delete(localFooPath)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestPut(t *testing.T) {
	defer tearDown(t)

	f, err := os.Open(fooPath)
	if err != nil {
		t.Fatal(err)
	}

	err = storage.Put(key, f)
	assert.Nil(t, err)

	exists, err := object.IsExist(context.Background(), key)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestPutFromFile(t *testing.T) {
	defer tearDown(t)

	err := storage.PutFromFile(key, fooPath)
	assert.Nil(t, err)

	exists, err := object.IsExist(context.Background(), key)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestGet(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	rc, err := storage.Get(key)
	defer func(rc io.ReadCloser) {
		err = rc.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(rc)
	assert.Nil(t, err)

	bs, err := io.ReadAll(rc)
	assert.Nil(t, err)
	assert.Equal(t, string(bs), "foo")

	rc, err = storage.Get(key + "not_exists")
	assert.Empty(t, rc)
	assert.Equal(t, http.StatusNotFound, err.(*cos.ErrorResponse).Response.StatusCode)
}

func TestGetString(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	content, err := storage.GetString(key)
	assert.Nil(t, err)
	assert.Equal(t, content, "foo")

	content, err = storage.GetString(key + "not_exists")
	assert.Empty(t, content)
	assert.Equal(t, http.StatusNotFound, err.(*cos.ErrorResponse).Response.StatusCode)
}

func TestGetBytes(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	bs, err := storage.GetBytes(key)
	assert.Nil(t, err)
	assert.Equal(t, string(bs), "foo")

	bs, err = storage.GetBytes(key + "not_exists")
	assert.Nil(t, bs)
	assert.Equal(t, http.StatusNotFound, err.(*cos.ErrorResponse).Response.StatusCode)
}

func TestSave(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	err := storage.GetToFile(key, localFooPath)
	assert.Nil(t, err)
	assert.Equal(t, "foo", fs.MustGetString(localFooPath))
}

func TestSize(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	size, err := storage.Size(key)

	var expectedSize int64 = 3
	assert.Nil(t, err)
	assert.Equal(t, expectedSize, size)

	var s int64 = 0
	size, err = storage.Size(key + "not_exists")
	assert.Equal(t, s, size)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*cos.ErrorResponse).Response.StatusCode)
}

func TestDelete(t *testing.T) {
	setUp(t)

	err := storage.Delete(key)
	assert.Nil(t, err)

	exists, err := object.IsExist(context.Background(), key)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestExists(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	exists, err := storage.Exists(key)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestFiles(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	files, err := storage.Files("test/")
	assert.Nil(t, err)
	assert.Len(t, files, 1)

	var expectedSize int64 = 3
	assert.Equal(t, "test/foo.txt", files[0].Key())
	assert.Equal(t, expectedSize, files[0].Size())
}

func TestFilesWithMultiPage(t *testing.T) {
	// Testdata was prepared before.
	dir := "test_all/"

	files, err := storage.Files(dir)
	assert.Nil(t, err)
	assert.Len(t, files, 200)
}
