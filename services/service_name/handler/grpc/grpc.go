package grpc

type GrpcHandler struct {
	// grpc consumer or consumers
	// srv service.Service - reference to the service layer to call business logic
}

func New( /*grpc consumer or consumers, srv Service*/ ) *GrpcHandler {
	return &GrpcHandler{
		// initialize grpc consumer or consumers and service reference
	}
}
