package grpc

import "context"

type GrpcGateway interface {
	GetUser(ctx context.Context, req gatewayentity.GetUserRequest) (gatewayentity.GetUserResponse, error)
	// other methods for other grpc calls
}
type grpcGateway struct {
	// grpc clients to other services
	// userServiceClient
	// orderServiceClient
	// etc.
}

func New( /*grpc clients*/ ) GrpcGateway {
	return &grpcGateway{}
}
