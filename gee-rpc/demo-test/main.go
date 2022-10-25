package main

import (
	gee_rpc "gee-rpc"
	"log"
	"net"
	"sync"
	"time"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}


func startServer(addr chan string) {
	var foo Foo
	err := gee_rpc.Register(&foo)
	if err != nil {
		log.Fatal("register error:", err)
	}


	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	gee_rpc.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)
	client, _ := gee_rpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

/*func main() {
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple gee-rpc client
	//conn, _ := net.Dial("tcp", <-addr)
	//defer func() { _ = conn.Close() }()
	client, _ := gee_rpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	time.Sleep(time.Second)
	time.Sleep(time.Second)
	// send options
	//_ = json.NewEncoder(conn).Encode(gee_rpc.DefaultOption)
	//cc := codec.NewGobCodec(conn)
	// send request & receive response
	//for i := 0; i < 5; i++ {
	//	h := &codec.Header{
	//		ServiceMethod: "Foo.Sum",
	//		Seq:           uint64(i),
	//	}
	//	_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
	//	//_ = cc.ReadHeader(h)
	//	var reply codec.Header
	//	_ = cc.ReadBody(&reply)
	//	log.Println("reply:", reply)
	//}
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			// 同步请求
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			// 异步请求
			//syncCall := client.Go("Foo.Sum", args, &reply, make(chan *gee_rpc.Call, 4))
			//replyDone := <-syncCall.Done
			//log.Println("reply:", replyDone)
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}*/