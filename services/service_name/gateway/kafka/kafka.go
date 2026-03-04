package kafka

type KafkaGateway interface {
	// methods for sending messages to kafka topics
	UserAdded(req gatewayentity.UserAddedRequest) error
}
type kafkaGateway struct {
	// need kafka producer or producers
}

func New() KafkaGateway {
	return &kafkaGateway{}
}
