package adapter

import (
	"github.com/sakuraapp/api/internal/config"
	sharedUtil "github.com/sakuraapp/shared/pkg/util"
)

type Adapters struct {
	S3 *S3Adapter
	Supervisor *SupervisorAdapter
}

func Init(conf *config.Config) (*Adapters, error) {
	var err error

	a := &Adapters{}

	a.S3 = NewS3Adapter(&sharedUtil.S3Config{
		Bucket: conf.S3Bucket,
		Region: conf.S3Region,
		Endpoint: conf.S3Endpoint,
		ForcePathStyle: conf.S3ForcePathStyle,
	})

	a.Supervisor, err = NewSupervisorAdapter(conf)

	if err != nil {
		return nil, err
	}

	return a, nil
}