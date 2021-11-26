package main

import (
	"log"
	"net/http"
	"net/rpc"
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
	// new a Rect service
	rect := new(Rect)
	// register Rect service
	if err := rpc.Register(rect); err != nil {
		log.Panicln(err)
	}
	// bind with http protocol
	rpc.HandleHTTP()
	// listen
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panicln(err)
	}
}
