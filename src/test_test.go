package src

import (
	"./client"
	"./server"
	"context"
	"log"
	"net"
	"testing"
	"time"
)

type Bar int

func (b Bar) Timeout(argv int, reply *int) error {
	time.Sleep(time.Second * 2)
	return nil
}

func startServer(addr chan string) {
	var b Bar
	_ = server.Register(&b)
	// pick a free port
	l, _ := net.Listen("tcp", ":0")
	addr <- l.Addr().String()
	server.Accept(l)
}

func TestClient_Call(t *testing.T) {
	t.Parallel()
	addrCh := make(chan string)
	go startServer(addrCh)
	addr := <-addrCh
	time.Sleep(time.Second)
	t.Run("client timeout", func(t *testing.T) {
		client, _ := client.Dial("tcp", addr)	//这里少写了一个参数：opt
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		var reply int
		_ = client.Call(ctx, "Bar.Timeout", 1, &reply)
		log.Println(reply)
	})
	t.Run("server handle timeout", func(t *testing.T) {
		client, _ := client.Dial("tcp", addr, &server.Option{
			Timeout: time.Second,
		})
		var reply int
		_ = client.Call(context.Background(), "Bar.Timeout", 1, &reply)

	})
}
