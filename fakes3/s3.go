package fakes3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/goamz/goamz/s3/s3test"
)

type FakeS3 struct {
	Server *s3test.Server
	Config *aws.Config
}

func New() (*FakeS3, error) {
	var err error
	f := FakeS3{}
	f.Server, err = s3test.NewServer(&s3test.Config{})
	if err != nil {
		return nil, err
	}
	f.Config = &aws.Config{
		Endpoint:         aws.String(f.Server.URL()),
		Region:           aws.String("fake-region"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	return &f, nil
}

func (f *FakeS3) Close() {
	f.Server.Quit()
}
