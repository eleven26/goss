//go:build integration

package goss

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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

var (
	storage Store

	s *store

	key          = "test/foo.txt"
	testdata     string
	fooPath      string
	localFooPath string
)

func init() {
	goss, err := New(WithConfig(&Config{
		Endpoint:          os.Getenv("GOSS_ENDPOINT"),
		AccessKey:         os.Getenv("GOSS_ACCESS_KEY"),
		SecretKey:         os.Getenv("GOSS_SECRET_KEY"),
		Region:            os.Getenv("GOSS_REGION"),
		Bucket:            os.Getenv("GOSS_BUCKET"),
		UseSsl:            aws.Bool(true),
		HostnameImmutable: aws.Bool(false),
	}))
	if err != nil {
		log.Fatal(err)
	}

	storage = goss.Store

	s = storage.(*store)

	testdata = filepath.Join("testdata")
	fooPath = filepath.Join(testdata, "foo.txt")
	localFooPath = filepath.Join(testdata, "foo1.txt")
}

func setUp(t *testing.T) {
	err := storage.PutFromFile(key, fooPath)
	if err != nil {
		t.Fatal(err)
	}
}

func tearDown(t *testing.T) {
	deleteLocal(t)
	deleteRemote(t)
}

func deleteRemote(t *testing.T) {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	_, err := s.s3.DeleteObject(context.TODO(), input)
	if err != nil {
		t.Fatal(err)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func deleteLocal(t *testing.T) {
	exists := fileExists(localFooPath)
	if exists {
		err := os.Remove(localFooPath)
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

	err = storage.Put(key, f)
	assert.Nil(t, err)

	_, err = s.getObject(key)
	assert.Nil(t, err)
}

func TestPutFromFile(t *testing.T) {
	defer tearDown(t)

	err := storage.PutFromFile(key, fooPath)
	assert.Nil(t, err)

	_, err = s.getObject(key)
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	rc, err := storage.Get(key)
	assert.Nil(t, err)

	bs, err := io.ReadAll(rc)
	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}

func TestGetString(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	content, err := storage.GetString(key)

	assert.Nil(t, err)
	assert.Equal(t, "foo", content)
}

func TestGetBytes(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	bs, err := storage.GetBytes(key)

	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}

func TestDelete(t *testing.T) {
	setUp(t)
	defer deleteLocal(t)

	err := storage.Delete(key)
	assert.Nil(t, err)

	_, err = s.getObject(key)
	assert.NotNil(t, err)
}

func TestSave(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	err := storage.GetToFile(key, localFooPath)
	assert.Nil(t, err)

	bs, err := os.ReadFile(localFooPath)
	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bs))
}

func TestExists(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	exists, err := storage.Exists(key)

	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = storage.Exists(key + "not_exists")

	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestSize(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	size, err := storage.Size(key)

	var siz int64 = 3
	assert.Nil(t, err)
	assert.Equal(t, siz, size)
}

func TestFiles(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	files, err := storage.Files("test/")
	assert.Nil(t, err)
	assert.Len(t, files, 1)

	var expectedSize int64 = 3
	assert.Equal(t, key, files[0].Key())
	assert.Equal(t, expectedSize, files[0].Size())

	today := time.Now().Format("2006-01-02")
	assert.Equal(t, today, files[0].LastModified().Format("2006-01-02"))
}

func sTestAb(t *testing.T) {
	dir := "test_all/"

	for i := 1; i <= 200; i++ {
		err := storage.Put(fmt.Sprintf("%s%s.txt", dir, strconv.Itoa(i)), strings.NewReader("foo"))
		assert.Nil(t, err)
	}
}

func TestFilesWithMultiPage(t *testing.T) {
	// Testdata was prepared before.
	dir := "test_all/"

	files, err := storage.Files(dir)
	assert.Nil(t, err)
	assert.Len(t, files, 200)
}
