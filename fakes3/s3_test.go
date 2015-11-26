package fakes3

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type FakeS3Test struct {
}

var _ = Suite(&FakeS3Test{})

func (s *FakeS3Test) TestStuff(c *C) {
	fakeS3, err := New()
	c.Assert(err, IsNil)
	defer fakeS3.Close()

	s3svc := s3.New(session.New(), fakeS3.Config)
	_, err = s3svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("mybucket"),
	})
	c.Assert(err, IsNil)
	_, err = s3svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("mybucket"),
		Key:    aws.String("path/to/my/file"),
		Body:   bytes.NewReader([]byte("Hello, World!")),
	})
	c.Assert(err, IsNil)
}
