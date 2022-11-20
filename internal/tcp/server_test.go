package tcp

// Basic imports
import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerTcpTestSuite struct {
	suite.Suite
}

func (s *ServerTcpTestSuite) TestServerReques() {
	go func() {
		time.Sleep(time.Second * 1)
		client := NewTcpClient(":3001")
		go client.Start()
		time.Sleep(time.Second * 3)
	}()
	server := NewTcpServer(10, 1, 1, 3001, "../../example.json")

	go server.Start()
	time.Sleep(3 * time.Second)
	assert.Equal(s.T(), len(server.clients), 1)

}

func (s *ServerTcpTestSuite) TestServerRequesFullPool() {
	go func() {
		time.Sleep(time.Second * 2)
		client := NewTcpClient(":3002")
		go client.Start()
		client2 := NewTcpClient(":3002")
		go client2.Start()
		time.Sleep(time.Second * 1)
	}()
	server := NewTcpServer(10, 1, 1, 3002, "../../example.json")

	go server.Start()
	time.Sleep(2 * time.Second)
}

func (s *ServerTcpTestSuite) TestServerDropNode() {
	go func() {
		time.Sleep(time.Second * 1)
		client := NewTcpClient(":3003")
		go client.Start()
		time.Sleep(time.Second * 1)
		client.conn.Close()
	}()

	server := NewTcpServer(10, 1, 1, 3003, "../../example.json")
	go server.Start()
	time.Sleep(2 * time.Second)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServerTcpTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTcpTestSuite))
}
