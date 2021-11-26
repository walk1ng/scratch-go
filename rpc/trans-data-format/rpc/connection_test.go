package rpc

import (
	"encoding/gob"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/pkg/errors"
)

type User struct {
	Name string
	Age  int
}

func queryUser(id int) (User, error) {
	db := make(map[int]User)
	db[0] = User{"wei", 100}
	db[1] = User{"fei", 100}
	db[2] = User{"kei", 98}

	if u, ok := db[id]; ok {
		return u, nil
	}

	return User{}, errors.New(fmt.Sprintf("user id %d not found", id))
}

func Test_RPC(t *testing.T) {
	gob.Register(User{})
	addr := "127.0.0.1:8080"
	srv := NewServer(addr)
	srv.Register("queryUser", queryUser)
	go srv.Run()

	time.Sleep(time.Second * 2)

	fmt.Println("dial server")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	cli := NewClient(conn)
	var query func(int) (User, error)
	cli.callRPC("queryUser", &query)
	u, err := query(1)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("user: %+v\n", u)
	time.Sleep(time.Second * 5)
}
