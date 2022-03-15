package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	e := &Engine{}
	http.Handle("/", e)
	log.Println("starting server")
	log.Fatalln(http.ListenAndServe(":8081", nil))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

type Engine struct {
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//判断请求是否为websocket升级请求。
	if websocket.IsWebSocketUpgrade(r) {
		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			fmt.Printf("upgrade failed: %v\n", err)
			return
		}
		conn.WriteMessage(websocket.TextMessage, []byte("wxm.alming"))
		conn.SetCloseHandler(func(code int, text string) error {
			fmt.Println(code, text)
			return nil
		})
		go func() {
			for {
				t, c, _ := conn.ReadMessage()
				fmt.Println(t, string(c))
				if t == -1 {
					return
				}
			}
		}()
	} else {
		//处理普通请求
		fmt.Println("normal request coming")
	}
}

// reference: https://www.jianshu.com/p/b5e289be5fa1
