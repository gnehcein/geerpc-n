package main

import (
	"../client"
	. "../server"
	"context"
	"fmt"
	"net"

)

type Doo struct {}
type Args1 struct {
	Num1, Num2 int
}

func (f Doo) Sum(args Args1, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}
func (f Doo) Mul(args Args1, reply *int) error {
	*reply = args.Num1 * args.Num2
	return nil
}

func startServer1(addrCh chan string) {
	var doo Doo
	l, _ := net.Listen("tcp", ":0")
	server :=NewServer()
	err := server.Register(&doo)
	if err!=nil {
		panic(err)
	}
	addrCh <- l.Addr().String()
	server.Accept(l)
}

func main()  {
	ch1 := make(chan string)
	go startServer1(ch1)
	addr1 := <-ch1
	client,err := client.Dial("tcp",addr1)
	if err!=nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	arg1 := Args1{1,2}
	arg2 := Args1{5,10}
	var res int
	client.Call(ctx,"Doo.Sum",&arg1,&res)
	fmt.Println("sum",res)
	client.Call(ctx,"Doo.Mul",&arg2,&res)
	fmt.Println("mul",res)
}