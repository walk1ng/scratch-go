package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Params struct {
	Height, Width int
}

type Rect struct{}

func (r *Rect) Area(p Params, ret *int) error {
	*ret = p.Height * p.Width
	return nil
}

func (r *Rect) Perimeter(p Params, ret *int) error {
	*ret = (p.Height + p.Width) * 2
	return nil
}

func main() {
	rpc.Register(new(Rect))
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panicln(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}

		go func(conn net.Conn) {
			log.Println("new client connection")
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}
