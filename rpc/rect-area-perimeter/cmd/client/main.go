package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type Params struct {
	Height, Width int
}

func main() {
	conn, err := rpc.DialHTTP("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}

	ret := 0
	err = conn.Call("Rect.Area", Params{5, 10}, &ret)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("area: %d\n", ret)

	err = conn.Call("Rect.Perimeter", Params{5, 10}, &ret)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("perimeter: %d\n", ret)
}
