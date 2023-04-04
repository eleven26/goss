//go:build integration

package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	config2 "github.com/eleven26/goss/internal/config"

	"github.com/minio/minio-go/v7"

	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/utils"

	fs "github.com/eleven26/go-filesystem"
	"github.com/stretchr/testify/assert"
)

var (
	storage2 core.Storage

	store *Store

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
	storage2, err = d.Storage()
	if err != nil {
		log.Fatal(err)
	}

	store = storage2.Store().(*Store)

	testdata = filepath.Join(utils.RootDir(), "testdata")
	fooPath = filepath.Join(testdata, "foo.txt")
	localFooPath = filepath.Join(testdata, "foo1.txt")

	err = store.CreateBucketIfNotExists()
	if err != nil {
		log.Fatal(err)
	}
}

func setUp(t *testing.T) {
	err := storage2.PutFromFile(key, fooPath)
	if err != nil {
		t.Fatal(err)
	}
}

func tearDown(t *testing.T) {
	deleteLocal(t)
	deleteRemote(t)
}

func deleteRemote(t *testing.T) {
	err := store.client.RemoveObject(context.Background(), store.Bucket, key, minio.RemoveObjectOptions{})
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
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(f)

	err = storage2.Put(key, f)
	assert.Nil(t, err)

	obj, err := store.getObject(key)
	assert.Nil(t, err)

	bs, err := io.ReadAll(obj)
	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}

func TestPutFromFile(t *testing.T) {
	defer tearDown(t)

	err := storage2.PutFromFile(key, fooPath)
	assert.Nil(t, err)

	_, err = store.getObject(key)
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	rc, err := storage2.Get(key)
	assert.Nil(t, err)

	bs, err := io.ReadAll(rc)
	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}

func TestGetString(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	content, err := storage2.GetString(key)

	assert.Nil(t, err)
	assert.Equal(t, "foo", content)
}

func TestGetBytes(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	bs, err := storage2.GetBytes(key)

	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}

func TestDelete(t *testing.T) {
	setUp(t)
	defer deleteLocal(t)

	err := storage2.Delete(key)
	assert.Nil(t, err)

	obj, err := store.getObject(key)
	assert.Nil(t, err)

	_, err = obj.Stat()
	assert.NotNil(t, err)
}

func TestSave(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	err := storage2.GetToFile(key, localFooPath)
	assert.Nil(t, err)
	assert.Equal(t, "foo", fs.MustGetString(localFooPath))
}

func TestExists(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	exists, err := storage2.Exists(key)

	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = storage2.Exists(key + "not_exists")

	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestSize(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	size, err := storage2.Size(key)

	var siz int64 = 3
	assert.Nil(t, err)
	assert.Equal(t, siz, size)
}

func TestFiles(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	files, err := storage2.Files("test/")
	assert.Nil(t, err)
	assert.Len(t, files, 1)

	var expectedSize int64 = 3
	assert.Equal(t, key, files[0].Key())
	assert.Equal(t, expectedSize, files[0].Size())

	today := time.Now().Format("2006-01-02")
	assert.Equal(t, today, files[0].LastModified().Format("2006-01-02"))
}

func putFiles(t *testing.T) {
	dir := "test_all/"

	for i := 1; i <= 200; i++ {
		err := storage2.Put(fmt.Sprintf("%s%s.txt", dir, strconv.Itoa(i)), strings.NewReader("foo"))
		assert.Nil(t, err)
	}
}

func TestFilesWithMultiPage(t *testing.T) {
	putFiles(t)

	// Testdata was prepared before.
	dir := "test_all/"

	files, err := storage2.Files(dir)
	assert.Nil(t, err)
	assert.Len(t, files, 200)
}
