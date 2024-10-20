package main

import (
	"fmt"
	"net"

	"github.com/hwg1999/go_gateway/gateway_demo/base/proxy/unpack/codec"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	defer conn.Close()

	codec.Encode(conn, "hello world 0!!!")
}
