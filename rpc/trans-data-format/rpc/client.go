package rpc

import (
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) callRPC(rpcName string, fPtr interface{}) {
	fn := reflect.ValueOf(fPtr).Elem()

	//
	f := func(inputArgs []reflect.Value) []reflect.Value {
		inArgs := make([]interface{}, len(inputArgs))
		for _, arg := range inputArgs {
			inArgs = append(inArgs, arg.Interface())
		}

		cliSession := NewSession(c.conn)
		reqRPC := RPCData{Name: rpcName, Args: inArgs}
		b, err := encode(reqRPC)
		if err != nil {
			panic(err)
		}
		// write rpc data to server
		err = cliSession.Write(b)
		if err != nil {
			panic(err)
		}

		// receive and read the responsed rpc data
		b, err = cliSession.Read()
		if err != nil {
			panic(err)
		}
		respRPC, err := decode(b)
		if err != nil {
			panic(err)
		}

		outArgs := make([]reflect.Value, len(respRPC.Args))
		for i, arg := range respRPC.Args {
			if arg == nil {
				// reflect.Zero()会返回类型的零值的value
				// .out()会返回函数输出的参数类型
				outArgs = append(outArgs, reflect.Zero(fn.Type().Out(i)))
				continue
			}

			outArgs = append(outArgs, reflect.ValueOf(arg))
		}

		return outArgs
	}

	v := reflect.MakeFunc(fn.Type(), f)
	fn.Set(v)
}
