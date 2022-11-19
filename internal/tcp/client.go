package tcp

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type TcpClient struct {
	Host        string
	Requests    uint32
	ReqDone     uint32
	Concurrency uint32
	State       string
	encoder     *gob.Encoder
	decoder     *gob.Decoder
	conn        net.Conn
}

func NewTcpClient(host string) *TcpClient {
	return &TcpClient{
		Host:    host,
		encoder: &gob.Encoder{},
		decoder: &gob.Decoder{},
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
	select {}
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

func (c *TcpClient) sendMessage(p Packet) {
	c.encoder.Encode(p)
}

func (c *TcpClient) handleMessage(p Packet) {
	switch p.Action {
	case INIT_INFO:
		c.handleInitInfo(p)
	case "checkInfo":
		fmt.Printf("total request: %d concurrency: %d\n", c.Requests, c.Concurrency)
	}
}

func (c *TcpClient) handleInitInfo(p Packet) {
	initInfo := p.Payload.(InitClientInfo)
	c.Requests = initInfo.Requests
	c.Concurrency = initInfo.Concurrency
	initInfoResponse := &Packet{Action: INIT_INFO, Payload: ""}
	c.sendMessage(*initInfoResponse)
}
