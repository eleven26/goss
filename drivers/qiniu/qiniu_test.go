//go:build integration

package qiniu

import (
	"path/filepath"
	"testing"
	"time"

	config2 "github.com/eleven26/goss/config"
	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/utils"

	fs "github.com/eleven26/go-filesystem"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	storage2 core.Storage

	store Store

	key          = "test/foo.txt"
	testdata     string
	fooPath      string
	localFooPath string
)

func init() {
	config2.ReadInUserHomeConfig()

	d := NewDriver()
	storage2 = d.Storage()

	store = storage2.(*Storage).store

	testdata = filepath.Join(utils.RootDir(), "testdata")
	fooPath = filepath.Join(testdata, "foo.txt")
	localFooPath = filepath.Join(testdata, "foo1.txt")

	viper.Set("show_progress_bar", false)
}

func setUp() {
	err := storage2.Put(key, fooPath)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	deleteLocal()
	deleteRemote()
}

func deleteRemote() {
	err := store.bucketManager.Delete(store.config.Bucket, key)
	if err != nil {
		panic(err)
	}
}

func deleteLocal() {
	exists, _ := fs.Exists(localFooPath)
	if exists {
		_ = fs.Delete(localFooPath)
	}
}

func TestPut(t *testing.T) {
	defer tearDown()

	err := storage2.Put(key, fooPath)

	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	setUp()
	defer tearDown()

	content, err := storage2.Get(key)

	assert.Nil(t, err)
	assert.Equal(t, "foo", content)
}

func TestDelete(t *testing.T) {
	setUp()
	defer deleteLocal()

	err := storage2.Delete(key)
	assert.Nil(t, err)

	_, err = store.bucketManager.Stat(store.config.Bucket, key)
	assert.NotNil(t, err)
}

func TestSave(t *testing.T) {
	setUp()
	defer tearDown()

	err := storage2.Save(key, localFooPath)
	assert.Nil(t, err)
	assert.Equal(t, "foo", fs.MustGetString(localFooPath))
}

func TestExists(t *testing.T) {
	setUp()
	defer tearDown()

	exists, err := storage2.Exists(key)

	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = storage2.Exists(key + "not_exists")

	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestFiles(t *testing.T) {
	setUp()
	defer tearDown()

	files, err := storage2.Files("test/")
	assert.Nil(t, err)
	assert.Len(t, files, 1)

	var expectedSize int64 = 3
	assert.Equal(t, key, files[0].Key())
	assert.Equal(t, expectedSize, files[0].Size())

	today := time.Now().Format("2006-01-02")
	assert.Equal(t, today, files[0].LastModified().Format("2006-01-02"))
}

func TestSize(t *testing.T) {
	setUp()
	defer tearDown()

	size, err := storage2.Size(key)

	var siz int64 = 3
	assert.Nil(t, err)
	assert.Equal(t, siz, size)
}
