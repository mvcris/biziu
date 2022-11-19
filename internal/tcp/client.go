package tcp

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"net"
	"sync/atomic"
)

type TcpClient struct {
	Host         string
	Requests     uint32
	ExecRequests uint32
	ReqDone      uint32
	Concurrency  uint32
	ReqPerGo     uint32
	ReqPerGoRem  uint32
	encoder      *gob.Encoder
	decoder      *gob.Decoder
	conn         net.Conn
	hasFinished  chan bool
}

func NewTcpClient(host string) *TcpClient {
	return &TcpClient{
		Host:         host,
		encoder:      &gob.Encoder{},
		decoder:      &gob.Decoder{},
		hasFinished:  make(chan bool, 1),
		ExecRequests: 0,
	}
}

func (c *TcpClient) Start() {
	gob.Register(InitClientInfo{})
	fmt.Println("start client")
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal("Connection error", err)
	}

	c.conn = conn
	c.encoder = gob.NewEncoder(conn)
	c.decoder = gob.NewDecoder(conn)
	go c.readMessages()
clientLoop:
	for {
		select {
		case <-c.hasFinished:
			c.sendMessage(Packet{Action: CLIENT_FINISH_REQUESTS})
			break clientLoop
		}
	}
	fmt.Println("node finalizou")
}

func (c *TcpClient) readMessages() {
	defer c.conn.Close()
	for {
		var p Packet
		if err := c.decoder.Decode(&p); err != nil {
			panic(err)
		}
		c.handleMessage(p)
	}
}

func (c *TcpClient) slipRequests() {
	c.ReqPerGo = uint32(math.Floor(float64(c.Requests) / float64(c.Concurrency)))
	c.ReqPerGoRem = c.Requests % c.Concurrency
}

func (c *TcpClient) sendMessage(p Packet) {
	c.encoder.Encode(p)
}

func (c *TcpClient) handleMessage(p Packet) {
	switch p.Action {
	case INIT_INFO:
		c.handleInitInfo(p)
	case START_REQUESTS:
		c.startRequests(p)
	}
}

func (c *TcpClient) handleInitInfo(p Packet) {
	initInfo := p.Payload.(InitClientInfo)
	c.Requests = initInfo.Requests
	c.Concurrency = initInfo.Concurrency
	c.slipRequests()
	initInfoResponse := &Packet{Action: INIT_INFO, Payload: ""}
	c.sendMessage(*initInfoResponse)
}

func (c *TcpClient) startRequests(p Packet) {
	for i := 0; i < int(c.Concurrency); i++ {
		for n := 0; n < int(c.ReqPerGo); n++ {
			go c.execute()
			if n == int(c.Concurrency) {
				if c.ReqPerGoRem > 0 {
					for m := 0; m < int(c.ReqPerGoRem); m++ {
						go c.execute()
					}
				}
			}
		}
	}

}

func (c *TcpClient) execute() {
	atomic.AddUint32(&c.ExecRequests, 1)
	fmt.Printf("total req: %d total exec: %d\n", c.Requests, c.ExecRequests)
	if c.Requests == c.ExecRequests {
		c.hasFinished <- true
	}
}
