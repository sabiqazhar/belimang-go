package client

import (
	"fmt"
	pb "github.com/sabiqazhar/belimang-go/proto/merchant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MerchantClient struct {
	Conn   *grpc.ClientConn
	Client pb.MerchantServiceClient
}

func NewMerchantClient(address string) (*MerchantClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}

	client := pb.NewMerchantServiceClient(conn)

	return &MerchantClient{
		Conn:   conn,
		Client: client,
	}, nil
}
