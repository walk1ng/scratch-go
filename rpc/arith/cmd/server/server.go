package main

import (
	"log"
	"net/http"
	"net/rpc"

	"github.com/pkg/errors"
)

type Arith struct{}

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

func (a *Arith) Multiply(req ArithRequest, resp *ArithResponse) error {
	resp.Pro = req.A * req.B
	return nil
}

func (a *Arith) Divide(req ArithRequest, resp *ArithResponse) error {
	if req.B == 0 {
		return errors.New("除数不能为0")
	}

	// 除
	resp.Quo = req.A / req.B
	// 取模
	resp.Rem = req.A % req.B
	return nil
}

func main() {
	rpc.Register(new(Arith))
	rpc.HandleHTTP()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panicln(err)
	}

}
