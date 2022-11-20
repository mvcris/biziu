package tcp

const (
	INIT_INFO              = "INIT_INFO"
	START_REQUESTS         = "START_REQUESTS"
	CLIENT_FINISH_REQUESTS = "CLIENT_FINISH_REQUESTS"
	REQUEST_RESPONSE       = "REQUEST_RESPONSE"
	CLOSE_NODE_CONNECTION  = "CLOSE_NODE_CONNECTION"
)

type Packet struct {
	Action  string
	Payload interface{}
}

type InitClientInfo struct {
	Requests    uint32
	Concurrency uint32
	Url         string            `json:"url"`
	Method      string            `json:"method"`
	Header      map[string]string `json:"headers"`
	Body        map[string]any    `json:"body"`
}
