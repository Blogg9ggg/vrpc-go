package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
	"vrpc/client"
	"vrpc/codec"
	"vrpc/server"
)

func startServer(addr chan string) {
	var foo Foo
	if err := server.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}

	// pick a free port
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error: ", err)
	}
	log.Println("start rpc server on", lis.Addr())
	addr <- lis.Addr().String()
	server.Accept(lis)
}

func day1() {
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple vrpc client
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// send options
	_ = json.NewEncoder(conn).Encode(codec.DefaultOption)
	cc := codec.NewGobCodec(conn)
	// send request & receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
	}
}

func day2() {
	log.SetFlags(0)
	addr := make(chan string)

	go startServer(addr)

	cli, err := client.Dial("tcp", <-addr)
	if err != nil {
		log.Fatal("dial error: ", err)
		return
	}

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	// send request & receive response
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			err = cli.Call(ctx, "Foo.Sum", args, &reply)
			if err != nil {
				log.Fatal("call error: ", err)
			}

			log.Println("receive reply:", reply)
		}(i)
	}
	wg.Wait()

	cli.Close()
}

type Foo int
type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}
func day3() {
	log.SetFlags(0)
	addr := make(chan string)

	go startServer(addr)

	cli, err := client.Dial("tcp", <-addr)
	if err != nil {
		log.Fatal("dial error: ", err)
		return
	}

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	// send request & receive response
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := Args{Num1: i, Num2: i * i}
			var reply int
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			err = cli.Call(ctx, "Foo.Sum", args, &reply)
			if err != nil {
				log.Fatal("call error:", err)
			}

			log.Println("receive reply:", reply)
		}(i)
	}
	wg.Wait()

	err = cli.Close()
	if err != nil {
		log.Fatal("client close error:", err)
	}
}

func startServer5(addr chan string) {
	var foo Foo
	_ = server.Register(&foo)
	server.HandleHTTP()

	l, _ := net.Listen("tcp", ":9999")
	addr <- l.Addr().String()
	_ = http.Serve(l, nil)
}

func call5(addrCh chan string) {
	client, _ := client.DialHTTP("tcp", <-addrCh)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := Args{Num1: i, Num2: i * i}
			var reply int
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			if err := client.Call(ctx, "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

func day5() {
	log.SetFlags(0)
	ch := make(chan string)
	go call5(ch)
	startServer5(ch)
}

func main() {
	day5()
}
