package server

import (
	"bufio"
	"fmt"
	"net"
	"testing"

	"github.com/pkar/blab/check"
)

func getFreePort() int {
	for i := 5000; i < 65535; i++ {
		conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", i))
		if err != nil {
			continue
		}
		conn.Close()
		return i
	}
	return 0
}

func TestNew(t *testing.T) {
	s, err := New("127.0.0.1", getFreePort(), "")
	check.Ok(t, err)
	go func() {
		<-s.shutdown
		s.done <- struct{}{}
	}()
	s.Close()
}

func TestStart(t *testing.T) {
	s, err := New("127.0.0.1", getFreePort(), "")
	check.Ok(t, err)
	go s.Start()

	// create a client which then goes to handleConn
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", s.port))
	check.Ok(t, err)
	status, err := bufio.NewReader(conn).ReadString('\n')
	check.Ok(t, err)
	check.Equals(t, "\n", status)
	conn.Write([]byte("bob\n"))

	s.Close()
}
