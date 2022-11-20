package tcp

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type TcpClient struct {
	Host         string
	Requests     uint32
	ExecRequests uint32
	ReqDone      uint32
	Concurrency  uint32
	ReqLoopTimes uint32
	ReqLoopRem   uint32
	encoder      *gob.Encoder
	decoder      *gob.Decoder
	conn         net.Conn
	hasFinished  chan bool
	wg           sync.WaitGroup
}

func NewTcpClient(host string) *TcpClient {
	return &TcpClient{
		Host:         host,
		encoder:      &gob.Encoder{},
		decoder:      &gob.Decoder{},
		hasFinished:  make(chan bool, 1),
		ExecRequests: 0,
		wg:           sync.WaitGroup{},
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

	<-c.hasFinished
	c.sendMessage(Packet{Action: CLIENT_FINISH_REQUESTS})
	c.wg.Wait()
	fmt.Println("node finalizou")
}

func (c *TcpClient) readMessages() {
	for {
		var p Packet
		fmt.Println(p.Action)
		if err := c.decoder.Decode(&p); err != nil {
			panic(err)
		}
		c.handleMessage(p)
	}
}

func (c *TcpClient) slipRequests() {
	c.ReqLoopTimes = uint32(math.Floor(float64(c.Requests) / float64(c.Concurrency)))
	c.ReqLoopRem = c.Requests % c.Concurrency
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
	fmt.Printf("loop times: %d", c.ReqLoopTimes)
	count := 0
	for i := 0; i < int(c.ReqLoopTimes); i++ {
		count++
		remainder := false
		if count == int(c.ReqLoopTimes) {
			remainder = true
		}
		c.startConcurrencyRequest(p, remainder)
	}
}

func (c *TcpClient) startConcurrencyRequest(p Packet, remainder bool) {
	total := c.Concurrency
	if remainder {
		total += c.ReqLoopRem
	}
	for i := 0; i < int(total); i++ {
		c.wg.Add(1)
		go c.execute()
	}
}

func (c *TcpClient) execute() {
	atomic.AddUint32(&c.ExecRequests, 1)
	p := Packet{Action: REQUEST_RESPONSE, Payload: time.Now().UTC().Unix()}
	c.sendMessage(p)
	c.wg.Done()
	if c.Requests == c.ExecRequests {
		c.hasFinished <- true
	}
}
