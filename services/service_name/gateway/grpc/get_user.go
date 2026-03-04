package grpc

import "context"

func (g *grpcGateway) GetUser(ctx context.Context, req gatewayentity.GetUserRequest) (gatewayentity.GetUserResponse, error) {
	//g.userServiceClient.GetUser(ctx, req)
	// return GetUserResponse{}, nil
	return GetUserResponse{}, nil
}
