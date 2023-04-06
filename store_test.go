//go:build integration

package goss

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

var (
	key          = "test/foo.txt"
	fooPath      string
	localFooPath string
)

func createTempFiles() {
	f1, err := os.CreateTemp("", "goss")
	if err != nil {
		log.Fatal(err)
	}
	_, err = f1.WriteString("foo")
	if err != nil {
		log.Fatal(err)
	}
	f2, err := os.CreateTemp("", "goss")
	if err != nil {
		log.Fatal(err)
	}

	fooPath = f1.Name()
	localFooPath = f2.Name()

	_ = f2.Close()
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

type GossTestSuite struct {
	suite.Suite

	config  Config
	storage Store
	store   *store
}

func (s *GossTestSuite) SetupTest() {
	// 临时文件创建
	createTempFiles()

	// 创建 goss 对象
	goss, err := New(WithConfig(&s.config))
	if err != nil {
		s.T().Fatal(err)
	}
	s.storage = goss.Store
	s.store = s.storage.(*store)

	// 临时文件创建
	err = s.storage.PutFromFile(key, fooPath)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *GossTestSuite) SetupSuite() {
	// github 集成测试 minio 的 bucket 创建
	_, err := s.store.s3.HeadBucket(context.TODO(), &s3.HeadBucketInput{Bucket: aws.String(s.config.Bucket)})
	if err != nil && strings.Contains(err.Error(), "404") {
		_, err = s.store.s3.CreateBucket(context.TODO(), &s3.CreateBucketInput{Bucket: aws.String(s.config.Bucket)})
		if err != nil {
			s.T().Fatal(err)
		}
		s.prepareTestData()
	}

	// Files 内部一次最多获取的 key 的数量，用于测试
	maxKeys = 100
}

func (s *GossTestSuite) TearDownTest() {
	// 本地临时文件删除
	exists := fileExists(localFooPath)
	if exists {
		err := os.Remove(localFooPath)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	// 远程临时文件删除
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.store.Bucket),
		Key:    aws.String(key),
	}
	_, err := s.store.s3.DeleteObject(context.TODO(), input)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *GossTestSuite) TestPut() {
	t := s.T()

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

	err = s.storage.Put(key, f)
	assert.Nil(t, err)

	_, err = s.store.getObject(key)
	assert.Nil(t, err)
}

func (s *GossTestSuite) TestPutFromFile() {
	err := s.storage.PutFromFile(key, fooPath)
	assert.Nil(s.T(), err)

	_, err = s.store.getObject(key)
	assert.Nil(s.T(), err)
}

func (s *GossTestSuite) TestGet() {
	rc, err := s.storage.Get(key)
	assert.Nil(s.T(), err)

	bs, err := io.ReadAll(rc)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "foo", string(bs))
}

func (s *GossTestSuite) TestGetString() {
	content, err := s.storage.GetString(key)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "foo", content)
}

func (s *GossTestSuite) TestGetBytes() {
	bs, err := s.storage.GetBytes(key)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "foo", string(bs))
}

func (s *GossTestSuite) TestDelete() {
	err := s.storage.Delete(key)
	assert.Nil(s.T(), err)

	_, err = s.store.getObject(key)
	assert.NotNil(s.T(), err)
}

func (s *GossTestSuite) TestSave() {
	err := s.storage.GetToFile(key, localFooPath)
	assert.Nil(s.T(), err)

	bs, err := os.ReadFile(localFooPath)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "foo", string(bs))
}

func (s *GossTestSuite) TestExists() {
	t := s.T()

	exists, err := s.storage.Exists(key)

	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = s.storage.Exists(key + "not_exists")

	assert.Nil(t, err)
	assert.False(t, exists)
}

func (s *GossTestSuite) TestSize() {
	size, err := s.storage.Size(key)

	var siz int64 = 3
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), siz, size)
}

func (s *GossTestSuite) TestFiles() {
	t := s.T()

	files, err := s.storage.Files("test/")
	assert.Nil(t, err)
	assert.Len(t, files, 1)

	var expectedSize int64 = 3
	assert.Equal(t, key, files[0].Key())
	assert.Equal(t, expectedSize, files[0].Size())

	today := time.Now().Format("2006-01-02")
	assert.Equal(t, today, files[0].LastModified().Format("2006-01-02"))
}

// prepare test data for testFilesWithMultiPage
func (s *GossTestSuite) prepareTestData() {
	dir := "test_all/"

	for i := 1; i <= 200; i++ {
		err := s.storage.Put(fmt.Sprintf("%s%s.txt", dir, strconv.Itoa(i)), strings.NewReader("foo"))
		assert.Nil(s.T(), err)
	}
}

func (s *GossTestSuite) TestFilesWithMultiPage() {
	// Testdata was prepared before.
	dir := "test_all/"

	files, err := s.storage.Files(dir)
	assert.Nil(s.T(), err)
	assert.Len(s.T(), files, 200)
}

func configs() map[string]Config {
	type yml struct {
		Configs map[string]Config `yaml:"configs"`
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	bs, err := os.ReadFile(filepath.Join(usr.HomeDir, ".goss.yml"))
	if err != nil {
		log.Fatal(err)
	}

	c := &yml{}
	err = yaml.NewDecoder(strings.NewReader(string(bs))).Decode(c)
	if err != nil {
		log.Fatal(err)
	}

	return c.Configs
}

func TestGossTestSuite(t *testing.T) {
	configs := configs()

	for k, val := range configs {
		t.Run(k, func(t *testing.T) {
			suite.Run(t, &GossTestSuite{config: val})
		})
	}
}
