package rpc

import (
	"fmt"
	"log"
	"net"
	"reflect"
)

type Server struct {
	addr  string
	funcs map[string]reflect.Value
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

func (s *Server) Register(rpcName string, f interface{}) {
	if _, ok := s.funcs[rpcName]; ok {
		return
	}

	fVal := reflect.ValueOf(f)
	s.funcs[rpcName] = fVal
}

func (s *Server) Run() {
	fmt.Println("Run server at", s.addr)
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("监听 %s err :%v", s.addr, err)
		return
	}

	for {
		fmt.Println("listening the incoming requests...")
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("listener accept failed: %v\n", err)
			return
		}
		session := NewSession(conn)
		b, err := session.Read()
		if err != nil {
			log.Printf("read rpc data failed: %v\n", err)
			return
		}

		rpcData, err := decode(b)
		if err != nil {
			log.Printf("decode rpc data failed: %v\n", err)
			return
		}

		f, ok := s.funcs[rpcData.Name]
		if !ok {
			log.Printf("func %s not found\n", rpcData.Name)
			return
		}

		// handle func args
		inArgs := make([]reflect.Value, 0, len(rpcData.Args))
		for _, arg := range rpcData.Args {
			inArgs = append(inArgs, reflect.ValueOf(arg))
		}

		out := f.Call(inArgs)
		outArgs := make([]interface{}, len(out))
		for _, arg := range out {
			outArgs = append(outArgs, arg.Interface())
		}

		respData := RPCData{Name: rpcData.Name, Args: outArgs}
		b, err = encode(respData)
		if err != nil {
			log.Printf("encode rpc data failed: %v\n", err)
			return
		}

		err = session.Write(b)
		if err != nil {
			log.Printf("write rpc data failed: %v\n", err)
			return
		}
	}
}
