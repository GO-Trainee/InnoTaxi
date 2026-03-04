package kafka

type KafkaHandler struct {
	// kafka consumer or consumers
	// srv service.Service - reference to the service layer to call business logic
}

func New( /*kafka consumer or consumers, srv Service*/ ) *KafkaHandler {
	return &KafkaHandler{
		// initialize kafka consumer or consumers and service reference
	}
}
