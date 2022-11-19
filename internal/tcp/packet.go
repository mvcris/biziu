package tcp

const (
	INIT_INFO              = "INIT_INFO"
	START_REQUESTS         = "START_REQUESTS"
	CLIENT_FINISH_REQUESTS = "CLIENT_FINISH_REQUESTS"
)

type Packet struct {
	Action  string
	Payload interface{}
}

type InitClientInfo struct {
	Requests    uint32
	Concurrency uint32
}
