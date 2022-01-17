package store

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"mime/multipart"
	"net/url"
)

type S3Config struct {
	Region *string
	Bucket *string
	Endpoint *string
	ForcePathStyle *bool
	Credentials *credentials.Credentials
}

type S3Adapter struct {
	baseUrl string
	conf *S3Config
	svc *s3.S3
	uploader *s3manager.Uploader
}

func (a *S3Adapter) Upload(key string, file multipart.File) (string, error) {
	res, err := a.uploader.Upload(&s3manager.UploadInput{
		Bucket: a.conf.Bucket,
		Key: &key,
		Body: file,
	})

	if err != nil {
		return "", err
	}

	return res.Location, nil
}

func (a *S3Adapter) Delete(key string) error {
	_, err := a.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: a.conf.Bucket,
		Key: &key,
	})

	return err
}

func (a *S3Adapter) ResolveURL(key string) string {
	return fmt.Sprintf("%s/%s", a.baseUrl, key)
}

func NewS3Adapter(conf *S3Config) *S3Adapter {
	awsConf := &aws.Config{
		Region: conf.Region,
		Credentials: conf.Credentials,
		Endpoint: conf.Endpoint,
		S3ForcePathStyle: conf.ForcePathStyle,
	}

	sess := session.Must(session.NewSession(awsConf))
	svc := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	baseUrl := svc.Endpoint

	if *conf.ForcePathStyle {
		baseUrl += *conf.Bucket
	} else {
		u, err := url.Parse(baseUrl)

		if err != nil {
			panic(err)
		}

		u.Host = fmt.Sprintf("%s.%s", *conf.Bucket, u.Host)
		baseUrl = u.String()
	}

	return &S3Adapter{
		baseUrl: baseUrl,
		conf: conf,
		svc: svc,
		uploader: uploader,
	}
}