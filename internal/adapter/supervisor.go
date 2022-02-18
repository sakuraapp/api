package adapter

import (
	"github.com/sakuraapp/api/internal/config"
	supervisorpb "github.com/sakuraapp/protobuf/supervisor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
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

	resolver.SetDefaultScheme("dns")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ] }`), // use round-robin LB
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