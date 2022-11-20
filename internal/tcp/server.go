package tcp

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"sync"
	"sync/atomic"
)

const (
	SERVER_READY                = "READY"
	SERVER_WAITING_NODES        = "SERVER_WAITING_NODES"
	SERVER_CLIENTS_STARTED_FLOW = "SERVER_STARTED_FLOW"
	SERVER_ALL_NODES_FINISH     = "SERVER_ALL_NODES_FINISH"
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
	FinishedNodes        uint32
	State                string
	reqRes               uint32
	stateCh              chan string
	registerCh           chan *Client
	unregisterCh         chan *Client
	clients              map[*Client]bool
	listener             net.Listener
	mu                   sync.Mutex
	hasFinished          chan bool
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
		hasFinished:  make(chan bool, 16),
		reqRes:       0,
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
serverLoop:
	for {
		select {
		case client := <-s.registerCh:
			s.addNode(client)
		case client := <-s.unregisterCh:
			s.removeNode(client)
		case state := <-s.stateCh:
			s.State = state
			s.handleServerState(state)
		case <-s.hasFinished:
			break serverLoop
		}
	}
	fmt.Println("Todos os processos acabaram")
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
	case CLIENT_FINISH_REQUESTS:
		s.handleClientFinishRequest(p, client)
	case REQUEST_RESPONSE:
		s.handleRequestResponse(p, client)
	}
}

func (s *TcpServer) sendMessage(packet *Packet, c *Client) {
	if err := c.encoder.Encode(packet); err != nil {
		fmt.Printf("%+v\n", packet)
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
	fmt.Println(state)
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

func (s *TcpServer) handleClientFinishRequest(p Packet, client *Client) {
	atomic.AddUint32(&s.FinishedNodes, 1)

	if s.FinishedNodes == s.Nodes {
		s.stateCh <- SERVER_ALL_NODES_FINISH
	}
}

func (s *TcpServer) handleRequestResponse(p Packet, client *Client) {
	atomic.AddUint32(&s.reqRes, 1)
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString("text to append\n"); err != nil {
		log.Println(err)
	}

	if s.reqRes == s.Requests {
		s.hasFinished <- true
	}
}
