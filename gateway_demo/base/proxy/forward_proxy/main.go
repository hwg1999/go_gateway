package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Proxy struct{}

func main() {
	fmt.Println("Serve on :8080")
	http.Handle("/", &Proxy{})
	http.ListenAndServe(":8080", nil)
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	// step1. 浅拷贝对象，然后就再新增属性数据
	outReq := new(http.Request)
	*outReq = *req
	clientIP, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// step2. 请求下游
	transport := http.DefaultTransport
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	// step3. 把下游请求内容返回给上游
	for key, value := range res.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}

	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}
