package store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	sharedUtil "github.com/sakuraapp/shared/pkg/util"
	"mime/multipart"
)

type S3Adapter struct {
	baseUrl string
	conf *sharedUtil.S3Config
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
	return sharedUtil.ResolveS3URL(a.baseUrl, key)
}

func NewS3Adapter(conf *sharedUtil.S3Config) *S3Adapter {
	awsConf := &aws.Config{
		Region: conf.Region,
		Endpoint: conf.Endpoint,
		S3ForcePathStyle: conf.ForcePathStyle,
	}

	sess := session.Must(session.NewSession(awsConf))
	svc := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	return &S3Adapter{
		baseUrl: sharedUtil.GetS3BaseUrl(conf),
		conf: conf,
		svc: svc,
		uploader: uploader,
	}
}