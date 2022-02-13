package adapter

import (
	"github.com/sakuraapp/api/config"
	supervisorpb "github.com/sakuraapp/protobuf/supervisor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type SupervisorAdapter struct {
	supervisorpb.SupervisorServiceClient
	conn *grpc.ClientConn
}

func (a *SupervisorAdapter) Conn() *grpc.ClientConn {
	return a.conn
}

func NewSupervisorAdapter(conf *config.Config) (*SupervisorAdapter, error) {
	creds, err := credentials.NewClientTLSFromFile(conf.SupervisorKeyPath, "")

	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(conf.SupervisorAddr, opts...)

	if err != nil {
		return nil, err
	}

	client := supervisorpb.NewSupervisorServiceClient(conn)

	return &SupervisorAdapter{
		SupervisorServiceClient: client,
		conn: conn,
	}, nil
}