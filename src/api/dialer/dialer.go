package dialer

import (
	pb "github.com/fahrizalfarid/user-service-grpc/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	dialer struct{}
	Dialer interface {
		UserValidatorClient(addr string) (*grpc.ClientConn, pb.UserValidatorClient, error)
		UserClient(addr string) (*grpc.ClientConn, pb.UserClient, error)
	}
)

func NewGrpcDialer() Dialer {
	return &dialer{}
}

func (d *dialer) UserValidatorClient(addr string) (*grpc.ClientConn, pb.UserValidatorClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	u := pb.NewUserValidatorClient(conn)
	return conn, u, nil
}

func (d *dialer) UserClient(addr string) (*grpc.ClientConn, pb.UserClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	u := pb.NewUserClient(conn)
	return conn, u, nil
}
