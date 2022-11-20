package tcp

// Basic imports
import (
	"encoding/gob"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ClientTcpTestSuite struct {
	suite.Suite
}

func (s *ClientTcpTestSuite) TestInitRequest() {
	go func() {
		net, _ := net.Listen("tcp", "localhost:3000")

		conn, _ := net.Accept()
		encoder := gob.NewEncoder(conn)
		initPacket := &InitClientInfo{Requests: 10, Concurrency: 1}
		encoder.Encode(&Packet{Action: INIT_INFO, Payload: initPacket})
		encoder.Encode(&Packet{Action: START_REQUESTS, Payload: initPacket})
	}()

	client := NewTcpClient("localhost:3000")
	go client.Start()
	time.Sleep(time.Second * 3)
	assert.IsType(s.T(), client.encoder, &gob.Encoder{})
	//assert.Panics(s.T(), func() { client.Start() }, "The code did not panic")
	assert.IsType(s.T(), client.decoder, &gob.Decoder{})
	assert.Equal(s.T(), client.Concurrency, uint32(1))
	assert.Equal(s.T(), client.ReqLoopTimes, uint32(10))
	assert.Equal(s.T(), client.ReqLoopRem, uint32(0))
	assert.Equal(s.T(), client.ExecRequests, uint32(10))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTcpTestSuite))
}
