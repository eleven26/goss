package goss

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config := Config{}
	assert.NotNil(t, config.validate())

	config = Config{
		Endpoint: "example.com",
		UseSsl:   aws.Bool(false),
	}
	assert.Equal(t, "http://example.com", config.url())
}
