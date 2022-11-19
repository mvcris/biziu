package internal

type Biziu struct {
	Requests    uint32
	Concurrency uint32
	Nodes       uint16
	Host        string
	port        string
}

func NewBiziuServerMode(requests uint32, concurrency uint32, nodes uint16, host string) *Biziu {
	return &Biziu{
		Requests:    requests,
		Concurrency: concurrency,
		Nodes:       nodes,
		Host:        host,
	}
}

func NewBiziuClientMode(host string) *Biziu {
	return &Biziu{
		Host: host,
	}
}

func (b *Biziu) ExecuteClient() {}
