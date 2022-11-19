package tcp

import (
	"encoding/gob"
	"fmt"
	"math"
	"net"
	"sync"
	"sync/atomic"
)

const (
	SERVER_READY                = "READY"
	SERVER_WAITING_NODES        = "SERVER_WAITING_NODES"
	SERVER_CLIENTS_STARTED_FLOW = "SERVER_STARTED_FLOW"
	SERVER_FINISH               = "SERVER_FINISH"
	SERVER_NODES_READY          = "SERVER_NODES_READY"
)

type Client struct {
	conn    net.Conn
	encoder *gob.Encoder
	decoder *gob.Decoder
}

type TcpServer struct {
	Requests             uint32
	Concurrency          uint32
	Nodes                uint32
	Port                 uint16
	ReqPerNode           uint32
	ReqDivisionRemainder uint32
	ConnectedNodes       uint32
	ReadyNodes           uint32
	FinishedNodes        uint16
	State                string
	stateCh              chan string
	registerCh           chan *Client
	unregisterCh         chan *Client
	clients              map[*Client]bool
	listener             net.Listener
	mu                   sync.Mutex
}

func NewTcpServer(requests uint32, concurrency uint32, nodes uint32, port uint16) *TcpServer {
	return &TcpServer{
		Requests:     requests,
		Concurrency:  concurrency,
		Nodes:        nodes,
		Port:         port,
		registerCh:   make(chan *Client),
		unregisterCh: make(chan *Client),
		clients:      make(map[*Client]bool),
		mu:           sync.Mutex{},
		State:        SERVER_WAITING_NODES,
		ReadyNodes:   0,
		stateCh:      make(chan string, 256),
	}
}

func (s *TcpServer) Start() {
	gob.Register(InitClientInfo{})
	port := fmt.Sprintf(":%d", s.Port)
	listener, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println("Error on creating server")
		panic(err)
	}
	s.splitRequests()
	s.listener = listener
	go s.Listen()
	for {
		select {
		case client := <-s.registerCh:
			s.addNode(client)
		case client := <-s.unregisterCh:
			s.removeNode(client)
		case state := <-s.stateCh:
			s.State = state
			s.handleServerState(state)
		}
	}

}

func (s *TcpServer) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("error on client connect: %v", err)
		}
		client := &Client{
			encoder: gob.NewEncoder(conn),
			decoder: gob.NewDecoder(conn),
			conn:    conn,
		}
		s.registerCh <- client
		go s.handleConnection(client)
	}
}

func (s *TcpServer) handleConnection(c *Client) {
	for {
		var p Packet
		if err := c.decoder.Decode(&p); err != nil {
			fmt.Println("erro ao receber mensagem")
			fmt.Println(err)
			s.unregisterCh <- c
			return
		}
		s.handleMessage(p, c)
	}
}

func (s *TcpServer) handleMessage(p Packet, client *Client) {
	switch p.Action {
	case INIT_INFO:
		s.handleInitInfo(p, client)
	}
}

func (s *TcpServer) sendMessage(packet *Packet, c *Client) {
	if err := c.encoder.Encode(packet); err != nil {
		fmt.Println(err)
		fmt.Println("erro ao enviar mensagem")
	}
}

func (s *TcpServer) addNode(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.clients) >= int(s.Nodes) {
		fmt.Println("pool is full")
		s.unregisterCh <- client
		return
	}

	s.clients[client] = true
	atomic.AddUint32(&s.ConnectedNodes, 1)

	requests := s.ReqPerNode
	if len(s.clients) == int(s.Nodes) {
		requests += s.ReqDivisionRemainder
	}
	initPacket := &InitClientInfo{Requests: requests, Concurrency: s.Concurrency}
	s.sendMessage(&Packet{Action: INIT_INFO, Payload: initPacket}, client)
}

func (s *TcpServer) splitRequests() {
	s.ReqPerNode = uint32(math.Floor(float64(s.Requests) / float64(s.Nodes)))
	s.ReqDivisionRemainder = s.Requests % s.Nodes
}

func (s *TcpServer) removeNode(client *Client) {
	if ok := s.clients[client]; ok {
		delete(s.clients, client)
		s.ConnectedNodes--
		fmt.Println("node removido")
	}

	//@NOTE: do something??? maybe when race condition ocurred?
}

func (s *TcpServer) handleInitInfo(p Packet, client *Client) {
	atomic.AddUint32(&s.ReadyNodes, 1)
	if s.ReadyNodes == s.Nodes {
		s.stateCh <- SERVER_NODES_READY
	}
}

func (s *TcpServer) handleServerState(state string) {
	switch state {
	case SERVER_NODES_READY:
		p := &Packet{Action: START_REQUESTS}
		for client := range s.clients {
			s.sendMessage(p, client)
		}
		s.stateCh <- SERVER_CLIENTS_STARTED_FLOW
	case SERVER_CLIENTS_STARTED_FLOW:
		fmt.Println("clientes inciaram")
	}
}
