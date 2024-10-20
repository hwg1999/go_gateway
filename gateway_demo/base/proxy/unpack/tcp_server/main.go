package main

import (
	"fmt"
	"net"

	"github.com/hwg1999/go_gateway/gateway_demo/base/proxy/unpack/codec"
)

func main() {
	// 1.监听端口
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("listen fail, err: %v\n", err)
		return
	}

	// 2.接收请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			continue
		}

		// 3.创建协程
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	for {
		bt, err := codec.Decode(conn)
		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		str := string(bt)
		fmt.Printf("receive from client, data: %v", str)
	}
}
