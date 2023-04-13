package goss

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type (
	PutOption func(input *s3.PutObjectInput)
)

func withPutOptions(o *s3.PutObjectInput, options ...PutOption) *s3.PutObjectInput {
	if o == nil {
		return nil
	}

	for _, with := range options {
		with(o)
	}

	return o
}

func WithBucketPutOption(b string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.Bucket = aws.String(b)
	}
}

func WithKeyPutOption(k string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.Key = aws.String(k)
	}
}

func WithACLPutOption(a types.ObjectCannedACL) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ACL = a
	}
}

func WithBodyPutOption(b io.Reader) PutOption {
	return func(input *s3.PutObjectInput) {
		input.Body = b
	}
}

func WithBucketKeyEnabledPutOption(b bool) PutOption {
	return func(input *s3.PutObjectInput) {
		input.BucketKeyEnabled = b
	}
}

func WithCacheControlPutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.CacheControl = aws.String(c)
	}
}

func withChecksumAlgorithmPutOption(a types.ChecksumAlgorithm) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ChecksumAlgorithm = a
	}
}

func WithChecksumCRC32PutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ChecksumCRC32 = aws.String(c)
	}
}

func WithChecksumCRC32CPutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ChecksumCRC32C = aws.String(c)
	}
}

func WithChecksumSHA1PutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ChecksumSHA1 = aws.String(c)
	}
}

func WithChecksumSHA256PutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ChecksumSHA256 = aws.String(c)
	}
}

func WithContentDispositionPutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ContentDisposition = aws.String(c)
	}
}

func WithContentDispositionIutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ContentDisposition = aws.String(c)
	}
}

func WithContentEncodingPutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ContentEncoding = aws.String(c)
	}
}

func WithContentLanguagePutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ContentLanguage = aws.String(c)
	}
}

func WithContentLengthPutOption(l int64) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ContentLength = l
	}
}

func withContentMD5PutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ContentMD5 = aws.String(c)
	}
}

func WithContentTypePutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ContentType = aws.String(c)
	}
}

func WithExpectedBucketOwnerPutOption(e string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ExpectedBucketOwner = aws.String(e)
	}
}

func WithExpiresPutOption(t time.Time) PutOption {
	return func(input *s3.PutObjectInput) {
		input.Expires = aws.Time(t)
	}
}

func WithGrantFullControlPutOption(g string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.GrantFullControl = aws.String(g)
	}
}

func WithGrantReadPutOption(g string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.GrantRead = aws.String(g)
	}
}

func WithGrantReadACPPutOption(g string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.GrantReadACP = aws.String(g)
	}

}

func WithGrantWriteACPPutOption(g string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.GrantWriteACP = aws.String(g)
	}
}

func WithMetadataPutOption(m map[string]string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.Metadata = m
	}
}
func WithObjectLockLegalHoldStatusPutOption(s types.ObjectLockLegalHoldStatus) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ObjectLockLegalHoldStatus = s
	}
}

func WithObjectLockModePutOption(m types.ObjectLockMode) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ObjectLockMode = m
	}
}

func WithObjectLockRetainUntilDatePutOption(t time.Time) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ObjectLockRetainUntilDate = aws.Time(t)
	}
}

func WithRequestPayerPutOption(p types.RequestPayer) PutOption {
	return func(input *s3.PutObjectInput) {
		input.RequestPayer = p
	}
}

func WithSSECustomerAlgorithmPutOption(a string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.SSECustomerAlgorithm = aws.String(a)
	}
}

func WithSSECustomerKeyPutOption(k string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.SSECustomerKey = aws.String(k)
	}
}

func WithSSECustomerKeyMD5PutOption(k string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.SSECustomerKeyMD5 = aws.String(k)
	}
}

func WithSSEKMSEncryptionContextPutOption(c string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.SSEKMSEncryptionContext = aws.String(c)
	}
}

func WithSSEKMSKeyIdPutOption(k string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.SSEKMSKeyId = aws.String(k)
	}
}

func WithServerSideEncryptionPutOption(e types.ServerSideEncryption) PutOption {
	return func(input *s3.PutObjectInput) {
		input.ServerSideEncryption = e
	}
}

func WithStorageClassPutOption(s types.StorageClass) PutOption {
	return func(input *s3.PutObjectInput) {
		input.StorageClass = s
	}
}

func WithTaggingPutOption(t string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.Tagging = aws.String(t)
	}
}

func WithWebsiteRedirectLocationPutOption(l string) PutOption {
	return func(input *s3.PutObjectInput) {
		input.WebsiteRedirectLocation = aws.String(l)
	}
}
