package rpc

import (
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/pkg/errors"
)

func TestSession_ReadWriter(t *testing.T) {

	addr := "127.0.0.1:8080"
	my_data := "hello"
	wg := sync.WaitGroup{}
	wg.Add(2)

	// goroutine for write
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		conn, err := lis.Accept()
		if err != nil {
			t.Fatal(err)
		}
		s := NewSession(conn)
		err = s.Write([]byte(my_data))
		if err != nil {
			t.Fatal(err)
		}
	}()

	// goroutine for read
	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		s := NewSession(conn)
		data, err := s.Read()
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != my_data {
			t.Fatal(errors.New("data not equal"))
		}

		fmt.Println("data:", string(data))
	}()

	wg.Wait()
}
