package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type ArithRequest struct {
	A, B int
}

type ArithResponse struct {
	// 乘积
	Pro int
	// 商
	Quo int
	// 余数
	Rem int
}

func main() {
	conn, err := rpc.DialHTTP("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}

	var res ArithResponse
	req := ArithRequest{11, 3}
	err = conn.Call("Arith.Multiply", req, &res)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)

	err = conn.Call("Arith.Divide", req, &res)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d / %d 商 %d，余数 = %d\n", req.A, req.B, res.Quo, res.Rem)

}
