package gateway

import (
	"awesomeProject/services/service_name/gateway/grpc"
	"awesomeProject/services/service_name/gateway/http"
	"awesomeProject/services/service_name/gateway/kafka"
)

type Gateway struct {
	HttpGateway  http.HttpGateway
	KafkaGateway kafka.KafkaGateway
	GrpcGateway  grpc.GrpcGateway
}

func New( /*some gateways*/ ) Gateway {
	return &Gateway{
		// HttpGateway:  httpGateway,
		// KafkaGateway: kafkaGateway,
		// GrpcGateway:  grpcGateway,
	}
}
