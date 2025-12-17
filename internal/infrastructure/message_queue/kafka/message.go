package kafka

type Message struct {
	Topic     string
	Key       string
	Value     []byte
	Headers   map[string]string
	Partition int32
}

type DeliveryReport struct {
	Topic     string
	Key       string
	Partition int32
	Offset    int64
	Error     error
}
